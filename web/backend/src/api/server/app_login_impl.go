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
	tm "pkg/msg"
	s "pkg/session"
	tu "pkg/user"
	tuf "pkg/user_friend"
	"server/common/src/cache"
	"server/common/src/db"
	"strconv"
	"sync"
	"time"
	"utils"
)

var StreamMap =  sync.Map{}


type TalkCloudService struct {
}

type DataContent struct {
	DataChan       chan interface{}
	Receiver       []int32
	senderId       int32
	KeepAliveClose chan int
	StreamMap      *sync.Map
}

type DataSource interface{}

// 推送初次登录应该有的群组，等信息 TODO 后期可能可以改变为使用im那个
func (tcs *TalkCloudService) LoginProcess(srv pb.TalkCloud_LoginProcessServer) error {
	return status.Errorf(codes.Unimplemented, "method LoginProcess not implemented")
}

func (tcs *TalkCloudService) KeepAlive(srv pb.TalkCloud_KeepAliveServer) error {
	return status.Errorf(codes.Unimplemented, "method KeepAlive not implemented")
}

func (tcs *TalkCloudService) ImMessagePublish(ctx context.Context, req *pb.ImMsgReqData) (*pb.ImMsgRespData, error) {
	var (
		errMap   sync.Map
		executor = make(chan interface{}, 1)
	)

	var data = DataContent{
		DataChan:       make(chan interface{}, 30),
		KeepAliveClose: make(chan int, 1),
		StreamMap:      &StreamMap,
	}

	go func(data *DataContent, ds DataSource) {
		if err := ImMessagePublishDispatcher(data, req); err != nil {
			errMap.Store(<-data.DataChan, err)
		}
	}(&data, req)

	go func(data *DataContent, e chan interface{}) {
		log.Println("start executor")
		if err := pushDataExecutor(data, executor); err != nil {
			errMap.Store(<-data.DataChan, err)
		}
	}(&data, executor)

	done := <-executor
	log.Println("ImMessagePublish ", done)
	return &pb.ImMsgRespData{Result: &pb.Result{Msg: "push data done", Code: 200}}, nil
}

// 获取
func ImMessagePublishDispatcher(dc *DataContent, ds DataSource) error {
	// 获取要发送的数据
	req := ds.(*pb.ImMsgReqData)
	// 获取在线离线用户id
	offlineMem := make([]int32, 0)
	onlineMem := make([]int32, 0)
	log.Println("ImMessagePublishDispatcher start")
	log.Printf("ImMessagePublishDispatcher now /*dc.*/StreamMap map have: %+v， %p", dc.StreamMap, &dc.StreamMap)
	log.Println("ImMessagePublishDispatcher end")

	if req.ReceiverType == 0 { // 发给单人
		// 判断是否在线
		v, ok := dc.StreamMap.Load(req.ReceiverId)
		log.Println(req.ReceiverId, v, ok)
		log.Printf("now dc.StreamMap map have: %+v， %p", dc.StreamMap, &dc.StreamMap)
		log.Println(req.ReceiverId)
		if !ok {
			log.Println("不在线")
			// 保存进数据库
			if err := tm.AddMsg(req, db.DBHandler); err != nil {
				log.Println("Add offline msg with error: ", err)
				return err
			}
		}
		onlineMem = append(onlineMem, req.ReceiverId)
	}

	if req.ReceiverType == 1 { // 发送给群组成员，然后区分离线在线
		res, err := tg.GetGroupMem(req.ReceiverId, cache.GetRedisClient())
		if err != nil {
			log.Println("Add offline msg get group member with error: ", err)
			return err
		}

		for _, v := range res {
			if _, ok := dc.StreamMap.Load(v); ok {
				onlineMem = append(onlineMem, int32(v))
			} else {
				offlineMem = append(offlineMem, int32(v))
			}
		}

		// 存储离线消息
		if err := tm.AddMultiMsg(req, offlineMem, db.DBHandler); err != nil {
			log.Println("Add offline msg with error: ", err)
			return err
		}
	}

	log.Println("web api want send to :", onlineMem)
	// 发送在线用户消息
	dc.senderId = req.Id
	dc.Receiver = onlineMem
	dc.DataChan <- &pb.StreamResponse{
		DataType:  2,
		ImMsgData: req,
		Res:       &pb.Result{Code: http.StatusOK, Msg: "receiver im message successful"},
	}
	log.Printf("dispatcher finish %+v |****| %+v", dc.Receiver, dc.DataChan)
	return nil
}

