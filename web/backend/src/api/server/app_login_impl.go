/*
@Time : 2019/4/12 19:29 
@Author : yanKoo
@File : talk_cloud_app_login_impl
@Software: GoLand
@Description:
*/
package server

import (
	pb "api/talk_cloud"
	"context"
	"database/sql"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"model"
	"net/http"
	tg "pkg/group"
	tgc "pkg/group_cache"
	tgm "pkg/group_member"
	tm "pkg/msg"
	s "pkg/session"
	tu "pkg/user"
	tuc "pkg/user_cache"
	tuf "pkg/user_friend"
	"server/common/src/cache"
	"server/common/src/db"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
	"utils"
)

const (
	FIRST_LOGIN_DATA                = 1 // 初次登录返回的数据。比如用户列表，群组列表，该用户的个人信息
	OFFLINE_IM_MSG                  = 2 // 用户离线时的IM数据
	IM_MSG_FROM_UPLOAD_OR_WS_OR_APP = 3 // APP和web通过httpClient上传的文件信息、在线时通信的im数据
	KEEP_ALIVE_MSG                  = 4 // 用户登录后，每隔interval秒向stream发送一个消息，测试能不能连通
	LOGOUT_NOTIFY_MSG               = 5 // 用户掉线之后，通知和他在一个组的其他成员

	IM_MSG_FROM_UPLOAD_RECEIVER_IS_USER  = 1 // APP和web通过httpClient上传的IM信息是发给个人
	IM_MSG_FROM_UPLOAD_RECEIVER_IS_GROUP = 2 // APP和web通过httpClient上传的IM信息是发给群组

	USER_OFFLINE = 1 // 用户离线
	USER_ONLINE  = 2 // 用户在线

	UNREAD_OFFLINE_IM_MSG = 1 // 用户离线消息未读
	READ_OFFLINE_IM_MSG   = 2 // 用户离线消息已读

	CLIENT_EXCEPTION_EXIT = -1 // 客户端异常终止

	NOTIFY = 1 // 通知完一个
)

type TalkCloudServiceImpl struct {
}

type Task struct {
	DataChan          *pb.StreamResponse
	Receiver          []int32
	SenderId          int32
	KeepAliveClose    chan int32
	Option            chan int32 // 加上就是普通调用（http上传文件）
	LogoutNotifyCount chan int32
	done              chan int32
}

type DataContent struct {
	UId              chan int32
	Task             chan Task
	ExceptionalLogin chan int32
	KeepAliveClose   chan int32
	Option           chan int32 // 加上就是普通调用（http上传文件）
}

type DataSource interface{}

// gen im task
func NewImTask(senderId int32, receiver []int32, resp *pb.StreamResponse, ) *Task {
	return &Task{
		SenderId:          senderId,
		Receiver:          receiver,
		DataChan:          resp,
		KeepAliveClose:    make(chan int32, 1),
		Option:            make(chan int32, 1),
		LogoutNotifyCount: make(chan int32, 1),
		done:              make(chan int32, 1),
	}
}

func NewDataContent() DataContent {
	return DataContent{
		UId:              make(chan int32, 1),
		Task:             make(chan Task, 10000),
		KeepAliveClose:   make(chan int32, 1),
		ExceptionalLogin: make(chan int32, 1),
		Option:           make(chan int32, 1),
	}
}

var StreamMap = sync.Map{}

// 上传文件方式产生的IM数据推送
func (tcs *TalkCloudServiceImpl) ImMessagePublish(ctx context.Context, req *pb.ImMsgReqData) (*pb.ImMsgRespData, error) {
	var (
	//errMap sync.Map
	)
	dc := NewDataContent()
	if err := imMessagePublishDispatcher(&dc, req); err != nil {
		//errMap.Store(<-dc.Task, err)
	}
	// scheduler
	task := <-dc.Task
	go pushDataExecutorByOnce(task)
	uid := <-task.Option
	log.Printf("# %d im once done", uid)
	return &pb.ImMsgRespData{Result: &pb.Result{Msg: "push data done", Code: 200}, MsgCode: req.MsgCode}, nil
}

