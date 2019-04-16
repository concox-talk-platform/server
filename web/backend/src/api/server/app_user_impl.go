package server

import (
	pb "api/talk_cloud"
	"cache"
	"context"
	"db"
	"errors"
	"google.golang.org/grpc/metadata"
	"log"
	"math/rand"
	"model"
	s "pkg/session"
	tu "pkg/user"
	"strconv"
	"time"
	"utils"
)


// 注册App
func (tcs *TalkCloudService) AppRegister(ctx context.Context, req *pb.AppRegReq) (*pb.AppRegRsp, error) {
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

	return &pb.AppRegRsp{Id: int32(res.Id), UserName: req.Name, Res: &pb.Result{Code: 200, Msg: "User registration successful"}}, nil
}

// 设备注册
func (tcs *TalkCloudService) DeviceRegister(ctx context.Context, req *pb.DeviceRegReq) (*pb.DeviceRegRsp, error) {
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

// 用户注销
func (tcs *TalkCloudService) Logout(ctx context.Context, req *pb.LogoutReq) (*pb.LogoutRsp, error) {
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

// 设置用户所在默认锁定组 TODO 对每一个用户设置操作都做鉴权
func (tcs *TalkCloudService) SetLockGroupId(ctx context.Context, req *pb.SetLockGroupIdReq) (*pb.SetLockGroupIdResp, error) {
	if !utils.CheckId(int(req.UId)) || !utils.CheckId(int(req.GId)) {
		err := errors.New("uid or gid is not valid")
		log.Println("service SetLockGroupId error :", err)
		return &pb.SetLockGroupIdResp{Res: &pb.Result{Msg: "User id or group id is not valid, please try again later.", Code: 500}}, nil
	}

	// TODO 用户id是否在组所传id中

	if err := tu.SetLockGroupId(req, db.DBHandler); err != nil {
		log.Println("service SetLockGroupId error :", err)
		return &pb.SetLockGroupIdResp{Res: &pb.Result{Msg: "Set lock default group error, please try again later.", Code: 500}}, nil
	}

	// 更新缓存
	if err := tu.UpdateLockGroupIdInCache(req, cache.GetRedisClient()); err != nil {
		log.Println("service SetLockGroupId error :", err)
		// TODO 去把数据库里的群组恢复？
		return &pb.SetLockGroupIdResp{Res: &pb.Result{Msg: "Set lock default group error, please try again later.", Code: 500}}, nil
	}

	return &pb.SetLockGroupIdResp{Res: &pb.Result{Msg: "Set lock default group success", Code: 200}}, nil
}
