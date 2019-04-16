package server

import (
	pb "api/talk_cloud"
	"cache"
	"context"
	"database/sql"
	"db"
	"errors"
	"log"
	"model"
	"pkg/group"
	tg "pkg/group" // table group
	tu "pkg/user"
	"strconv"
	"strings"
)

func (serv *TalkCloudService) CreateGroup(ctx context.Context, req *pb.CreateGroupReq) (*pb.CreateGroupResp, error) {
	deviceIds := make([]int, 0)
	if req.DeviceIds == "-1"{
		deviceIds = append(deviceIds, -1)
		req.GroupName = req.GroupInfo.GroupName
		req.AccountId = req.GroupInfo.AccountId
	} else {
		for _, v := range strings.Split(req.DeviceIds, ",") {
			id, _ := strconv.Atoi(v)
			deviceIds = append(deviceIds, id)
		}
	}
	log.Println("create group member id : ", deviceIds)
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
		GroupName: req.GroupName,
		AccountId: int(req.AccountId),
		Status:    strconv.FormatInt(int64(64), 10),
	}

	log.Println("req.GroupName:", req.GroupName)

	gl := &model.GroupList{DeviceIds: deviceIds, DeviceInfo: deviceInfos, GroupInfo: groupInfo}
	var userType = -1
	if deviceIds[0] < 0 {
		userType = 1 // 管理员创建
	} else {
		userType = 0 // 普通app用户创建
	}

	log.Printf("%+v", gl.GroupInfo.GroupName)
	if gid, err := tg.CreateGroup(gl, userType);
		err != nil {
		log.Println("create group error :", err)
		return &pb.CreateGroupResp{Res: &pb.Result{Msg: "create group unsuccessful, please try again later", Code: 422}}, err
	} else {
		gl.GroupInfo.Id = int(gid)
	}

	// 增加到缓存
	if err := tg.AddGroupAndUserInCache(gl, cache.GetRedisClient()); err != nil {
		log.Println("insert cache error")
		return &pb.CreateGroupResp{Res: &pb.Result{Msg: "create group unsuccessful, please try again later", Code: 422}}, err
	}

	return &pb.CreateGroupResp{
		GroupInfo: &pb.GroupInfo{Gid: int32(gl.GroupInfo.Id), GroupName: gl.GroupInfo.GroupName,},
		Res:       &pb.Result{Msg: "create group successful.", Code: 200},
	}, nil
}

func (serv *TalkCloudService) JoinGroup(ctx context.Context, req *pb.GrpUserAddReq) (*pb.GrpUserAddRsp, error) {
	// 如果已经在群组里，就直接返回
	_, gMap, err := group.GetGroupList(req.Uid, db.DBHandler)
	if err != nil {
		return &pb.GrpUserAddRsp{Res: &pb.Result{Msg: "Join group unsuccessful, please try again later", Code: 422}}, err
	}
	if _, ok := (*gMap)[req.Gid]; ok {
		log.Println("User join this group already")
		return &pb.GrpUserAddRsp{Res: &pb.Result{Code: 422, Msg: "User join this group already"}}, nil
	}

	// TODO 判断要id是不是有没有权限加群
	err = group.AddGroupMember(int64(req.Uid), int64(req.Gid), group.GROUP_NORMAL_USER, db.DBHandler)
	if err != nil {
		return &pb.GrpUserAddRsp{Res: &pb.Result{Msg: "Join group unsuccessful, please try again later", Code: 422}}, err
	}

	return &pb.GrpUserAddRsp{Res: &pb.Result{Msg: "Join group successful", Code: 200}}, err
}

