/*
@Time : 2019/4/12 19:29 
@Author : yanKoo
@File : talk_cloud_app_login_impl
@Software: GoLand
@Description:
*/
package server

import (
	"database/sql"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	pb "server/grpc-server/api/talk_cloud"
	"server/grpc-server/cache"
	cfgGs "server/grpc-server/configs/grpc_server"
	tg "server/grpc-server/dao/group"
	tgc "server/grpc-server/dao/group_cache"
	tgm "server/grpc-server/dao/group_member"
	tlc "server/grpc-server/dao/location"
	"server/grpc-server/log"
	"server/web-api/model"
	//tlc "server/web-api/dao/location"
	tm "server/grpc-server/dao/msg"
	s "server/grpc-server/dao/session"
	tu "server/grpc-server/dao/user"
	tuc "server/grpc-server/dao/user_cache"
	tuf "server/grpc-server/dao/user_friend"
	"server/grpc-server/db"
	"server/grpc-server/utils"
	"strconv"
	"sync"
	"time"
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
}{Tasks: make(chan Task, 100)}

// 全局map，记录在线人 // TODO 后期修改
var StreamMap = StreamMspSt{
	Streams: make(map[int32]interface{}),
}

func init() {
	ImEngine{}.Run()
}

func (ie ImEngine) Run() {
	// 消息推送exec TODO 暂时只用一个scher
	go ExcecScheduler()

	// 根据时间间隔，检查stream map和redis
	go syncStreamWithRedis(cfgGs.Interval)

	// redis持续获取im数据，dispatcher
	JanusPttMsgPublish()
}

//间隔interval 更新stream map
func syncStreamWithRedis(Interval int) {
	// 遍历map，检测redis是否
	// TODO 存在很严重的问题，比如，在这急短的时间内同步stream和redis
	for {
		for k, v := range StreamMap.Streams {
			log.Log.Printf("check %d, %+v", k, v)
			res, _ := tuc.GetUserStatusFromCache(k, cache.GetRedisClient())
			log.Log.Printf("Get uid:%d status: %+v\n", k, res)
			if res == USER_OFFLINE {
				log.Log.Printf("Will del uid: %d\n", k)
				StreamMap.Del(k)
				if err := tuc.UpdateOnlineInCache(&pb.Member{Id: k, Online: USER_OFFLINE}, cache.GetRedisClient()); err != nil {
					log.Log.Println("Update user online state error:", err)
				}
				// 往q里面写离线通知
				go notifyToOther(TQ.Tasks, k, LOGOUT_NOTIFY_MSG)
			}
		}
		time.Sleep(time.Second * time.Duration(cfgGs.Interval))
	}
}

type StreamMspSt struct {
	Streams map[int32]interface{}
	Lock    sync.RWMutex
}

func (s StreamMspSt) Get(k int32) interface{} {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	return s.Streams[k]
}

func (s StreamMspSt) Len() int {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	return len(s.Streams)
}

func (s StreamMspSt) Set(k int32, v interface{}) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.Streams[k] = v
}

