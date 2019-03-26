package api

import (
	"api/talk_cloud"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"math/rand"
	"net"
	pb "server/web/backend/src/api/talk_cloud"
	"server/web/backend/src/model"
	d "server/web/backend/src/pkg/device"
	s "server/web/backend/src/pkg/session"
	"server/web/backend/src/utils"
	"strconv"
	"time"
)

const (
	SERVICEPORT = "9000"
)

type TalkCloudService struct{}

// 注册App
func (serv *TalkCloudService) AppRegister(ctx context.Context, req *pb.AppRegReq) (*pb.AppRegRsp, error) {
	iMei := strconv.FormatInt(int64(rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(1000000000000000)), 10)
	appRegResp := &pb.AppRegRsp{}

	// 查重
	user, err := d.GetAdviceByName(req.Name)
	if err != nil {
		log.Printf("app register error : %s", err)
		appRegResp.Err = &pb.ErrorMsg{
			Code: 500,
			Msg:  "User registration failed. Please try again later",
		}
		return appRegResp, nil
	}
	if user > 0 {
		appRegResp.Err = &pb.ErrorMsg{
			Code: 500,
			Msg:  "User name has been registered",
		}
		return appRegResp, nil
	}

	device := &model.Device{
		UserName:  req.Name,
		PassWord:  req.Password,
		AccountId: 1, // TODO 默认给谁
		IMei:      iMei,
	}

	if err := d.AddDevice(device); err != nil {
		log.Printf("app register error : %s", err)
		appRegResp.Err = &pb.ErrorMsg{
			Code: 500,
			Msg:  "User registration failed. Please try again later",
		}
		return appRegResp, nil
	}

	res, err := d.SelectDeviceByKey(req.Name)
	if err != nil {
		log.Printf("app register error : %s", err)
		appRegResp.Err = &pb.ErrorMsg{
			Code: 500,
			Msg:  "User registration Process failed. Please try again later",
		}
		return appRegResp, nil
	}

	return &pb.AppRegRsp{Id: int64(res.Id)}, nil
}

// 设备注册
func (serv *TalkCloudService) DeviceRegister(ctx context.Context, req *pb.DeviceRegReq) (*pb.DeviceRegRsp, error) {
	// TODO 设备串号和账户id进行校验

	name := string([]byte(req.DeviceList)[9:len(req.DeviceList)])
	device := &model.Device{
		UserName:  name,
		PassWord:  "123456",
		AccountId: int(req.AccountId),
		IMei:      req.DeviceList,
	}

	if err := d.AddDevice(device); err != nil {
		log.Printf("app register error : %s", err)
		return &pb.DeviceRegRsp{Err: &pb.ErrorMsg{Code: 500, Msg: "Device registration failed. Please try again later"}}, err
	}

	return &pb.DeviceRegRsp{Err: &pb.ErrorMsg{Code: 200, Msg: "Device registration successful"}}, nil
}

func (serv *TalkCloudService) Login(ctx context.Context, req *pb.LoginReq) (*talk_cloud.LoginRsp, error) {
	//　验证用户名是否存在以及密码是否正确，然后就生成一个uuid session, 把sessionid放进metadata返回给客户端，
	//  然后之后的每一次连接都需要客户端加入这个metadata，使用拦截器，对用户进行鉴权
	if req.Name == "" || req.Passwd == "" {
		return &pb.LoginRsp{Err: &pb.ErrorMsg{Code: 422, Msg: "用户名或密码不能为空"}}, nil
	}

	sessionId, err := utils.NewUUID();
	if err != nil {
		log.Panicf("session id is error%s", err)
		return &pb.LoginRsp{Err: &pb.ErrorMsg{Code: 500, Msg: "server internal error"}}, err
	}

	ct := time.Now().UnixNano() / 1000000
	ttl := ct + 30*60*1000 // Severside session valid time: 30 min
	ttlStr := strconv.FormatInt(ttl, 10)

	sInfo := &model.SessionInfo{SessionID: sessionId, UserName: req.Name, UserPwd: req.Passwd, TTL: ttlStr}
	if err := s.InsertSession(sInfo); err != nil {
		log.Panicf("session id insert is error%s", err)
		return &pb.LoginRsp{Err: &pb.ErrorMsg{Code: 500, Msg: "server internal error"}}, err
	}

	// create and send header
	header := metadata.Pairs("sessionId", sessionId)
	if err := grpc.SendHeader(ctx, header); err != nil {
		log.Panicf("sessionid metadata set  error%s", err)
		return &pb.LoginRsp{Err: &pb.ErrorMsg{Code: 500, Msg: "server internal error"}}, err
	}
	return &pb.LoginRsp{Err: &pb.ErrorMsg{Code: 200, Msg: req.Name + "login successful"}}, nil
}

func (serv *TalkCloudService) Logout(ctx context.Context, req *pb.LogoutReq) (*talk_cloud.LogoutRsp, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Panicf("sessionid metadata set  error%s")
		return &pb.LogoutRsp{Err: &pb.ErrorMsg{Code: 403, Msg: "server internal error"}}, nil
	}
	// TODO 考虑要不要验证sessionInfo中的name和password
	if err := s.DeleteSession(md.Get("sessionId")[1]); err != nil {
		log.Panicf("sessionid metadata delete  error%s", err)
		return &pb.LogoutRsp{Err: &pb.ErrorMsg{Code: 500, Msg: "server internal error"}}, err
	}
	return &pb.LogoutRsp{Err: &pb.ErrorMsg{Code: 200, Msg: req.Name + "logout successful"}}, nil
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

// Test
func maint() {
	groupServer := grpc.NewServer()
	talk_cloud.RegisterTalkCloudServer(groupServer, &TalkCloudService{})

	lis, err := net.Listen("tcp", ":"+SERVICEPORT)
	if err != nil {
		log.Fatalf("group net listen err: %v", err)
	}

	if err := groupServer.Serve(lis); err != nil {
		log.Fatalf("监听失败")
	} else {
		log.Println("listing")
	}
}