func (serv *TalkCloudService) GetGroupList(ctx context.Context, req *pb.GrpListReq) (*pb.GroupListRsp, error) {
	// 先去缓存取，取不出来再去mysql取 TODO redis取出一部分，但是没有更新
	log.Println("Get GroupList start")
	res, err := tu.GetGroupList(req.Uid, cache.GetRedisClient())
	if err != nil && err != sql.ErrNoRows {
		log.Println("cache.NofindInCacheError")
		return &pb.GroupListRsp{Res: &pb.Result{Code: 500, Msg: "process error, please try again"}}, err
	}
	if err == sql.ErrNoRows {
		log.Println("get")
		res, _, err = group.GetGroupList(req.Uid, db.DBHandler)
		if err != nil {
			return &pb.GroupListRsp{Res: &pb.Result{Code: 500, Msg: "process error, please try again"}}, err
		}
		// 新增到缓存 更新两个地方，首先，每个组的信息要更新，就是group data，记录了群组的id和名字 TODO 后期应该要把群组里有哪些人也在这里查出来，更新。
		if err := tg.AddGroupInCache(res, cache.GetRedisClient()); err != nil {
			return &pb.GroupListRsp{Res: &pb.Result{Code: 500, Msg: "process error, please try again"}}, err
		}
		// 其次更新一个userSet
		if err := tu.AddUserInGroupToCache(res, cache.GetRedisClient()); err != nil {
			return &pb.GroupListRsp{Res: &pb.Result{Code: 500, Msg: "process error, please try again"}}, err
		}
	}
	res.Res = &pb.Result{Msg: "get group list successful", Code: 200}
	return res, err
}

// 通过关键字返回群组，区分在群组和不在的群组
func (serv *TalkCloudService) SearchGroup(ctx context.Context, req *pb.GrpSearchReq) (*pb.GroupListRsp, error) {
	// 判空
	if req.Target == "" {
		return &pb.GroupListRsp{Res: &pb.Result{Code: 422, Msg: "process error, please input target"}}, errors.New("target is nil")
	}

	// 模糊查询群组 TODO 暂时这么写吧，感觉有点蠢
	groups, err := group.SearchGroup(req.Target, db.DBHandler)
	if err != nil {
		return &pb.GroupListRsp{Res: &pb.Result{Code: 500, Msg: "process error, please try again"}}, err
	}
	// 查找用户所在组
	_, gMap, err := group.GetGroupList(req.Uid, db.DBHandler)
	if err != nil {
		return &pb.GroupListRsp{Res: &pb.Result{Code: 500, Msg: "process error, please try again"}}, err
	}

	for _, v := range groups.GroupList {
		if _, ok := (*gMap)[v.Gid]; ok {
			v.IsExist = true
		}
	}
	log.Printf("server search group: %+v", groups)
	groups.Res = &pb.Result{Msg: "search group success", Code: 200}
	return groups, nil
}

func (serv *TalkCloudService) RemoveGrpUser(ctx context.Context, req *pb.GrpUserDelReq) (*pb.GrpUserDelRsp, error) {
	err := group.RemoveGroupMember(req.Uid, req.Gid, db.DBHandler)
	rsp := new(pb.GrpUserDelRsp)
	rsp.Res = new(pb.Result)
	rsp.Res.Code = 0
	rsp.Res.Msg = ""
	if err != nil {
		rsp.Res.Code = -1
		rsp.Res.Msg = ""
	}

	return rsp, err
}

func (serv *TalkCloudService) ExitGrp(ctx context.Context, req *pb.UserExitReq) (*pb.UserExitRsp, error) {
	return nil, nil
}

func (serv *TalkCloudService) RemoveGrp(ctx context.Context, req *pb.GroupDelReq) (*pb.GroupDelRsp, error) {
	// clear group user first
	err := group.ClearGroupMember(req.Gid, db.DBHandler)
	rsp := new(pb.GroupDelRsp)
	rsp.Res = new(pb.Result)

	if err != nil {
		rsp.Res.Code = -1
		rsp.Res.Msg = err.Error()
	}

	// then remove group
	err = group.RemoveGroup(req.Gid, db.DBHandler)
	if err != nil {
		rsp.Res.Code = -2
		rsp.Res.Msg = err.Error()
	}

	return rsp, err
}