func (s StreamMspSt) Del(k int32) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	delete(s.Streams, k)
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
		Task: make(chan Task, 100),
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
			log.Log.Debugf("now task queue len:%d", len(tasks))
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

	log.Log.Printf("grpc receive im from web : %+v", req)

	// 获取在线、离线用户id
	if req.ReceiverType == IM_MSG_FROM_UPLOAD_RECEIVER_IS_USER { // 发给单人
		// 判断是否在线
		v := StreamMap.Get(req.ReceiverId)
		log.Log.Println(req.ReceiverId, v)
		log.Log.Printf("now dc.StreamMap map have: %+v， %p", StreamMap, &StreamMap)
		if v != nil {
			log.Log.Println(req.ReceiverId, " is online")
			// 保存进数据库
			if err := tm.AddMsg(req, db.DBHandler); err != nil {
				log.Log.Errorf("Add offline msg with error: ", err)
			}
			onlineMem = append(onlineMem, req.ReceiverId)
		} else {
			log.Log.Println(req.ReceiverId, " is offline")
			offlineMem = append(offlineMem, req.ReceiverId)
		}
		onlineMem = append(onlineMem, req.Id)
	}

	if req.ReceiverType == IM_MSG_FROM_UPLOAD_RECEIVER_IS_GROUP { // 发送给群组成员，然后区分离线在线
		log.Log.Printf("want send msg to group %d", req.ReceiverId)
		res, err := tgm.SelectDeviceIdsByGroupId(int(req.ReceiverId))
		if err != nil {
			log.Log.Errorln("Add offline msg get group member with error: ", err)
		}

		log.Log.Printf("the group %d has %+v", req.ReceiverId, res)
		log.Log.Printf("now dc.StreamMap map have: %+v， %p", StreamMap, &StreamMap)
		for _, v := range res {
			log.Log.Printf("now stream map have:%+v", StreamMap)
			srv := StreamMap.Get(int32(v))
			log.Log.Infof("the group # %d member %d online state: %+v", req.ReceiverId, v, srv)
			if srv != nil {
				onlineMem = append(onlineMem, int32(v))
			} else {
				offlineMem = append(offlineMem, int32(v))
			}
		}
		//onlineMem = append(onlineMem, req.Id)
		// 存储离线消息
		log.Log.Printf("the offline: %+v， the length is %d", offlineMem, len(offlineMem))
		if len(offlineMem) != 0 {
			if err := tm.AddMultiMsg(req, offlineMem, db.DBHandler); err != nil {
				log.Log.Println("Add offline msg with error: ", err)
			}
		}
	}

	log.Log.Println("web api want send to :", onlineMem)
	resp := &pb.StreamResponse{
		DataType:  IM_MSG_FROM_UPLOAD_OR_WS_OR_APP,
		ImMsgData: req,
		Res:       &pb.Result{Code: http.StatusOK, Msg: "receiver im message successful"},
	}
	// 发送在线用户消息
	if onlineMem != nil {
		dc.Task <- *NewImTask(req.Id, onlineMem, resp)
		log.Log.Printf("dispatcher finish %+v <-||||-> %+v", req.Id, resp)
	}
}

// Dispatcher 根据数据类型，调用不同的函数，产生不同的数据，然后把数据放到发送chan中
func pushDataDispatcher(dc *DataContext, ds DataSource) {
	for {
		log.Log.Infof("dispatcher client msg")
		errMap := &sync.Map{}
		srv := ds.(pb.TalkCloud_DataPublishServer)

		data, _ := srv.Recv()

		if data == nil {
			log.Log.Println("this stream is no data, maybe is offline")
			return
		}

		// 如果再次登录的用户在map中已存在，并且srv不同，那么就给dc的chan里面写一个终止信号
		if v := StreamMap.Get(data.Uid); v != nil && v != srv {
			log.Log.Printf("this here %+v, %+v", v, srv)
			log.Log.Printf("this user # %d is login already", data.Uid)
			dc.ExceptionalLogin <- data.Uid
			return
		}

		if v := StreamMap.Get(data.Uid); v == nil {
			// 更新stream和redis状态
			StreamMap.Set(data.Uid, srv)
			log.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> login data.Uid :", data.Uid)
			if err := tuc.UpdateOnlineInCache(&pb.Member{Id: data.Uid, Online: USER_ONLINE}, cache.GetRedisClient()); err != nil {
				log.Log.Println("Update user online state error:", err)
			}
		}

		log.Log.Printf(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> start gen data for dispatcher")
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

			//// 如果连接存在，就3s往channel里面放一个数据
		//	go sendHeartbeat(dc, data, cfgGs.Interval, srv)

			// 往dc里面写上线通知
			go notifyToOther(dc.Task, data.Uid, LOGIN_NOTIFY_MSG)

		case IM_MSG_FROM_UPLOAD_OR_WS_OR_APP:
			imMessagePublishDispatcher(dc, data.ImMsg)

		case KEEP_ALIVE_MSG:
			if err := tuc.UpdateOnlineInCache(&pb.Member{Id: data.Uid, Online: USER_ONLINE}, cache.GetRedisClient()); err != nil {
				log.Log.Println("KEEP_ALIVE_MSG Update user online state error:", err)
			}
		}
	}
}