// 分发登录返回数据、IM离线数据、IM离线数据、Heartbeat
func (tcs *TalkCloudServiceImpl) DataPublish(srv pb.TalkCloud_DataPublishServer) error {
	var dc = NewDataContent()

	// msg dispatcher
	go pushDataDispatcher(&dc, srv)

	// scheduler
	go pushDataScheduler(&dc)

	// 重复登录就直接返回
	uid := <-dc.ExceptionalLogin
	err := srv.Send(&pb.StreamResponse{
		Res: &pb.Result{
			Msg:  "The user with id " + strconv.FormatInt(int64(uid), 10) + " is login already. please try again",
			Code: http.StatusUnauthorized,
		},
	})
	return err
}

// 分发上传文件方式产生的IM数据
func imMessagePublishDispatcher(dc *DataContent, ds DataSource) error {
	// 获取要发送的数据
	req := ds.(*pb.ImMsgReqData)

	offlineMem := make([]int32, 0)
	onlineMem := make([]int32, 0)

	log.Printf("grpc receive im from web : %+v", req)

	// 获取在线、离线用户id
	if req.ReceiverType == IM_MSG_FROM_UPLOAD_RECEIVER_IS_USER { // 发给单人
		// 判断是否在线
		v, ok := StreamMap.Load(req.ReceiverId)
		log.Println(req.ReceiverId, v, ok)
		//log.Printf("now dc.StreamMap map have: %+v， %p", dc.StreamMap, &dc.StreamMap)
		log.Println(req.ReceiverId)
		if !ok {
			log.Println(req.ReceiverId, " is offline")
			// 保存进数据库
			if err := tm.AddMsg(req, db.DBHandler); err != nil {
				log.Println("Add offline msg with error: ", err)
				return err
			}
		} else {
			onlineMem = append(onlineMem, req.ReceiverId)
		}
		onlineMem = append(onlineMem, req.Id)
	}

	if req.ReceiverType == IM_MSG_FROM_UPLOAD_RECEIVER_IS_GROUP { // 发送给群组成员，然后区分离线在线
		log.Printf("want send msg to group %d", req.ReceiverId)
		res, err := tgm.SelectDeviceIdsByGroupId(int(req.ReceiverId))
		if err != nil {
			log.Println("Add offline msg get group member with error: ", err)
			return err
		}

		log.Printf("the group %d has %+v", req.ReceiverId, res)
		for _, v := range res {
			log.Printf("now stream map have:%+v", StreamMap)
			_, ok := StreamMap.Load(int32(v))
			log.Println("the group member online state:", req.ReceiverId, v, ok)
			if ok {
				onlineMem = append(onlineMem, int32(v))
			} else {
				offlineMem = append(offlineMem, int32(v))
			}
		}

		// 存储离线消息
		log.Printf("the offline: %+v， the length is %d", offlineMem, len(offlineMem))
		if len(offlineMem) != 0 {
			if err := tm.AddMultiMsg(req, offlineMem, db.DBHandler); err != nil {
				log.Println("Add offline msg with error: ", err)
				return err
			}
		}
	}

	log.Println("web api want send to :", onlineMem)
	resp := &pb.StreamResponse{
		DataType:  IM_MSG_FROM_UPLOAD_OR_WS_OR_APP,
		ImMsgData: req,
		Res:       &pb.Result{Code: http.StatusOK, Msg: "receiver im message successful"},
	}
	// 发送在线用户消息
	if onlineMem != nil {
		dc.Task <- *NewImTask(req.Id, onlineMem, resp)
		log.Printf("dispatcher finish %+v <-||||-> %+v", req.Id, resp)
	}
	return nil
}

