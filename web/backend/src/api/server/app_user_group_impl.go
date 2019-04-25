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
	"net/http"
	"pkg/group"
	tg "pkg/group"        // table group
	tgc "pkg/group_cache" // table group cache
	tuc "pkg/user_cache"
	"strconv"
	"strings"
)

const (
	CREATE_GROUP_BY_DISPATCHER = 1
	CREATE_GROUP_BY_USER = 0

)
func (serv *TalkCloudServiceImpl) CreateGroup(ctx context.Context, req *pb.CreateGroupReq) (*pb.CreateGroupResp, error) {
	deviceIds := make([]int, 0)
	if req.DeviceIds == "-1" {
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
		userType = CREATE_GROUP_BY_DISPATCHER // 管理员创建
	} else {
		userType = CREATE_GROUP_BY_USER // 普通app用户创建
	}

	log.Printf("Create group is name: %s", gl.GroupInfo.GroupName)
	if gid, err := tg.CreateGroup(gl, userType);
		err != nil {
		log.Println("create group error :", err)
		return &pb.CreateGroupResp{Res: &pb.Result{Msg: "create group unsuccessful, please try again later", Code: 422}}, err
	} else {
		gl.GroupInfo.Id = int(gid)
	}

	// 增加到缓存  TODO 出错就清空缓存算了
	if err := tgc.AddGroupAndUserInCache(gl, cache.GetRedisClient()); err != nil {
		log.Printf("CreateGroup AddGroupAndUserInCache error: %v", err)
	}

	// 增加所创建群所含成员也要加进缓存,因为每个成员都新加了一个群组,还要把每个人的信息也加入缓存
	for _, v := range gl.DeviceIds {
		if err := tgc.AddGroupSingleMemCache(int32(gl.GroupInfo.Id), int32(v), cache.GetRedisClient()); err != nil {
			log.Printf("CreateGroup AddGroupAndUserInCache error: %v", err)
		}

		if err := tuc.AddUserForSingleGroupCache(int32(v), int32(gl.GroupInfo.Id), cache.GetRedisClient()); err != nil {
			log.Println("CreateGroup add group member into single group into cache error:", err)
		}

		u := &pb.UserRecord{}
		tuc.UpdateUserFromDBToRedis(u, v)
	}

	return &pb.CreateGroupResp{
		GroupInfo: &pb.GroupInfo{Gid: int32(gl.GroupInfo.Id), GroupName: gl.GroupInfo.GroupName,},
		Res:       &pb.Result{Msg: "create group successful.", Code: 200},
	}, nil
}

func (serv *TalkCloudServiceImpl) InviteUserIntoGroup(ctx context.Context, req *pb.InviteUserReq) (*pb.InviteUserResp, error) {
	uIdStrs := strings.Split(req.Uids, ",")
	uIds := make([]int32, 0)
	resp := &pb.InviteUserResp{
		Res: &pb.Result{
			Msg:  "Invite user into group unsuccessful, please try again later",
			Code: http.StatusInternalServerError,
		},
	}
	for _, v := range uIdStrs {
		uId, err := strconv.Atoi(v)
		if err != nil {
			log.Printf("Invite user into group range uIdStrs have error: %v", err)
			return resp, nil
		}
		uIds = append(uIds, int32(uId))
	}

	log.Printf("%v", uIds)
	for _, v := range uIds {
		err := group.AddGroupMember(v, req.Gid, group.GROUP_MEMBER, db.DBHandler)
		if err != nil {
			log.Printf("Invite user into group range uIds have error: %v", err)
			return resp, nil
		}
	}
	// 添加进缓存
	// 1. 更新用户的group那个set
	if err := tuc.AddUsersGroupInCache(uIds, req.Gid, cache.GetRedisClient()); err != nil {
		log.Println("Invite user into group AddUsersGroupInCache error: ", err)
		return resp, nil
	}
	// 2. 更新群组里有哪些用户那个set AddGroupSingleMemCache
	if err := tgc.AddGroupMemsInCache(req.Gid, uIds, cache.GetRedisClient()); err != nil {
		log.Println("Invite user into group AddGroupMemsInCache error: ", err)
		return resp, nil
	}

	resp.Res.Code = http.StatusOK
	resp.Res.Msg = "Invite user into group successful"
	gMem, err := tuc.GetGroupMemDataFromCache(req.Gid, cache.GetRedisClient())
	if err != nil {
		log.Println("Invite user into group GetGroupMemDataFromCache error: ", err)
		return resp, nil

	}
	resp.UsrList = gMem
	resp.Res.Msg = "Invite user into group successful"
	resp.Res.Code = http.StatusOK
	return resp, err

}