func (tcs *TalkCloudService) DataPublish(srv pb.TalkCloud_DataPublishServer) error {
	var errMap sync.Map
	var data = DataContent{
		DataChan:       make(chan interface{}, 30),
		KeepAliveClose: make(chan int, 1),
		StreamMap:      &StreamMap,
	}

	go func(data *DataContent, ds DataSource) {

		if err := pushDataDispatcher(data, srv); err != nil {
			errMap.Store(<-data.DataChan, err)
		}
		log.Println()
	}(&data, srv)

	go func(data *DataContent, ds DataSource) {

		if err := pushDataExecutor(data); err != nil {
			errMap.Store(<-data.DataChan, err)
		}
	}(&data, srv)

	// TODO 错误处理

	c := make(chan int)
	<-c
	return nil
}

//用户登录
func (tcs *TalkCloudService) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRsp, error) {
	//// 登录就开启任务四个协程，一个用来去返回初步登录应该获取的数据，一个用来推送发送保活数据，一个用来发送im数据
	//
	//// 1. TODO 用sync.map将上下文保存? 记录在线用户
	//
	//// 2.goroutine，建立一个客户端，调用loginProess方法
	//
	//// 3.goroutine，建立一个客户端，调用KeepAlive方法
	//
	//// 4.goroutine，建立一个客户端，调用ImMessagePublish方法
	////loginProscess(ctx context.Context, req *pb.LoginReq) (*pb.LoginRsp, error)
	//log.Println("enter login")
	////time.Sleep(time.Second*20)
	////　验证用户名是否存在以及密码是否正确，然后就生成一个uuid session, 把sessionid放进metadata返回给客户端，
	////  然后之后的每一次连接都需要客户端加入这个metadata，使用拦截器，对用户进行鉴权
	//if req.Name == "" || req.Passwd == "" {
	//	return &pb.LoginRsp{Res: &pb.Result{Code: 422, Msg: "用户名或密码不能为空"}}, nil
	//}
	//
	//res, err := tu.SelectUserByKey(req.Name)
	//if err != nil && err != sql.ErrNoRows {
	//	log.Printf("App login error : %s", err)
	//	loginRsp := &pb.LoginRsp{
	//		Res: &pb.Result{
	//			Code: 500,
	//			Msg:  "User Login Process failed. Please try again later"},
	//	}
	//	return loginRsp, nil
	//}
	//
	//if err == sql.ErrNoRows {
	//	log.Printf("App login error : %s", err)
	//	loginRsp := &pb.LoginRsp{
	//		Res: &pb.Result{
	//			Code: 500,
	//			Msg:  "User is not exist error. Please try again later"},
	//	}
	//	return loginRsp, nil
	//}
	//
	//if res.PassWord != req.Passwd {
	//	log.Printf("App login error : %s", err)
	//	loginRsp := &pb.LoginRsp{
	//		Res: &pb.Result{
	//			Code: 500,
	//			Msg:  "User Login pwd error. Please try again later"},
	//	}
	//	return loginRsp, nil
	//}
	//userInfo := &pb.Member{
	//	Id:          int32(res.Id),
	//	IMei:        res.IMei,
	//	UserName:    res.UserName,
	//	NickName:    res.NickName,
	//	UserType:    int32(res.UserType),
	//	LockGroupId: int32(res.LockGroupId),
	//	Online:      1, // 登录就在线
	//}
	//
	//var (
	//	errMap   = &sync.Map{}
	//	wg       sync.WaitGroup
	//	fList    = make(chan *pb.FriendsRsp, 1)
	//	gList    = make(chan *pb.GroupListRsp, 1)
	//	existErr bool
	//)
	//defer func() {
	//	close(fList)
	//	close(gList)
	//}()
	//wg.Add(4) // 获取用户列表 用户群组列表，加入缓存, 插入session
	//
	//// 1. 处理登录session
	//go processSession(&ctx, req, errMap, &wg)
	//
	//// 2. 获取好友列表
	//go getFriendList(int32(res.Id), fList, errMap, &wg)
	//
	//// 3. 群组列表
	//go getGroupList(int32(res.Id), gList, errMap, &wg)
	//
	//// 4. 将用户信息添加进redis
	//go addUserInfoToCache(userInfo, &wg)
	//
	//wg.Wait()
	//
	////遍历该map，参数是个函数，该函数参的两个参数是遍历获得的key和value，返回一个bool值，当返回false时，遍历立刻结束。
	//errMap.Range(func(k, v interface{}) bool {
	//	err = v.(error)
	//	if err != nil {
	//		log.Println(k, " gen error: ", err)
	//		existErr = true
	//		return false
	//	}
	//	return true
	//})
	//if existErr {
	//	return &pb.LoginRsp{Res: &pb.Result{Code: 500, Msg: "process error, please try again"}}, nil
	//}
	//
	//loginRep := &pb.LoginRsp{
	//	UserInfo:   userInfo,
	//	FriendList: (<-fList).FriendList,
	//	GroupList:  (<-gList).GroupList,
	//	Res:        &pb.Result{Code: 200, Msg: req.Name + " login successful"},
	//}

	//return loginRep, nil
	return nil, nil
}