// Dispatcher 根据数据类型，调用不同的函数，产生不同的数据，然后把数据放到发送chan中
func pushDataDispatcher(dc *DataContent, ds DataSource) {
	interval := 5 // TODO 写入配置
	for {
		log.Println("dispatcher client msg")
		errMap := &sync.Map{}
		srv := ds.(pb.TalkCloud_DataPublishServer)

		data, _ := srv.Recv()

		if data == nil {
			log.Println("this stream is no data, maybe is offline")
			return
		}

		// 如果再次登录的用户在map中已存在，并且srv不同，那么就给dc的chan里面写一个终止信号
		if v, ok := StreamMap.Load(data.Uid); ok && v != srv {
			log.Printf("this here %+v, %+v", v, srv)
			log.Printf("this user # %d is login already", data.Uid)
			dc.ExceptionalLogin <- data.Uid
			return
		}
		// 更新stream和redis状态
		StreamMap.Store(data.Uid, srv)
		log.Println("login data.Uid:", data.Uid)

		if err := tuc.UpdateOnlineInCache(&pb.Member{Id: data.Uid, Online: USER_ONLINE}, cache.GetRedisClient()); err != nil {
			log.Println("Update user online state error:", err)
		}

		switch data.DataType {
		case FIRST_LOGIN_DATA:
			res, err := firstLoginData(dc, data, srv)
			if err != nil {
				errMap.Store(data.Uid, err)
			}
			re := make([]int32, 0)

			dc.Task <- *NewImTask(data.Uid, append(re, data.Uid), res)

		case OFFLINE_IM_MSG:
			res, err := GetOfflineImMsgFromDB(data)
			if err != nil {
				errMap.Store(data.Uid, err)
			}

			re := make([]int32, 0)
			dc.Task <- *NewImTask(data.Uid, append(re, data.Uid), res)

			// 如果连接存在，就3s往channel里面放一个数据
			go sendHeartbeat(dc, data, interval)

		case IM_MSG_FROM_UPLOAD_OR_WS_OR_APP:
			if err := imMessagePublishDispatcher(dc, data.ImMsg); err != nil {
				errMap.Store(data.Uid, err)
			}
		}
	}
}

// temp  临时用一下，后期统统会改用RabbitMQ
func pushDataScheduler(dContent *DataContent) {
	log.Printf("start Scheduler im msg")

	var (
		executor    = createExecutor()
		notifyTotal = int32(0) // 掉线通知的在线的人数
		exit        = make(chan int32, 1)
		tasks       = make([]Task, 0)
	)
	defer close(exit)

	tick := time.Tick(time.Second * 10)
	for {
		var activeExecutor chan<- Task
		var activeValue Task
		if len(tasks) > 0 {
			activeValue = tasks[0] // 取出第一个task
			activeExecutor = executor
			// 对每一个消息都进行监听，看是否需要退出
			go func(e chan int32, total *int32) {
				for {
					select {
					case uId := <-activeValue.KeepAliveClose:
						dContent.KeepAliveClose <- uId
						//e <- uId //不立刻关闭调度器，还得发完所有和他相关联的人
						//close(activeValue.KeepAliveClose)
						return
					case count := <-activeValue.LogoutNotifyCount:
						atomic.AddInt32(total, count)
						if atomic.LoadInt32(total) == activeValue.DataChan.NotifyTotal {
							e <- activeValue.DataChan.NotifyTotal
							//close(activeValue.LogoutNotifyCount)
							return
						}
					case <-activeValue.done:
						return
					}
				}
			}(exit, &notifyTotal)
		}

		// 接收任务
		select {
		case task := <-dContent.Task:
			tasks = append(tasks, task)
		case activeExecutor <- activeValue:
			tasks = tasks[1:]
		case u := <-exit:
			log.Printf("user # %d exit", u)
			return
		case <-tick:
			log.Printf("%d im task queue len = %d", activeValue.SenderId, len(tasks)) //TODO 合理退出，关闭调度器
		}
	}
}

// 创建执行者
func createExecutor() chan<- Task {
	ct := make(chan Task)
	go pushDataExecutor(ct)
	return ct
}

// Executor 推送登录返回数据、IM离线数据、IM离线数据、Heartbeat
func pushDataExecutor(ct chan Task) {
	var wg sync.WaitGroup
	for task := range ct {
		log.Println("receiver: ", task.Receiver)
		for _, receiverId := range task.Receiver {
			wg.Add(1)
			pushData(task, receiverId, task.DataChan, &wg)
			log.Println("will send to ", receiverId)
		}
		wg.Wait()
		task.Option <- task.SenderId
		task.done <- task.SenderId
	}

}

func pushDataExecutorByOnce(task Task) {
	var wg sync.WaitGroup
	log.Println("receiver: ", task.Receiver)
	for _, receiverId := range task.Receiver {
		wg.Add(1)
		pushData(task, receiverId, task.DataChan, &wg)
		log.Println("will send to ", receiverId)
	}
	wg.Wait()
	log.Printf("will ---------------------------------")
	task.Option <- task.SenderId

}

