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
	cfgGs "configs/grpc_server"
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
	tlc "pkg/location"
	//tlc "pkg/location"
	tm "pkg/msg"
	s "pkg/session"
	tu "pkg/user"
	tuc "pkg/user_cache"
	tuf "pkg/user_friend"
	"server/common/src/cache"
	"server/common/src/db"
	"strconv"
	"sync"
	"time"
	"utils"
)

const (
	FIRST_LOGIN_DATA                = 1 // 初次登录返回的数据。比如用户列表，群组列表，该用户的个人信息
	OFFLINE_IM_MSG                  = 2 // 用户离线时的IM数据
	IM_MSG_FROM_UPLOAD_OR_WS_OR_APP = 3 // APP和web通过httpClient上传的文件信息、在线时通信的im数据
	KEEP_ALIVE_MSG                  = 4 // 用户登录后，每隔interval秒向stream发送一个消息，测试能不能连通
	LOGOUT_NOTIFY_MSG               = 5 // 用户掉线之后，通知和他在一个组的其他成员
	LOGIN_NOTIFY_MSG                = 6 // 用户上线之后，通知和他在一个组的其他成员
	SOS_MSG                         = 7 // 用户按SOS按键呼救

	IM_MSG_FROM_UPLOAD_RECEIVER_IS_USER  = 1 // APP和web通过httpClient上传的IM信息是发给个人
	IM_MSG_FROM_UPLOAD_RECEIVER_IS_GROUP = 2 // APP和web通过httpClient上传的IM信息是发给群组

	USER_OFFLINE = 1 // 用户离线
	USER_ONLINE  = 2 // 用户在线

	UNREAD_OFFLINE_IM_MSG = 1 // 用户离线消息未读
	READ_OFFLINE_IM_MSG   = 2 // 用户离线消息已读

	CLIENT_EXCEPTION_EXIT = -1 // 客户端异常终止

	NOTIFY = 1 // 通知完一个

	WORK_BY_GORONTINE = 2
	WORK_BY_NORMAL    = 1 // 普通调用
)

type ImEngine struct {
}

// 分发的任务
type Task struct {
	Data     *pb.StreamResponse // 具体的消息
	Receiver []int32
	SenderId int32
}

// im推送上下文
type DataContext struct {
	UId              chan int32
	Task             chan Task
	ExceptionalLogin chan int32 // 重复登录
}

// im推送的client
type Client struct {
	WorkType  int32 // 开启协程还是不用  1是不用 2是开协程
	LongLived bool  // 分发多次数据
	Dc        *DataContext
	Ds        DataSource
	Cf        ClientFunc
}

// client需要实现的方法
type ClientFunc interface {
	Dispatcher(dc *DataContext, ds DataSource)
	DispatcherScheduler(dc *DataContext, longLived bool)
}

// 消息数据来源
type DataSource interface{}

// 全局消息队列
var TQ = struct {
	Tasks chan Task
	m     sync.Mutex
}{Tasks: make(chan Task, 10000000)}

// 全局map，记录在线人 // TODO 后期修改
var StreamMap = sync.Map{}

func init() {
	ImEngine{}.Run()
}

func (ie ImEngine) Run() {
	// 消息推送exec TODO 暂时只用一个scher
	go ExcecScheduler()

	// redis持续获取im数据，dispatcher
	JanusPttMsgPublish()
}

// gen im task
func NewImTask(senderId int32, receiver []int32, resp *pb.StreamResponse, ) *Task {
	return &Task{
		SenderId: senderId,
		Receiver: receiver,
		Data:     resp,
	}
}

// gen 数据上下文
func NewDataContent() *DataContext {
	return &DataContext{
		UId:  make(chan int32, 1),
		Task: make(chan Task, 1000),
	}
}

// client运行分发数据
func (c *Client) Run() {
	if c.WorkType == WORK_BY_GORONTINE {
		go c.Cf.Dispatcher(c.Dc, c.Ds)
		go c.Cf.DispatcherScheduler(c.Dc, c.LongLived)
	}
	if c.WorkType == WORK_BY_NORMAL {
		c.Cf.Dispatcher(c.Dc, c.Ds)
		c.Cf.DispatcherScheduler(c.Dc, c.LongLived)
	}
}