// 获取群组成员信息
func (serv *TalkCloudServiceImpl) GetGroupInfo(ctx context.Context, req *pb.GetGroupInfoReq) (*pb.GetGroupInfoResp, error) {
	// 直接去缓存获取了 TODO
	res, err := tgc.GetGroupInfoFromCache(req.Gid, cache.GetRedisClient())
	if err != nil {
		log.Println("GetGroupInfo has error:", err)
		return &pb.GetGroupInfoResp{Res: &pb.Result{Msg: "Get group info unsuccessful, please try again later", Code: 500}}, err
	}
	return &pb.GetGroupInfoResp{
		Res: &pb.Result{
			Msg:  "Get group info successful",
			Code: http.StatusOK,
		},
		GroupInfo: res,
	}, nil
}

func (serv *TalkCloudServiceImpl) JoinGroup(ctx context.Context, req *pb.GrpUserAddReq) (*pb.GrpUserAddRsp, error) {
	// 如果已经在群组里，就直接返回
	resp := &pb.GrpUserAddRsp{Res: &pb.Result{Msg: "Join group unsuccessful, please try again later", Code: 500}}
	_, gMap, err := group.GetGroupListFromDB(req.Uid, db.DBHandler)
	if err != nil {
		log.Printf("JoinGroup GetGroupListFromDB error: %+v", err)
		return resp, err
	}
	if _, ok := (*gMap)[req.Gid]; ok {
		log.Println("User join this group already")
		return resp, err
	}

	// TODO 判断要id是不是有没有权限加群?
	err = group.AddGroupMember(req.Uid, req.Gid, group.GROUP_MEMBER, db.DBHandler)
	if err != nil {
		log.Printf("JoinGroup AddGroupMember error: %+v", err)
		return resp, err
	}
	// 添加进缓存
	// 1. 更新用户的group那个set
	if err := tuc.AddUserForSingleGroupCache(req.Uid, req.Gid, cache.GetRedisClient()); err != nil {
		log.Println("JoinGroup AddUserForSingleGroupCache error: ", err)
		return resp, err
	}
	// 2. 更新群组里有哪些用户那个set AddGroupSingleMemCache
	if err := tgc.AddGroupSingleMemCache(req.Gid, req.Uid, cache.GetRedisClient()); err != nil {
		log.Println("JoinGroup AddGroupSingleMemCache error: ", err)
		return resp, err
	}
	// 3. 添加这个群组的信息进缓存，因为这个是模糊搜索的结果
	gInfo, err := tg.GetGroupInfoFromDB(req.Gid, req.Uid)
	if err != nil {
		log.Println("JoinGroup GetGroupInfoFromDB error: ", err)
		return resp, err
	}

	//3.1 每个用户的信息
	for _, u := range gInfo.UsrList {
		if err := tuc.AddUserDataInCache(&pb.Member{
			Id:          u.Uid,
			IMei:        u.Imei,
			UserName:    u.Name,
			Online:      u.Online,
			LockGroupId: u.LockGroupId,
		}, cache.GetRedisClient()); err != nil {
			log.Println("Add user information to cache with error: ", err)
		}
	}

	//3.2  每一个群组拥有的成员
	if err := tgc.AddGroupCache(gInfo.UsrList, gInfo, cache.GetRedisClient()); err != nil {
		log.Println("JoinGroup AddGroupCache error: ", err)
		return resp, err
	}

	return &pb.GrpUserAddRsp{Res: &pb.Result{Msg: "Join group successful", Code: 200}}, err
}