// 推送数据
func pushData(task Task, receiverId int32, resp *pb.StreamResponse, wg *sync.WaitGroup) {
	defer wg.Done()

	if task.DataChan.DataType == LOGOUT_NOTIFY_MSG {
		task.LogoutNotifyCount <- NOTIFY
	}

	log.Printf("the stream map have: %+v", StreamMap)
	if value, ok := StreamMap.Load(receiverId); ok {
		srv := value.(pb.TalkCloud_DataPublishServer)
		log.Printf("# %d receiver response: %+v", receiverId, resp)
		if err := srv.Send(resp); err != nil {
			// 发送失败处理
			processErrorSendMsg(err, task, receiverId, resp)
		} else {
			// 发送成功如果是离线数据（接收者等于stream id自己）就去更新状态
			log.Printf("send success. dc.senderId: %d, receiverId: %d", task.SenderId, receiverId)
			if task.SenderId == receiverId && resp.DataType == OFFLINE_IM_MSG {
				//  更新数据库里面的消息的状态
				if err := tm.SetMsgStat(receiverId, READ_OFFLINE_IM_MSG, db.DBHandler); err != nil {
					log.Println("Add offline msg with error: ", err)
				}
			}
		}
	} else {
		log.Println("can't find stream")
		task.KeepAliveClose <- receiverId // 退出这个连接
		// 不在线，存入数据库  //TODO 存储即时发送失败的消息
		if resp.DataType == IM_MSG_FROM_UPLOAD_OR_WS_OR_APP && task.SenderId != receiverId {
			// 把发送数据保存进数据库, 如果是离线数据就忽略
			// TODO 有个问题，如果发送的时候，还在线，函数走到这里，就掉线了， 但是发送方收到的回复仍然是已发送，发送方能看见的只有该用户已离线的状态
			log.Printf("send fail. dc.senderId: %d, receiverId: %d", task.SenderId, receiverId)
			if err := tm.AddMsg(resp.ImMsgData, db.DBHandler); err != nil {
				log.Println("Send fail and add offline msg with error: ", err)
			}
		}
	}
}

// 处理推送数据失败的情况
func processErrorSendMsg(err error, task Task, receiverId int32, resp *pb.StreamResponse) {
	log.Println("send msg fail with error: ", err)

	// 判断错误类型
	if errSC, _ := status.FromError(err); errSC.Code() == codes.Unavailable || errSC.Code() == codes.Canceled {
		// 1. 只要是发送失败，就认为对方离线，如果是发送Im数据，就保存到数据库
		task.KeepAliveClose <- receiverId

		// 删除map中的stream
		StreamMap.Delete(receiverId)
		log.Printf("# user %d is logout", receiverId)

		// 更新redis状态
		if err := tuc.UpdateOnlineInCache(&pb.Member{Id: receiverId, Online: USER_OFFLINE}, cache.GetRedisClient()); err != nil {
			log.Println("Update user online state error:", err)
		}

		log.Printf("%=v ---- %+v start add offline msg: ", resp, task)
		if resp.DataType == IM_MSG_FROM_UPLOAD_OR_WS_OR_APP && task.SenderId != receiverId {
			// 把发送数据保存进数据库, 如果是离线数据就忽略
			log.Printf("send fail. dc.senderId: %d, receiverId: %d", task.SenderId, receiverId)
			if err := tm.AddMsg(resp.ImMsgData, db.DBHandler); err != nil {
				log.Println("Send fail and add offline msg with error: ", err)
			}
		}
	}
}