func ExcecScheduler() {
	var tasks []Task
	var executor = CreateExecutor()
	tick := time.NewTicker(time.Second * time.Duration(5))
	for {
		var activeExecu chan Task
		var activeTask Task
		if len(tasks) > 0 {
			activeExecu = executor
			activeTask = tasks[0]
		}
		select {
		case t := <-TQ.Tasks:
			tasks = append(tasks, t)
		case activeExecu <- activeTask:
			tasks = tasks[1:]
		case <-tick.C:
			log.Printf("now task queue len:%d", len(tasks))
		}
	}
}

func CreateExecutor() chan Task {
	tc := make(chan Task)
	go pushDataExecutor(tc)
	return tc
}

// 分发上传文件方式产生的IM数据
func imMessagePublishDispatcher(dc *DataContext, ds DataSource) {
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
}

// Dispatcher 根据数据类型，调用不同的函数，产生不同的数据，然后把数据放到发送chan中
func pushDataDispatcher(dc *DataContext, ds DataSource) {
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
		log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>login data.Uid>>>>:", data.Uid)
		if err := tuc.UpdateOnlineInCache(&pb.Member{Id: data.Uid, Online: USER_ONLINE}, cache.GetRedisClient()); err != nil {
			log.Println("Update user online state error:", err)
		}
		log.Printf(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>start gen data for dispatcher")
		switch data.DataType {

		/*case FIRST_LOGIN_DATA:
			res, err := firstLoginData(dc, data, srv)
			if err != nil {
				errMap.Store(data.Uid, err)
			}
			re := make([]int32, 0)
			dc.Task <- *NewImTask(data.Uid, append(re, data.Uid), res)*/

		case OFFLINE_IM_MSG:
			res, err := GetOfflineImMsgFromDB(data)
			if err != nil {
				errMap.Store(data.Uid, err)
			}

			re := make([]int32, 0)
			dc.Task <- *NewImTask(data.Uid, append(re, data.Uid), res)

			// 如果连接存在，就3s往channel里面放一个数据
			go sendHeartbeat(dc, data, cfgGs.Interval, srv)

			// 往dc里面写上线通知
			go notifyToOther(dc, data.Uid, LOGIN_NOTIFY_MSG)
		case IM_MSG_FROM_UPLOAD_OR_WS_OR_APP:
			imMessagePublishDispatcher(dc, data.ImMsg)
		}
	}
}

// temp  临时用一下，后期统统会改用RabbitMQ
func dispatcherScheduler(dContent *DataContext, multiSend bool) {
	log.Printf("start Scheduler im msg")
	var notify int32
	tick := time.Tick(time.Minute * 5)
	for {
		// 接收任务
		select {
		case t := <-dContent.Task:
			go func() { TQ.Tasks <- t }()
			//log.Printf("///////////%T///////%+v", t, t)
			if t.Data.DataType == LOGOUT_NOTIFY_MSG {
				notify++
				if notify == t.Data.NotifyTotal {
					log.Printf("notify: %d, total: %d", notify, t.Data.NotifyTotal)
					return
				}

			}
			if !multiSend {
				log.Printf("only scheduler once")
				return
			}
		case <-tick:
			log.Printf("single im task queue len = %d", len(dContent.UId)) //TODO 合理退出，关闭调度器
		}
	}
}

// Executor 推送登录返回数据、IM离线数据、IM离线数据、Heartbeat
func pushDataExecutor(ct chan Task) {
	var wg sync.WaitGroup
	for {
		select {
		case task := <-ct:
			log.Println("global executor receiver: ", task.Receiver)
			for _, receiverId := range task.Receiver {
				wg.Add(1)
				go pushData(task, receiverId, task.Data, &wg)
				log.Println("will send to ", receiverId)
			}
			wg.Wait()
		}
	}
}

