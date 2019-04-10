/*
@Time : 2019/3/28 15:33
@Author : yanKoo
@File : TalkCloudRegisterImpl
@Software: GoLand
@Description: 目前主要供web端调用 protoc -I . talk_cloud_web.proto --go_out=plugins=grpc:.
*/
package server

import (
	pb "api/talk_cloud"
	"cache"
	"context"
	"log"
	"model"
	tg "pkg/group"
	"strconv"
)

type WebServiceServerImpl struct {
}

// 更新组成员
func (wssu *WebServiceServerImpl) UpdateGroup(ctx context.Context, req *pb.UpdateGroupReq) (*pb.UpdateGroupResp, error) {
	log.Println("enter update group")
	deviceIds := make([]int, 0)
	for _, v := range req.DeviceIds {
		deviceIds = append(deviceIds, int(v))
	}

	deviceInfos := make([]interface{}, 0)
	for _, v := range req.DeviceInfos {
		log.Println("web impl vid:", v.Id)
		deviceInfos = append(deviceInfos, map[string]interface{}{
			"id":         int(v.Id),
			"imei":       v.IMei,
			"user_name":  v.UserName,
			"nick_name":  v.NickName,
			"password":   v.Pwd,
			"user_type":  int(v.UserType),
			"account_id": int(v.AccountId),
			"parent_id":  v.ParentId,
		})
	}

	groupInfo := &model.GroupInfo{
		Id:        int(req.GroupInfo.Id),
		GroupName: req.GroupInfo.GroupName,
		AccountId: int(req.GroupInfo.AccountId),
		Status:    strconv.FormatInt(int64(64), 10),
	}

	gl := &model.GroupList{DeviceIds: deviceIds, DeviceInfo: deviceInfos, GroupInfo: groupInfo}
	var userType = -1
	if deviceIds[0] < 0 {
		userType = 1 // 管理员创建
	} else {
		userType = 0 // 普通用户创建
	}
	if gid, err := tg.UpdateGroupMember(gl, userType);
		err != nil {
		log.Println("create group error :", err)
		return &pb.UpdateGroupResp{ResultMsg: &pb.Result{Msg: "Update group unsuccessful, please try again later", Code: 422}}, err
	} else {
		gl.GroupInfo.Id = int(gid)
	}

	// TODO 增加到缓存
	if err := tg.AddGroupAndUserInCache(gl, cache.GetRedisClient()); err != nil {
		log.Println("insert cache error")
		return &pb.UpdateGroupResp{ResultMsg: &pb.Result{Msg: "Update group unsuccessful, please try again later", Code: 422}}, err
	}

	return &pb.UpdateGroupResp{ResultMsg: &pb.Result{Msg: "Update group successful.", Code: 200}}, nil
}

// 删除组
func (wssu *WebServiceServerImpl) DeleteGroup(ctx context.Context, req *pb.Group) (*pb.DeleteGroupResp, error) {
	if err := tg.DeleteGroup(&model.GroupInfo{Id: int(req.Id)}); err != nil {
		log.Printf("delete group fail , error: %s", err)
		return &pb.DeleteGroupResp{
			ResultMsg: &pb.Result{
				Msg:  "delete group error ",
				Code: 500}}, err
	}

	// TODO 更新缓存
	return &pb.DeleteGroupResp{ResultMsg: &pb.Result{Msg: "delete group success ", Code: 200}}, nil
}