// 初次登录应该返回的数据
func firstLoginData(dc *DataContent, data *pb.StreamRequest, srv pb.TalkCloud_DataPublishServer) (*pb.StreamResponse, error) {
	//time.Sleep(time.Second*20)
	//　验证用户名是否存在以及密码是否正确，然后就生成一个uuid session, 把sessionid放进metadata返回给客户端，
	//  然后之后的每一次连接都需要客户端加入这个metadata，使用拦截器，对用户进行鉴权
	if data.Name == "" || data.Passwd == "" {
		return &pb.StreamResponse{
			Res: &pb.Result{Code: 422, Msg: "用户名或密码不能为空"},
		}, errors.New("username or password can't be empty")
	}

	res, err := tu.SelectUserByKey(data.Name)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("App login error : %s", err)
		return &pb.StreamResponse{
			Res: &pb.Result{Code: 500, Msg: "User Login Process failed. Please try again later"},
		}, err
	}

	if err == sql.ErrNoRows {
		log.Printf("App login error : %s", err)
		return &pb.StreamResponse{
			Res: &pb.Result{Code: 500, Msg: "User is not exist error. Please try again later"},
		}, err
	}

	if res.PassWord != data.Passwd {
		log.Printf("App login error : %s", err)
		return &pb.StreamResponse{
			Res: &pb.Result{Code: 500, Msg: "User Login pwd error. Please try again later"},
		}, err
	}
	userInfo := &pb.Member{
		Id:          int32(res.Id),
		IMei:        res.IMei,
		UserName:    res.UserName,
		NickName:    res.NickName,
		UserType:    int32(res.UserType),
		LockGroupId: int32(res.LockGroupId),
		Online:      tuc.USER_ONLINE, // 登录就在线
	}

	var (
		errMap   = &sync.Map{}
		wg       sync.WaitGroup
		fList    = make(chan *pb.FriendsRsp, 1)
		gList    = make(chan *pb.GroupListRsp, 1)
		existErr bool
	)
	defer func() {
		close(fList)
		close(gList)
	}()

	// 如果再次登录的用户在map中已存在，并且srv不同，那么就给dc的chan里面写一个终止信号
	if v, ok := StreamMap.Load(int32(res.Id)); ok && v != srv {
		log.Printf("this here %+v, %+v", v, srv)
		dc.ExceptionalLogin <- int32(res.Id)
		log.Println("this user is login already")
		return &pb.StreamResponse{
			Res: &pb.Result{Code: 500, Msg: "this user is login already"},
		}, errors.New("the user is login already")
	}
	// 更新stream和redis状态
	log.Println("login data.Uid:", int32(res.Id))
	StreamMap.Store(int32(res.Id), srv)
	if err := tuc.UpdateOnlineInCache(&pb.Member{Id: int32(res.Id), Online: USER_ONLINE}, cache.GetRedisClient()); err != nil {
		log.Println("Update user online state error:", err)
	}

	//wg.Add(3) // 获取用户列表 用户群组列表，加入缓存, 插入session

	// 1. 处理登录session
	//go processSession(data, errMap, &wg)

	// 2. 获取好友列表
	getFriendList(int32(res.Id), fList, errMap, &wg)

	// 3. 群组列表
	getGroupList(int32(res.Id), gList, errMap, &wg)

	// 4. 将用户信息添加进redis
	addUserInfoToCache(userInfo, &wg)

	//wg.Wait()

	//遍历该map，参数是个函数，该函数参的两个参数是遍历获得的key和value，返回一个bool值，当返回false时，遍历立刻结束。
	errMap.Range(func(k, v interface{}) bool {
		err = v.(error)
		if err != nil {
			log.Println(k, " gen error: ", err)
			existErr = true
			return false
		}
		return true
	})
	if existErr {
		return &pb.StreamResponse{
			DataType: FIRST_LOGIN_DATA,
			Res:      &pb.Result{Code: 500, Msg: "User Login pwd error. Please try again later"},
		}, err
	}
	return &pb.StreamResponse{
		DataType: FIRST_LOGIN_DATA,
		LoginResp: &pb.FirstLoginData{
			UserInfo:   userInfo,
			FriendList: (<-fList).FriendList,
			GroupList:  (<-gList).GroupList,
		},
		Res: &pb.Result{Code: 200, Msg: data.Name + " login successful"},
	}, nil
}

// 处理登录session validate TODO 给stream模式加metadata
func processSession(req *pb.StreamRequest, errMap *sync.Map, wg *sync.WaitGroup) {
	defer wg.Done()
	sessionId, err := utils.NewUUID()
	if err != nil {
		log.Panicf("session id is error%s", err)
		errMap.Store("processSession", err)
		return
	}

	ct := time.Now().UnixNano() / 1000000
	ttl := ct + 30*60*1000 // Severside session valid time: 30 min
	ttlStr := strconv.FormatInt(ttl, 10)

	sInfo := &model.SessionInfo{SessionID: sessionId, UserName: req.Name, UserPwd: req.Passwd, TTL: ttlStr}
	if err := s.InsertSession(sInfo); err != nil {
		log.Panicf("session id insert is error%s", err)
		errMap.Store("processSession", err)
		return
	}
}

