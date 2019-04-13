/*
@Time : 2019/4/12 19:29 
@Author : yanKoo
@File : talk_cloud_app_login_impl
@Software: GoLand
@Description:
*/
package server

import (
	sl "api/server/server_login"
	pb "api/talk_cloud"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type TalkCloudService struct {
}

// 推送初次登录应该有的群组，等信息 TODO 后期可能可以改变为使用im那个
func (tcs *TalkCloudService) LoginProcess(srv pb.TalkCloud_LoginProcessServer) error {
	return status.Errorf(codes.Unimplemented, "method LoginProcess not implemented")
}

func (tcs *TalkCloudService) KeepAlive(srv pb.TalkCloud_KeepAliveServer) error {
	return status.Errorf(codes.Unimplemented, "method KeepAlive not implemented")
}
func (tcs *TalkCloudService) ImMessagePublish(srv pb.TalkCloud_ImMessagePublishServer) error {
	return status.Errorf(codes.Unimplemented, "method ImMessagePublish not implemented")
}

func (tcs *TalkCloudService) DataPublish(srv pb.TalkCloud_DataPublishServer) error {
	// 读数据，运行一个runner，根据数据不同类型分发给不同的executor
	//go func() {
	//	for {
	data, _ := srv.Recv()
	log.Println(data)
	r := sl.NewRunner(30, false, srv,
		sl.ClientMsgDispatcher,
		sl.FirstLoginDataExecutor,
		sl.ImMsgDataDataExecutor,
		sl.KeepAliveExecutor)
	r.StartAll()
	//}
	//}()

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