// 根据数据类型，调用不同的函数，产生不同的数据，然后把数据放到发送chan中
func pushDataDispatcher(dc *DataContent, ds DataSource) error {
	interval := 2
	for {
		log.Println("dispatcher client msg")
		errMap := &sync.Map{}
		srv := ds.(pb.TalkCloud_DataPublishServer)

		data, _ := srv.Recv()

		if data == nil {
			return errors.New("this stream is no data")
		}

		// 更新stream和redis状态
		dc.StreamMap.Store(data.Uid, srv)
		if err := tu.UpdateOnlineInCache(&pb.Member{Id: data.Uid, Online: 1}, cache.GetRedisClient()); err != nil {
			log.Println("Update user online state error:", err)
		}

		switch data.DataType {
		case 1:
			res, err := firstLoginData(data) // 去数据库拉取数据考虑放在这里
			if err != nil {
				errMap.Store(data.Uid, err)
			}
			dc.DataChan <- res

		case 2:
			res, err := ImMsgData(data)
			if err != nil {
				errMap.Store(data.Uid, err)
			}
			for _, v := range res {
				dc.DataChan <- v
			}
		}

		re := make([]int32, 0)
		dc.senderId = data.Uid
		dc.Receiver = append(re, data.Uid) // TODO

		// 如果连接存在，就3s往channel里面放一个数据
		go sendHeartbeat(dc, data, interval)
	}
}

func sendHeartbeat(dc *DataContent, data *pb.StreamRequest, interval int) {
	// 使用定时器，收到登录请求之后每隔5s后发送一次数据
	timerTask := time.NewTicker(time.Second * time.Duration(interval))
timerLoop:
	for {
		select {
		case <-timerTask.C:
		loop:
			for {
				select {
				case uid := <-dc.KeepAliveClose:
					log.Printf("client %d close", uid)
					log.Println(dc.StreamMap.Load(uid))
					break timerLoop
				default:
					dc.senderId = data.Uid
					dc.DataChan <- keepAliveData(data)
					break loop
				}
			}
		}
	}

}

func pushDataExecutor(dContent *DataContent, option ...interface{}) error {
	var wg sync.WaitGroup

	for response := range dContent.DataChan {
		resp := response.(*pb.StreamResponse)
		// 读取接收者
		log.Println("receiver: ", dContent.Receiver)
		for _, receiverId := range dContent.Receiver {
			log.Println("will send to ", receiverId)
			wg.Add(1)
			go pushMessage(dContent, receiverId, resp, &wg)
		}
		wg.Wait()

		for _, v := range option {
			v.(chan interface{}) <- "done,ha ha ha"
		}
	}
	return nil
}