// 获取好友列表
func getFriendList(uid int32, fList chan *pb.FriendsRsp, errMap *sync.Map, wg *sync.WaitGroup) {
	log.Println("get FriendList start")
	var err error
	fl, _, err := tuf.GetFriendReqList(int32(uid), db.DBHandler)
	if err != nil {
		errMap.Store("getFriendList", err)
		fList <- nil
	} else {
		fList <- fl
	}
	log.Println("get FriendList done")
	//wg.Done()
}

// 获取群组列表
func getGroupList(uid int32, gList chan *pb.GroupListRsp, errMap *sync.Map, wg *sync.WaitGroup) {
	log.Println("Get group list start")
	// 先去缓存取，取不出来再去mysql取
	gl, err := tuc.GetGroupListFromRedis(int32(uid), cache.GetRedisClient())
	if err != nil && err != sql.ErrNoRows {
		log.Println("No find In CacheError")
		errMap.Store("getGroupList", err)
		log.Printf("get GroupList%v", err)
		gList <- gl
		return
	}

	if err == sql.ErrNoRows {
		log.Println("redis is not find")
		for {
			gl, _, err = tg.GetGroupListFromDB(int32(uid), db.DBHandler)
			if err != nil {
				errMap.Store("getGroupList", err)
				break
			}
			log.Println("start update redis")
			// 新增到缓存 更新两个地方，首先，每个组的信息要更新，就是group data，记录了群组的id和名字
			if err := tgc.AddGroupInCache(gl, cache.GetRedisClient()); err != nil {
				errMap.Store("getGroupList", err)
				break
			}

			// 其次更新一个userSet  就是一个组里有哪些用户
			if err := tuc.AddUserInGroupToCache(gl, cache.GetRedisClient()); err != nil {
				errMap.Store("getGroupList", err)
				break
			}

			// 每个用户的信息
			for _, g := range gl.GroupList {
				for _, u := range g.UsrList {
					if err := tuc.AddUserDataInCache(&pb.Member{
						Id:          u.Uid,
						IMei:        u.Imei,
						UserName:    u.Name,
						Online:      u.Online,
						LockGroupId: u.LockGroupId,
					}, cache.GetRedisClient()); err != nil {
						log.Println("Add user information to cache with error: ", err)
					}
				}
			}

			// 每一个群组拥有的成员
			for _, v := range gl.GroupList {
				if err := tgc.AddGroupCache(v.UsrList, v, cache.GetRedisClient()); err != nil {
					errMap.Store("AddGroupCache", err)
					break
				}
			}
			break
		}
	}

	gList <- gl
	log.Println("Get group list done")
}

// 增加缓存 TODO 没有错误
func addUserInfoToCache(userInfo *pb.Member, wg *sync.WaitGroup) {
	log.Println("Add User Info into cache start")

	if err := tuc.AddUserDataInCache(userInfo, cache.GetRedisClient()); err != nil {
		log.Println("Add user information to cache with error: ", err)
	}
	log.Println("Add User Info into cache done")
}

