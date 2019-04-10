package server

import (
	"api/talk_cloud"
	pb "api/talk_cloud"
	"cache"
	"context"
	"database/sql"
	"db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"math/rand"
	"model"
	s "pkg/session"
	tu "pkg/user"
	"pkg/user_friend"
	"server/web/backend/src/pkg/group"
	"strconv"
	"time"
	"utils"
)

type TalkCloudService struct{}

// 注册App
func (serv *TalkCloudService) AppRegister(ctx context.Context, req *pb.AppRegReq) (*pb.AppRegRsp, error) {
	iMei := strconv.FormatInt(int64(rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(1000000000000000)), 10)
	appRegResp := &pb.AppRegRsp{}

	// 查重
	ifExist, err := tu.GetUserByName(req.Name)
	if err != nil {
		log.Printf("app register error : %s", err)
		appRegResp.Res = &pb.Result{
			Code: 500,
			Msg:  "User registration failed. Please try again later",
		}
		return appRegResp, nil
	}
	if ifExist > 0 {
		appRegResp.Res = &pb.Result{
			Code: 500,
			Msg:  "User name has been registered",
		}
		return appRegResp, nil
	}

	user := &model.User{
		UserName:  req.Name,
		PassWord:  req.Password,
		AccountId: -1, //  app用户是-1 调度员是0，设备用户是受管理的账号id
		IMei:      iMei,
		UserType:  1,
	}

	log.Println("app register start")
	if err := tu.AddUser(user); err != nil {
		log.Printf("app register error : %s", err)
		appRegResp.Res = &pb.Result{
			Code: 500,
			Msg:  "User registration failed. Please try again later",
		}
		return appRegResp, nil
	}

	res, err := tu.SelectUserByKey(req.Name)
	if err != nil {
		log.Printf("app register error : %s", err)
		appRegResp.Res = &pb.Result{
			Code: 500,
			Msg:  "User registration Process failed. Please try again later",
		}
		return appRegResp, nil
	}

	return &pb.AppRegRsp{Id: int32(res.Id), UserName:req.Name, Res:&pb.Result{Code:200, Msg:"User registration successful"}}, nil
}

// 设备注册
func (serv *TalkCloudService) DeviceRegister(ctx context.Context, req *pb.DeviceRegReq) (*pb.DeviceRegRsp, error) {
	// TODO 设备串号和账户id进行校验
	name := string([]byte(req.DeviceList)[9:len(req.DeviceList)])
	user := &model.User{
		UserName:  name,
		PassWord:  "123456",
		AccountId: int(req.AccountId),
		IMei:      req.DeviceList,
	}

	if err := tu.AddUser(user); err != nil {
		log.Printf("app register error : %s", err)
		return &pb.DeviceRegRsp{Res: &pb.Result{Code: 500, Msg: "Device registration failed. Please try again later"}}, err
	}

	return &pb.DeviceRegRsp{Res: &pb.Result{Code: 200, Msg: "Device registration successful"}}, nil
}

// 用户登录
func (serv *TalkCloudService) Login(ctx context.Context, req *pb.LoginReq) (*talk_cloud.LoginRsp, error) {
	log.Println("enter login")
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

	sessionId, err := utils.NewUUID()
	if err != nil {
		log.Panicf("session id is error%s", err)
		return &pb.LoginRsp{Res: &pb.Result{Code: 500, Msg: "server internal error"}}, err
	}

	ct := time.Now().UnixNano() / 1000000
	ttl := ct + 30*60*1000 // Severside session valid time: 30 min
	ttlStr := strconv.FormatInt(ttl, 10)

	sInfo := &model.SessionInfo{SessionID: sessionId, UserName: req.Name, UserPwd: req.Passwd, TTL: ttlStr}
	if err := s.InsertSession(sInfo); err != nil {
		log.Panicf("session id insert is error%s", err)
		return &pb.LoginRsp{Res: &pb.Result{Code: 500, Msg: "server internal error"}}, err
	}

	// create and send header
	header := metadata.Pairs("sessionId", sessionId)
	if err := grpc.SendHeader(ctx, header); err != nil {
		log.Panicf("sessionid metadata set  error%s", err)
		return &pb.LoginRsp{Res: &pb.Result{Code: 500, Msg: "server internal error"}}, err
	}

	userInfo := &pb.Member{
		Id:       int32(res.Id),
		IMei:     res.IMei,
		UserName: res.UserName,
		NickName: res.NickName,
		UserType: int32(res.UserType),
	}

	// 好友列表
	fList, _, err := user_friend.GetFriendReqList(int32(res.Id), db.DBHandler)
	if err != nil {
		return &pb.LoginRsp{Res: &pb.Result{Code: 500, Msg: "process error, please try again"}}, err
	}

	// 群组列表
	// 先去缓存取，取不出来再去mysql取
	gList, err := tu.GetGroupList(int32(res.Id), cache.GetRedisClient())
	if err != nil && err != sql.ErrNoRows {
		log.Println("cache.NofindInCacheError")
		return &pb.LoginRsp{Res: &pb.Result{Code: 500, Msg: "process error, please try again"}}, err
	}
	if err == sql.ErrNoRows {
		log.Println("get")
		gList, _, err = group.GetGroupList(int32(res.Id), db.DBHandler)
		if err != nil {
			return &pb.LoginRsp{Res: &pb.Result{Code: 500, Msg: "process error, please try again"}}, err
		}
	}

	return &pb.LoginRsp{
		UserInfo:   userInfo,
		FriendList: fList.FriendList,
		GroupList: gList.GroupList,
		Res:      &pb.Result{Code: 200, Msg: req.Name + " login successful"},
	}, nil
}

// 用户注销 TODO
func (serv *TalkCloudService) Logout(ctx context.Context, req *pb.LogoutReq) (*talk_cloud.LogoutRsp, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Panicf("sessionid metadata set  error%s")
		return &pb.LogoutRsp{Res: &pb.Result{Code: 403, Msg: "server internal error"}}, nil
	}
	// TODO 考虑要不要验证sessionInfo中的name和password
	if err := s.DeleteSession(md.Get("sessionId")[1], nil); err != nil {
		log.Panicf("sessionid metadata delete  error%s", err)
		return &pb.LogoutRsp{Res: &pb.Result{Code: 500, Msg: "server internal error"}}, err
	}
	return &pb.LogoutRsp{Res: &pb.Result{Code: 200, Msg: req.Name + "logout successful"}}, nil
}

//// authenticateClient check the client credentials
//func authenticateClient(ctx context.Context) (string, error) {
//	if md, ok := metadata.FromIncomingContext(ctx); ok {
//		clientUsername := strings.Join(md["username"], "")
//		//if clientUsername != "valineliu" {
//		//	return "", fmt.Errorf("unknown user %s", clientUsername)
//		//}
//		clientSessionId := strings.Join(md["sessionId"], "")
//		sInfo, err := pkg.GetSessionValue(clientSessionId)
//		if err != nil {
//			log.Printf("authenticated client: %s", clientUsername)
//			return "", fmt.Errorf("missing credentials")
//		}
//		log.Printf("authenticated client: %s", clientUsername)
//		return "9527", nil
//	}
//	return "", fmt.Errorf("missing credentials")
//}
//
//// unaryInterceptor calls authenticateClient with current context
//func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
//	clientID, err := authenticateClient(ctx)
//	if err != nil {
//		return nil, err
//	}
//	ctx = context.WithValue(ctx, "clientID", clientID)
//	return handler(ctx, req)
//}