// temp  临时用一下，后期统统会改用RabbitMQ
func dispatcherScheduler(dContent *DataContext, multiSend bool) {
	log.Log.Printf("start Scheduler im msg")
	var notify int32
	tick := time.Tick(time.Minute * 5)
	for {
		// 接收任务
		select {
		case t := <-dContent.Task:
			go func() { TQ.Tasks <- t }()
			//log.Log.Printf("///////////%T///////%+v", t, t)
			if t.Data.DataType == LOGOUT_NOTIFY_MSG {
				notify++
				if notify == t.Data.NotifyTotal {
					log.Log.Printf("notify: %d, total: %d", notify, t.Data.NotifyTotal)
					return
				}
			}
			if !multiSend {
				log.Log.Printf("only scheduler once")
				return
			}
		case <-tick:
			log.Log.Printf("single im task queue len = %d", len(dContent.UId)) //TODO 合理退出，关闭调度器
		}
	}
}

// Executor 推送登录返回数据、IM离线数据、IM离线数据、Heartbeat
func pushDataExecutor(ct chan Task) {
	var wg sync.WaitGroup
	for {
		select {
		case task := <-ct:
			log.Log.Println("global executor receiver: ", task.Receiver)
			for _, receiverId := range task.Receiver {
				wg.Add(1)
				go pushData(task, receiverId, task.Data, &wg)
				log.Log.Println("will send to ", receiverId)
			}
			wg.Wait()
		}
	}
}

// 推送数据
func pushData(task Task, receiverId int32, resp *pb.StreamResponse, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Log.Printf("the stream map have: %+v", StreamMap)
	if value := StreamMap.Get(receiverId); value != nil {
		srv := value.(pb.TalkCloud_DataPublishServer)
		log.Log.Printf("# %d receiver response: %+v", receiverId, resp)
		if err := srv.Send(resp); err != nil {
			// 发送失败处理
			processErrorSendMsg(err, task, receiverId, resp)
		} else {
			// 发送成功如果是离线数据（接收者等于stream id自己）就去更新状态
			log.Log.Printf("send success. dc.senderId: %d, receiverId: %d", task.SenderId, receiverId)
			if task.SenderId == receiverId && resp.DataType == OFFLINE_IM_MSG {
				//  更新数据库里面的消息的状态
				if err := tm.SetMsgStat(receiverId, READ_OFFLINE_IM_MSG, db.DBHandler); err != nil {
					log.Log.Println("Add offline msg with error: ", err)
				}
			}
		}
	} else {
		log.Log.Errorf("Send to %d  im that can't find stream", receiverId) //TODO 就依靠那边心跳了，这里就不管发送失败了
		// 存储即时发送失败的消息
		if resp.DataType == IM_MSG_FROM_UPLOAD_OR_WS_OR_APP && task.SenderId != receiverId {
			// 把发送数据保存进数据库, 如果是离线数据就忽略
			log.Log.Printf("send fail. dc.senderId: %d, receiverId: %d", task.SenderId, receiverId)
			if err := tm.AddMsg(resp.ImMsgData, db.DBHandler); err != nil {
				log.Log.Errorf("Send fail and add offline msg with error: ", err)
			}
		}
	}
}