// 返回的IM离线数据
func GetOfflineImMsgFromDB(req *pb.StreamRequest) (*pb.StreamResponse, error) {
	// 去数据库拉取离线数据
	log.Println("start get offline im msg")
	offlineMsg, err := tm.GetMsg(req.Uid, UNREAD_OFFLINE_IM_MSG, db.DBHandler)
	if err != nil {
		log.Println("Get offline msg fail with error:", err)
		return nil, err
	}
	log.Printf("get offline msg %+v", offlineMsg)

	var (
		idIndexSMap   = map[int32]int{}
		idIndexGMap   = map[int32]int{}
		respPkgSingle = make([]*pb.OfflineImMsg, 0)
		respPkgGroup  = make([]*pb.OfflineImMsg, 0)
		idxG          = 0
		idxS          = 0
	)
	// 遍历离线数据集，记录数据用户id和位置

	for _, msg := range offlineMsg {
		if msg.ReceiverType == IM_MSG_FROM_UPLOAD_RECEIVER_IS_USER {
			if v, ok := idIndexSMap[msg.Id]; ok {
				// 已经发现了这个用户的一条消息，那么就把消息加到对应的切片下的
				respPkgSingle[v].ImMsgData = append(respPkgSingle[v].ImMsgData, msg)
			} else {
				// 首次找到这个用户的第一条单人消息，就respPackage添加一个slice，并记录index
				var userMsgs = &pb.OfflineImMsg{
					SenderId:        msg.Id,
					Name:            msg.SenderName,
					MsgReceiverType: IM_MSG_FROM_UPLOAD_RECEIVER_IS_USER,
				}
				userMsgs.ImMsgData = append(make([]*pb.ImMsgReqData, 0), msg)
				respPkgSingle = append(respPkgSingle, userMsgs)
				idIndexSMap[msg.Id] = idxS
				idxS++

			}
		}

		// 群组
		if msg.ReceiverType == IM_MSG_FROM_UPLOAD_RECEIVER_IS_GROUP {
			if v, ok := idIndexGMap[msg.ReceiverId]; ok {
				// 已经发现了这个用户的一条消息，那么就把消息加到对应的切片下的
				log.Printf("v %d, s %v msg %+v", v, ok, msg)
				respPkgGroup[v].ImMsgData = append(respPkgGroup[v].ImMsgData, msg)
			} else {
				// 首次找到这个用户的第一条单人消息，就respPackage添加一个slice，并记录index
				var userMsgs = &pb.OfflineImMsg{
					GroupId:         msg.ReceiverId,
					Name:            msg.ReceiverName,
					MsgReceiverType: IM_MSG_FROM_UPLOAD_RECEIVER_IS_GROUP,
				}
				userMsgs.ImMsgData = append(make([]*pb.ImMsgReqData, 0), msg)
				respPkgGroup = append(respPkgGroup, userMsgs)
				idIndexGMap[msg.ReceiverId] = idxG
				idxG++

			}
		}
	}

	log.Printf("%+v \n %+v", respPkgSingle, respPkgGroup)

	return &pb.StreamResponse{
		OfflineImMsgResp: &pb.OfflineImMsgResp{
			OfflineSingleImMsgs: respPkgSingle,
			OfflineGroupImMsgs:  respPkgGroup},
		DataType: OFFLINE_IM_MSG}, nil
}

// 服务器主动发送Heartbeat
func sendHeartbeat(dc *DataContent, data *pb.StreamRequest, interval int) {
	// 使用定时器，收到登录请求之后每隔5s后发送一次数据
	timerTask := time.NewTicker(time.Second * time.Duration(interval))
	receiverId := make([]int32, 0)
timerLoop:
	for {
		select {
		case <-timerTask.C:
		loop:
			for {
				select {
				case uid := <-dc.KeepAliveClose:
					log.Printf("client %d close", uid)
					log.Printf("now dc stream : %+v", StreamMap)
					// 往dc里面写掉线通知
					go notifyOfflineToOther(dc, uid)
					break timerLoop
				default:
					dc.Task <- *NewImTask(data.Uid, append(receiverId, data.Uid), keepAliveData(data))
					break loop
				}
			}
		}
	}

}

// 掉线通知
func notifyOfflineToOther(dc *DataContent, uId int32) {
	// 3. 群组列表
	var (
		errMap      = &sync.Map{}
		selfGList   = make(chan *pb.GroupListRsp, 1)
		notifyTotal int32
	)

	uInfo, _ := tuc.GetUserFromCache(uId)
	//uLocation, _ := tlc.GetUserLocationInCache(uId, cache.GetRedisClient())

	getGroupList(uId, selfGList, errMap, nil)
	gl := <-selfGList

	for _, g := range gl.GroupList {
		for _, u := range g.UsrList {
			if u.Uid == uId || u.Online == tuc.USER_ONLINE {
				notifyTotal++
			}
		}
	}

	for _, g := range gl.GroupList {
		for _, u := range g.UsrList {
			if u.Uid == uId || u.Online == tuc.USER_ONLINE {
				// 对于群里每一位都要通知到
				recvId := make([]int32, 0)
				recvId = append(recvId, u.Uid)

				gList := make(chan *pb.GroupListRsp, 1)
				getGroupList(uId, gList, errMap, nil)
				resp := &pb.StreamResponse{
					DataType:    LOGOUT_NOTIFY_MSG,
					NotifyTotal: notifyTotal,
					LogoutNotify: &pb.LogoutNotify{
						UserInfo: uInfo,
						//UserLocation: uLocation.GpsInfo,
						GroupList: (<- gList).GroupList,
					},
					Res: &pb.Result{Code: 200, Msg: strconv.FormatInt(int64(uId), 10) + " notify successful"},
				}

				log.Printf("will send logout notify to %+v", recvId)
				dc.Task <- *NewImTask(uId, recvId, resp)
			}
		}
	}

}