func pushMessage(dc *DataContent, receiverId int32, resp *pb.StreamResponse, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("the stream map have: %+v", dc.StreamMap)
	if value, ok := dc.StreamMap.Load(receiverId); ok {
		srv := value.(pb.TalkCloud_DataPublishServer)
		log.Printf("%d server send response: %+v", receiverId, resp)
		if err := srv.Send(resp); err != nil {
			// 发送失败处理
			processErrorSendMsg(err, dc, receiverId, resp)
		} else {
			// 发送成功如果是离线数据（接收者等于stream id自己）就去更新状态
			log.Printf("send success. dc.senderId: %d, receiverId: %d",dc.senderId, receiverId)
			if dc.senderId == receiverId {
				//  更新数据库里面的消息的状态
				msgStat := 1
				if err := tm.SetMsgStat(receiverId, msgStat,  db.DBHandler); err != nil {
					log.Println("Add offline msg with error: ", err)
				}
			}
		}
	} else {
		log.Println("can't find stream")
	}
}

// 处理发送数据失败的情况
func processErrorSendMsg(err error, dc *DataContent, receiverId int32, resp *pb.StreamResponse){
	log.Println("send msg fail with error: ", err)
	// 判断错误类型
	errSC, _ := status.FromError(err)
	log.Println(errSC.Code())
	if errSC.Code() == codes.Unavailable {
		// 1. 只要是发送失败，就认为对方离线，如果是发送Im数据，就保存到数据库
		dc.KeepAliveClose <- int(receiverId)

		// 删除map中的stream
		dc.StreamMap.Delete(receiverId)
		//log.Printf("now dc.StreamMap map have: %+v， %p", dc.StreamMap, &dc.StreamMap)

		// 更新redis状态
		if err := tu.UpdateOnlineInCache(&pb.Member{Id: receiverId, Online: 0}, cache.GetRedisClient()); err != nil {
			log.Println("Update user online state error:", err)
		}

		if resp.DataType == 2 && dc.senderId != receiverId {
			// 把发送数据保存进数据库 TODO 如果是离线数据就忽略
			// TODO 有个问题，如果发送的时候，还在线，函数走到这里，就掉线了，
			//  但是发送方收到的回复仍然是已发送，发送方能看见的只有该用户已离线的状态
			log.Printf("send success. dc.senderId: %d, receiverId: %d",dc.senderId, receiverId)
			if err := tm.AddMsg(resp.ImMsgData, db.DBHandler); err != nil {
				log.Println("Add offline msg with error: ", err)
			}
		}
	}
}

// TODO 返回请求信息
func firstLoginData(req *pb.StreamRequest) (*pb.StreamResponse, error) {

	//time.Sleep(time.Second*20)
	//　验证用户名是否存在以及密码是否正确，然后就生成一个uuid session, 把sessionid放进metadata返回给客户端，
	//  然后之后的每一次连接都需要客户端加入这个metadata，使用拦截器，对用户进行鉴权
	if req.Name == "" || req.Passwd == "" {
		return &pb.StreamResponse{
			Res: &pb.Result{Code: 422, Msg: "用户名或密码不能为空"},
		}, errors.New("username or password can't be empty")
	}

	res, err := tu.SelectUserByKey(req.Name)
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

	if res.PassWord != req.Passwd {
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
		Online:      1, // 登录就在线
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
	wg.Add(4) // 获取用户列表 用户群组列表，加入缓存, 插入session

	// 1. 处理登录session
	go processSession(req, errMap, &wg)

	// 2. 获取好友列表
	go getFriendList(int32(res.Id), fList, errMap, &wg)

	// 3. 群组列表
	go getGroupList(int32(res.Id), gList, errMap, &wg)

	// 4. 将用户信息添加进redis
	go addUserInfoToCache(userInfo, &wg)

	wg.Wait()

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
			Res: &pb.Result{Code: 500, Msg: "User Login pwd error. Please try again later"},
		}, err
	}
	return &pb.StreamResponse{
		LoginResp: &pb.FirstLoginData{
			UserInfo:   userInfo,
			FriendList: (<-fList).FriendList,
			GroupList:  (<-gList).GroupList,
		},
		Res: &pb.Result{Code: 200, Msg: req.Name + " login successful"},
	}, nil
}