// 处理推送数据失败的情况
func processErrorSendMsg(err error, task Task, receiverId int32, resp *pb.StreamResponse) {
	log.Log.Debugf("send msg fail with error: ", err)

	// 判断错误类型
	if errSC, _ := status.FromError(err); errSC.Code() == codes.Unavailable || errSC.Code() == codes.Canceled {

		//log.Log.Printf("%+v ---- %+v start add offline msg: ", resp, task)
		if resp.DataType == IM_MSG_FROM_UPLOAD_OR_WS_OR_APP && task.SenderId != receiverId {
			// 把发送数据保存进数据库, 如果是离线数据就忽略
			log.Log.Infof("send fail. dc.senderId: %d, receiverId: %d", task.SenderId, receiverId)
			if err := tm.AddMsg(resp.ImMsgData, db.DBHandler); err != nil {
				log.Log.Errorf("Send fail and add offline msg with error: ", err)
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
		log.Log.Printf("App login error : %s", err)
		return &pb.StreamResponse{
			Res: &pb.Result{Code: 500, Msg: "User Login Process failed. Please try again later"},
		}, err
	}

	if err == sql.ErrNoRows {
		log.Log.Printf("App login error : %s", err)
		return &pb.StreamResponse{
			Res: &pb.Result{Code: 500, Msg: "User is not exist error. Please try again later"},
		}, err
	}

	if res.PassWord != data.Passwd {
		log.Log.Printf("App login error : %s", err)
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
	if v := StreamMap.Get(int32(res.Id)); v != srv {
		log.Log.Printf("this here %+v, %+v", v, srv)
		dc.ExceptionalLogin <- int32(res.Id)
		log.Log.Println("this user is login already")
		return &pb.StreamResponse{
			Res: &pb.Result{Code: 500, Msg: "this user is login already"},
		}, errors.New("the user is login already")
	}
	// 更新stream和redis状态
	log.Log.Println("login data.Uid:", int32(res.Id))
	//StreamMap.Store(int32(res.Id), srv)

	StreamMap.Set(int32(res.Id), srv)

	if err := tuc.UpdateOnlineInCache(&pb.Member{Id: int32(res.Id), Online: USER_ONLINE}, cache.GetRedisClient()); err != nil {
		log.Log.Println("Update user online state error:", err)
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
			log.Log.Println(k, " gen error: ", err)
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
		log.Log.Printf("session id is error%s", err)
		errMap.Store("processSession", err)
		return
	}

	ct := time.Now().UnixNano() / 1000000
	ttl := ct + 30*60*1000 // Severside session valid time: 30 min
	ttlStr := strconv.FormatInt(ttl, 10)

	sInfo := &model.SessionInfo{SessionID: sessionId, UserName: req.Name, UserPwd: req.Passwd, TTL: ttlStr}
	if err := s.InsertSession(sInfo); err != nil {
		log.Log.Printf("session id insert is error%s", err)
		errMap.Store("processSession", err)
		return
	}
}

// 获取好友列表
func getFriendList(uid int32, fList chan *pb.FriendsRsp, errMap *sync.Map, wg *sync.WaitGroup) {
	log.Log.Println("get FriendList start")
	var err error
	fl, _, err := tuf.GetFriendReqList(int32(uid), db.DBHandler)
	if err != nil {
		errMap.Store("getFriendList", err)
		fList <- nil
	} else {
		fList <- fl
	}
	log.Log.Println("get FriendList done")
}

// 获取群组列表
func getGroupList(uid int32, gList chan *pb.GroupListRsp, errMap *sync.Map, wg *sync.WaitGroup) {
	log.Log.Debugf("Get group list start")
	// 先去缓存取，取不出来再去mysql取
	gl, err := tuc.GetGroupListFromRedis(int32(uid), cache.GetRedisClient())
	if err != nil && err != sql.ErrNoRows {
		log.Log.Println("No find In CacheError")
		errMap.Store("getGroupList", err)
		log.Log.Printf("get GroupList%v", err)
		gList <- gl
		return
	}

	if err == sql.ErrNoRows {
		log.Log.Println("redis is not find")
		for {
			gl, _, err = tg.GetGroupListFromDB(int32(uid), db.DBHandler)
			if err != nil {
				errMap.Store("getGroupList", err)
				break
			}
			log.Log.Println("start update redis GetGroupListFromDB")
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
						log.Log.Println("Add user information to cache with error: ", err)
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
	log.Log.Debugf("Get group list done")
}

// 增加缓存
func addUserInfoToCache(userInfo *pb.Member, wg *sync.WaitGroup) {

	log.Log.Printf("will get rediscli, now redis pool info :%+v |<<<>>>| idleCount: %+v", cache.RedisPool, cache.RedisPool.IdleCount())
	redisCli := cache.GetRedisClient()
	if err := tuc.AddUserDataInCache(userInfo, redisCli); err != nil {
		log.Log.Println("Add user information to cache with error: ", err)
	}
	log.Log.Println("addUserInfoToCache done")
}

// 返回的IM离线数据
func GetOfflineImMsgFromDB(req *pb.StreamRequest) (*pb.StreamResponse, error) {
	// 去数据库拉取离线数据
	log.Log.Println("start get offline im msg")
	offlineMsg, err := tm.GetMsg(req.Uid, UNREAD_OFFLINE_IM_MSG, db.DBHandler)
	if err != nil {
		log.Log.Println("Get offline msg fail with error:", err)
		return nil, err
	}
	log.Log.Printf("get offline msg %+v", offlineMsg)

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
				log.Log.Printf("v %d, s %v msg %+v", v, ok, msg)
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

	log.Log.Printf("%+v \n %+v", respPkgSingle, respPkgGroup)

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
			log.Log.Debugf("# %d receiver response: %+v", data.Uid, resp)
			//if value, ok := StreamMap.Load(data.Uid); ok {
			if value := StreamMap.Get(data.Uid); value != nil {
				srv := value.(pb.TalkCloud_DataPublishServer)
				if err := srv.Send(resp); err != nil {
					//if errSC, _ := status.FromError(err); errSC.Code() == codes.Unavailable || errSC.Code() == codes.Canceled {
					// 只要是发送失败，就认为对方离线
					log.Log.Infof("client %d close with %+v", data.Uid, err)
					log.Log.Infof("now dc stream : %+v", StreamMap)

					// 删除map中的stream
					StreamMap.Del(data.Uid)
					log.Log.Infof("# user %d is logout successfully", data.Uid)

					// 更新redis状态
					if err := tuc.UpdateOnlineInCache(&pb.Member{Id: data.Uid, Online: USER_OFFLINE}, cache.GetRedisClient()); err != nil {
						log.Log.Println("Update user online state error:", err)
					}
					// 往dc里面写掉线通知
					go notifyToOther(dc.Task, data.Uid, LOGOUT_NOTIFY_MSG)
					return
					//}
				}
			}
		}

	}

}

// 上线通知所有人，掉线通知所有人、sos通知
func notifyToOther(dcTask chan Task, uId int32, notifyType int32) {
	var (
		errMap      = &sync.Map{}
		selfGList   = make(chan *pb.GroupListRsp, 1)
		notifyTotal = int32(0)
		notifyId    = make([]int32,0)
	)
	log.Log.Printf("notify root id :%d", uId)
	uInfo, _ := tuc.GetUserFromCache(uId)
	_, uLocation, _ := tlc.GetUserLocationInCache(uId, cache.GetRedisClient())

	getGroupList(uId, selfGList, errMap, nil)
	gl := <-selfGList
	if gl != nil && uInfo != nil {
		for _, g := range gl.GroupList {
			for _, u := range g.UsrList {
				if u.Uid != uId && u.Online == tuc.USER_ONLINE {
					//log.Log.Printf("will notify *******************------------------------------------------------%d", u.Uid)
					notifyTotal++
					notifyId = append(notifyId, u.Uid)
				}
			}
		}
		// add self or not add self
		notifyTotal++
		notifyId = append(notifyId, uId)

		log.Log.Printf("notify total: %d and notify all id: %+v", notifyTotal, notifyId)
		for _, g := range gl.GroupList {
			for _, u := range g.UsrList {
				if u.Uid != uId && u.Online == tuc.USER_ONLINE {
					doNotify(u, errMap, notifyType,
						notifyTotal,
						uInfo,
						uLocation,
						dcTask,
						uId)

				}
			}
		}
		// send to self
		doNotify(uInfo, errMap, notifyType, notifyTotal, uInfo, uLocation, dcTask, uId)
	}
	log.Log.Printf("cant load notify self info %v----------%v", gl, uInfo)
}

func doNotify(u *pb.UserRecord, errMap *sync.Map, notifyType int32, notifyTotal int32,
	uInfo *pb.UserRecord, uLocation *pb.GPS, dcTask chan Task, uId int32) {
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
	log.Log.Printf("will send %d notify to %+v", notifyType, recvId)
	dcTask <- *NewImTask(uId, recvId, resp)
	//}
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