// 心跳数据
func keepAliveData(req *pb.StreamRequest) *pb.StreamResponse {
	return &pb.StreamResponse{
		DataType: KEEP_ALIVE_MSG,
		KeepAlive: &pb.KeepAlive{
			Uid: req.Uid,
			SYN: 1,
		},
		Res: &pb.Result{
			Msg:  "The user with id " + strconv.FormatInt(int64(req.Uid), 10) + " send heartbeat",
			Code: http.StatusOK,
		},
	}
}

// 用户登录
func (tcs *TalkCloudServiceImpl) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRsp, error) {
	// 登录就开启任务四个协程，一个用来去返回初步登录应该获取的数据，一个用来推送发送保活数据，一个用来发送im数据

	// 1. TODO 用sync.map将上下文保存? 记录在线用户

	// 2.goroutine，建立一个客户端，调用loginProess方法

	// 3.goroutine，建立一个客户端，调用KeepAlive方法

	// 4.goroutine，建立一个客户端，调用ImMessagePublish方法
	//loginProscess(ctx context.Context, req *pb.LoginReq) (*pb.LoginRsp, error)
	log.Println("enter login")
	//time.Sleep(time.Second*20)
	//　验证用户名是否存在以及密码是否正确，然后就生成一个uuid session, 把sessionid放进metadata返回给客户端，
	//  然后之后的每一次连接都需要客户端加入这个metadata，使用拦截器，对用户进行鉴权
	if req.Name == "" || req.Passwd == "" {
		return &pb.LoginRsp{Res: &pb.Result{Code: 422, Msg: "用户名或密码不能为空"}}, nil
	}

	res, err := tu.SelectUserByKey(req.Name)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("App login error : %s", err)
		loginRsp := &pb.LoginRsp{
			Res: &pb.Result{
				Code: 500,
				Msg:  "User Login Process failed. Please try again later"},
		}
		return loginRsp, nil
	}

	if err == sql.ErrNoRows {
		log.Printf("App login error : %s", err)
		loginRsp := &pb.LoginRsp{
			Res: &pb.Result{
				Code: 500,
				Msg:  "User is not exist error. Please try again later"},
		}
		return loginRsp, nil
	}

	if res.PassWord != req.Passwd {
		log.Printf("App login error : %s", err)
		loginRsp := &pb.LoginRsp{
			Res: &pb.Result{
				Code: 500,
				Msg:  "User Login pwd error. Please try again later"},
		}
		return loginRsp, nil
	}
	userInfo := &pb.Member{
		Id:          int32(res.Id),
		IMei:        res.IMei,
		UserName:    res.UserName,
		NickName:    res.NickName,
		UserType:    int32(res.UserType),
		LockGroupId: int32(res.LockGroupId),
		Online:      tuc.USER_ONLINE, // 登录就在线
	}

	var (
		errMap   = &sync.Map{}
		wg       sync.WaitGroup
		fList    = make(chan *pb.FriendsRsp, 1)
		gList    = make(chan *pb.GroupListRsp, 1)
		existErr bool
	)
	defer func() {
		close(fList)
		close(gList)
	}()
	//wg.Add(3) // 获取用户列表 用户群组列表，加入缓存, 插入session

	// 1. 处理登录session
	//go processSession(req, errMap, &wg)

	// 2. 获取好友列表
	getFriendList(int32(res.Id), fList, errMap, &wg)

	// 3. 群组列表
	getGroupList(int32(res.Id), gList, errMap, &wg)

	// 4. 将用户信息添加进redis
	addUserInfoToCache(userInfo, &wg)

	//wg.Wait()

	//遍历该map，参数是个函数，该函数参的两个参数是遍历获得的key和value，返回一个bool值，当返回false时，遍历立刻结束。
	errMap.Range(func(k, v interface{}) bool {
		err = v.(error)
		if err != nil {
			log.Println(k, " gen error: ", err)
			existErr = true
			return false
		}
		return true
	})
	if existErr {
		return &pb.LoginRsp{Res: &pb.Result{Code: 500, Msg: "process error, please try again"}}, nil
	}

	loginRep := &pb.LoginRsp{
		UserInfo:   userInfo,
		FriendList: (<-fList).FriendList,
		GroupList:  (<-gList).GroupList,
		Res:        &pb.Result{Code: 200, Msg: req.Name + " login successful"},
	}

	return loginRep, nil
}