// TODO 给stream模式加metadata
func processSession(req *pb.StreamRequest, errMap *sync.Map, wg *sync.WaitGroup) {
	sessionId, err := utils.NewUUID()
	if err != nil {
		log.Panicf("session id is error%s", err)
		errMap.Store("processSession", err)
		wg.Done()
		return
	}

	ct := time.Now().UnixNano() / 1000000
	ttl := ct + 30*60*1000 // Severside session valid time: 30 min
	ttlStr := strconv.FormatInt(ttl, 10)

	sInfo := &model.SessionInfo{SessionID: sessionId, UserName: req.Name, UserPwd: req.Passwd, TTL: ttlStr}
	if err := s.InsertSession(sInfo); err != nil {
		log.Panicf("session id insert is error%s", err)
		errMap.Store("processSession", err)
		wg.Done()
		return
	}

	/*// create and send header
	header := metadata.Pairs("sessionId", sessionId)
	if err := grpc.SendHeader(*ctx, header); err != nil {
		log.Panicf("sessionid metadata set  error%s", err)
		errMap.Store("processSession", err)
		wg.Done()
		return
	}*/

	wg.Done()
}

// 获取好友列表
func getFriendList(uid int32, fList chan *pb.FriendsRsp, errMap *sync.Map, wg *sync.WaitGroup) {
	var err error
	fl, _, err := tuf.GetFriendReqList(int32(uid), db.DBHandler)
	if err != nil {
		errMap.Store("getFriendList", err)
		fList <- nil
	} else {
		fList <- fl
	}
	wg.Done()
}

// 获取群组列表
func getGroupList(uid int32, gList chan *pb.GroupListRsp, errMap *sync.Map, wg *sync.WaitGroup) {
	// 先去缓存取，取不出来再去mysql取
	gl, err := tu.GetGroupList(int32(uid), cache.GetRedisClient())
	if err != nil && err != sql.ErrNoRows {
		log.Println("cache.Nofind In CacheError")
		errMap.Store("getGroupList", err)
		gList <- gl
		wg.Done()
		return
	}

	// TODO 有一个隐藏问题，redis如果只有一部分数据
	if err == sql.ErrNoRows {
		log.Println("redis is not find")
		for {
			gl, _, err = tg.GetGroupList(int32(uid), db.DBHandler)
			if err == nil {
				err := errors.New("hello errors")
				errMap.Store("getGroupList", err)
				break
			}
			// 新增到缓存 更新两个地方，首先，每个组的信息要更新，就是group data，记录了群组的id和名字 TODO 后期应该要把群组里有哪些人也在这里查出来，更新。
			if err := tg.AddGroupInCache(gl, cache.GetRedisClient()); err != nil {
				errMap.Store("getGroupList", err)
				break
			}
			// 其次更新一个userSet
			if err := tu.AddUserInGroupToCache(gl, cache.GetRedisClient()); err != nil {
				errMap.Store("getGroupList", err)
				break
			}
			break
		}
	}

	gList <- gl
	wg.Done()
}

// 增加缓存 TODO 没有错误
func addUserInfoToCache(userInfo *pb.Member, wg *sync.WaitGroup) {
	if err := tu.AddUserDataInCache(userInfo, cache.GetRedisClient()); err != nil {
		log.Println("Add user information to cache with error: ", err)
	}

	wg.Done()
}

// 刚登陆的时候根据请求数据，生成相应的response
func ImMsgData(req *pb.StreamRequest) ([]*pb.StreamResponse, error) {
	// 去数据库拉取离线数据
	offlineMsg, err := tm.GetMsg(req.Uid, int32(0), db.DBHandler)
	if err != nil {
		log.Println("Get offline msg fail with error:", err)
		return nil, err
	}
	// 发送组和个人的消息过去
	offMsgResp := make([]*pb.StreamResponse, 0)

	// TODO 看app需要怎么接，是全部人给他，还是单个数据单个数据过去。多半就是使用单个数据
	for _, v := range offlineMsg {
		offMsgResp = append(offMsgResp, &pb.StreamResponse{
			ImMsgData: v,
		})
	}

	return offMsgResp, nil
}

// 生成保活请求
func keepAliveData(req *pb.StreamRequest) *pb.StreamResponse {
	return &pb.StreamResponse{
		KeepAlive: &pb.KeepAlive{
			Uid: req.Uid,
			SYN: 1,
		},
	}
}