func (serv *TalkCloudServiceImpl) GetGroupList(ctx context.Context, req *pb.GrpListReq) (*pb.GroupListRsp, error) {
	// 先去缓存取，取不出来再去mysql取
	log.Println("Get GroupList start")
	// 先去缓存取，取不出来再去mysql取
	gl, err := tuc.GetGroupListFromRedis(int32(req.Uid), cache.GetRedisClient())
	if err != nil && err != sql.ErrNoRows {
		log.Println("No find In CacheError")
		log.Printf("get GroupList%v", err)
		return &pb.GroupListRsp{Res: &pb.Result{Code: 500, Msg: "process error, please try again"}}, err
	}

	// TODO 有一个隐藏问题，redis如果只有一部分数据
	if err == sql.ErrNoRows {
		log.Println("redis is not find， start to mysql query")
		for {
			gl, _, err = tg.GetGroupListFromDB(int32(req.Uid), db.DBHandler)
			if err != nil {
				log.Printf("get GroupList %+v", err)
				return &pb.GroupListRsp{Res: &pb.Result{Code: 500, Msg: "process error, please try again"}}, err
			}
			log.Println("start update redis")
			// 新增到缓存 更新两个地方，首先，每个组的信息要更新，就是group data，记录了群组的id和名字
			if err := tgc.AddGroupInCache(gl, cache.GetRedisClient()); err != nil {
				log.Printf("get GroupList %+v", err)
				return &pb.GroupListRsp{Res: &pb.Result{Code: 500, Msg: "process error, please try again"}}, err
			}

			// 其次更新一个userSet  就是一个组里有哪些用户
			if err := tuc.AddUserInGroupToCache(gl, cache.GetRedisClient()); err != nil {
				log.Printf("get GroupList %+v", err)
				return &pb.GroupListRsp{Res: &pb.Result{Code: 500, Msg: "process error, please try again"}}, err
			}

			// 每个用户的信息
			for _, g := range gl.GroupList {
				for _, u := range g.UsrList {
					if err := tuc.AddUserDataInCache(&pb.Member{
						Id:          u.Uid,
						IMei:        u.Imei,
						UserName:    u.Name,
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
					return &pb.GroupListRsp{Res: &pb.Result{Code: 500, Msg: "process error, please try again"}}, err
				}
			}
			break
		}
	}

	res := &pb.GroupListRsp{
		Res:       &pb.Result{Msg: "get group list successful", Code: 200},
		GroupList: gl.GroupList,
	}
	return res, err
}

// 通过关键字返回群组，区分在群组和不在的群组
func (serv *TalkCloudServiceImpl) SearchGroup(ctx context.Context, req *pb.GrpSearchReq) (*pb.GroupListRsp, error) {
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
	_, gMap, err := group.GetGroupListFromDB(req.Uid, db.DBHandler)
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

func (serv *TalkCloudServiceImpl) RemoveGrpUser(ctx context.Context, req *pb.GrpUserDelReq) (*pb.GrpUserDelRsp, error) {
	err := group.RemoveGroupMember(req.Uid, req.Gid, db.DBHandler)
	resp := &pb.GrpUserDelRsp{
		Res: &pb.Result{
			Msg:  "remove Group User error, please try again later.",
			Code: http.StatusInternalServerError,
		},
	}
	if err != nil {
		log.Println("Remove Group error: ", err)
		return resp, nil
	}

	// 清空缓存
	// 1. 更新该用户在哪些组的那个set
	if err := tuc.RemoveUserForSingleGroupCache(req.Uid, req.Gid, cache.GetRedisClient()); err != nil {
		log.Println("JoinGroup AddUserForSingleGroupCache error: ", err)
	}
	// 2. 更新群组里有哪些用户那个set AddGroupSingleMemCache
	if err := tgc.RemoveGroupSingleMemCache(req.Uid, req.Gid, cache.GetRedisClient()); err != nil {
		log.Println("JoinGroup AddUserForSingleGroupCache error: ", err)
	}
	return resp, err
}

func (serv *TalkCloudServiceImpl) ExitGrp(ctx context.Context, req *pb.UserExitReq) (*pb.UserExitRsp, error) {
	return nil, nil
}

func (serv *TalkCloudServiceImpl) RemoveGrp(ctx context.Context, req *pb.GroupDelReq) (*pb.GroupDelRsp, error) {
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
