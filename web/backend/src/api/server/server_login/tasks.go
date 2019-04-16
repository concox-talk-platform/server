package server_login

import (
	pb "api/talk_cloud"
	"database/sql"
	"errors"
	"log"
	"model"
	tg "pkg/group"
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

func ClientMsgDispatcher(dc DataChan, dt *int, srv pb.TalkCloud_DataPublishServer) error {
	log.Println("dispatcher client msg")
	if data, err := srv.Recv(); err != nil {
		log.Println("ClientMsgDispatcher fail with error:", err)
		return err
	} else {
		switch data.DataType {
		case 1:
			dc <- &pb.LoginReq{Name: data.Name, Passwd: data.Passwd}
		case 2:
			dc <- &pb.ImMsgReqData{}
		case 3:
			dc <-&pb.KeepAliveReq{Uid:data.Uid, ACK:data.ACK}
		}
		dType := int(data.DataType)
		dt = &dType
	}
	//log.Println("dc:", <-dc, <-dt)
	return nil
}

// TODO 返回请求信息
func FirstLoginDataExecutor(dc DataChan, srv pb.TalkCloud_DataPublishServer) error {
	streamReq := <- dc
	req := streamReq.(*pb.LoginReq)
	//time.Sleep(time.Second*20)
	//　验证用户名是否存在以及密码是否正确，然后就生成一个uuid session, 把sessionid放进metadata返回给客户端，
	//  然后之后的每一次连接都需要客户端加入这个metadata，使用拦截器，对用户进行鉴权

	/*
	if (req.Name == "" || req.Passwd == "" {

		return &pb.LoginRsp{Res: &pb.Result{Code: 422, Msg: "用户名或密码不能为空"}}, nil
	}*/

	res, err := tu.SelectUserByKey(req.Name)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("App login error : %s", err)
		_ = srv.Send(&pb.StreamResponse{
			Res: &pb.Result{Code: 500, Msg: "User Login Process failed. Please try again later"},
		})
		return err
	}

	if err == sql.ErrNoRows {
		log.Printf("App login error : %s", err)
		_ = srv.Send(&pb.StreamResponse{
			Res: &pb.Result{Code: 500, Msg: "User is not exist error. Please try again later"},
		})
		return err
	}

	if res.PassWord != req.Passwd {
		log.Printf("App login error : %s", err)
		_ = srv.Send(&pb.StreamResponse{
			Res: &pb.Result{Code: 500, Msg: "User Login pwd error. Please try again later"},
		})
		return err
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
	go processSession(srv, req, errMap, &wg)

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
		_ = srv.Send(&pb.StreamResponse{
			Res: &pb.Result{Code: 500, Msg: "User Login pwd error. Please try again later"},
		})
		return err
	}
	_ = srv.Send(&pb.StreamResponse{
		LoginResp: &pb.FirstLoginData{
			UserInfo:   userInfo,
			FriendList: (<-fList).FriendList,
			GroupList:  (<-gList).GroupList,
		},
		Res: &pb.Result{Code: 200, Msg: req.Name + " login successful"},
	})
	return nil
}

// TODO 给stream模式加metadata
func processSession(_ pb.TalkCloud_DataPublishServer, req *pb.LoginReq, errMap *sync.Map, wg *sync.WaitGroup) {
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

func ImMsgDataDataExecutor(dc DataChan, srv pb.TalkCloud_DataPublishServer) error {
	return nil
}

// TODO 处理保活请求
func KeepAliveExecutor(dc DataChan, srv pb.TalkCloud_DataPublishServer) error {
	keepAliveReq := (<-dc).(*pb.KeepAliveReq)
	log.Println(keepAliveReq.Uid, keepAliveReq.ACK)
	_ = srv.Send(&pb.StreamResponse{
		KeepAlive: &pb.KeepAlive{
			Uid: keepAliveReq.Uid,
			SYN: 1,
		},
	})
	return nil
}