// 推送数据
func pushData(task Task, receiverId int32, resp *pb.StreamResponse, wg *sync.WaitGroup) {
	defer wg.Done()

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
		log.Println("can't find stream") //TODO 就依靠那边心跳了，这里就不管发送失败了
		// 存储即时发送失败的消息
		if resp.DataType == IM_MSG_FROM_UPLOAD_OR_WS_OR_APP && task.SenderId != receiverId {
			// 把发送数据保存进数据库, 如果是离线数据就忽略
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

		//log.Printf("%+v ---- %+v start add offline msg: ", resp, task)
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
func firstLoginData(dc *DataContext, data *pb.StreamRequest, srv pb.TalkCloud_DataPublishServer) (*pb.StreamResponse, error) {
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
			log.Println("start update redis GetGroupListFromDB")
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
						NickName:    u.Name,
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

// 增加缓存
func addUserInfoToCache(userInfo *pb.Member, wg *sync.WaitGroup) {

	log.Printf("will get rediscli, now redis pool info :%+v |<<<>>>| idleCount: %+v", cache.RedisPool, cache.RedisPool.IdleCount())
	redisCli := cache.GetRedisClient()
	if err := tuc.AddUserDataInCache(userInfo, redisCli); err != nil {
		log.Println("Add user information to cache with error: ", err)
	}
	log.Println("addUserInfoToCache done")
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
func sendHeartbeat(dc *DataContext, data *pb.StreamRequest, interval int, srv pb.TalkCloud_DataPublishServer) {
	// 使用定时器，收到登录请求之后每隔5s后发送一次数据
	timerTask := time.NewTicker(time.Second * time.Duration(interval))
	resp := keepAliveData(data)
	receiverId := make([]int32, 0)
	receiverId = append(receiverId, data.Uid)

	for {
		select {
		case <-timerTask.C:
			log.Printf("# %d receiver response: %+v", data.Uid, resp)
			if value, ok := StreamMap.Load(data.Uid); ok {
				srv := value.(pb.TalkCloud_DataPublishServer)
				if err := srv.Send(resp); err != nil {
					//if errSC, _ := status.FromError(err); errSC.Code() == codes.Unavailable || errSC.Code() == codes.Canceled {
					// 只要是发送失败，就认为对方离线
					log.Printf("client %d close with %+v", data.Uid, err)
					log.Printf("now dc stream : %+v", StreamMap)

					// 删除map中的stream
					StreamMap.Delete(data.Uid)
					log.Printf("# user %d is logout", data.Uid)

					// 更新redis状态
					if err := tuc.UpdateOnlineInCache(&pb.Member{Id: data.Uid, Online: USER_OFFLINE}, cache.GetRedisClient()); err != nil {
						log.Println("Update user online state error:", err)
					}
					// 往dc里面写掉线通知
					go notifyToOther(dc, data.Uid, LOGOUT_NOTIFY_MSG)
					return
					//}
				}
			}
		}

	}

}

// 上线通知所有人，掉线通知所有人、sos通知
func notifyToOther(dc *DataContext, uId int32, notifyType int32) {
	var (
		errMap      = &sync.Map{}
		selfGList   = make(chan *pb.GroupListRsp, 1)
		notifyTotal = int32(0)
	)
	log.Printf("notify root id :%d", uId)
	uInfo, _ := tuc.GetUserFromCache(uId)
	_, uLocation, _ := tlc.GetUserLocationInCache(uId, cache.GetRedisClient())

	getGroupList(uId, selfGList, errMap, nil)
	gl := <-selfGList
	if gl != nil && uInfo != nil {
		for _, g := range gl.GroupList {
			for _, u := range g.UsrList {
				if u.Uid == uId || u.Online == tuc.USER_ONLINE {
					//log.Printf("will notify *******************------------------------------------------------%d", u.Uid)
					notifyTotal++
				}
			}
		}
		log.Printf("notify total: %d", notifyTotal)
		for _, g := range gl.GroupList {
			for _, u := range g.UsrList {
				if u.Uid == uId || u.Online == tuc.USER_ONLINE {
					// 对于群里每一位都要通知到
					recvId := make([]int32, 0)
					recvId = append(recvId, u.Uid)

					gList := make(chan *pb.GroupListRsp, 1)
					getGroupList(u.Uid, gList, errMap, nil)
					//if gList != nil {
					resp := &pb.StreamResponse{
						DataType:    notifyType,
						NotifyTotal: notifyTotal,
						Notify: &pb.LoginOrLogoutNotify{
							UserInfo:     uInfo,
							UserLocation: uLocation,
							GroupList:    (<-gList).GroupList,
						},
						Res: &pb.Result{Code: 200, Msg: strconv.FormatInt(int64(u.Uid), 10) + " notify successful"},
					}
					if notifyType == SOS_MSG {
						resp.Notify.GroupList = nil
					}
					log.Printf("will send %d notify to %+v", notifyType, recvId)
					dc.Task <- *NewImTask(uId, recvId, resp)
					//}

				}
			}
		}
	}
	log.Printf("cant load notify self info %v----------%v", gl, uInfo)
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
