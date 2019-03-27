package api

import (
	"context"
	"db"
	"pkg/group"
	pb "api/talk_cloud"
)

func (serv *TalkCloudService) CreateGroup(ctx context.Context, req *pb.GroupNewReq) (*pb.GroupNewRsp, error) {
	err := group.CreateGroup(int64(req.Uid), req.GroupName, db.DBHandler)
	rsp := new(pb.GroupNewRsp)
	rsp.Err = new(pb.ErrorMsg)
	rsp.Err.Code = 0
	rsp.Err.Msg = ""

	if err != nil {
		rsp.Err.Code = -1
		rsp.Err.Msg = err.Error()
	}

	return rsp, err
}

func (serv *TalkCloudService) JoinGroup(ctx context.Context, req *pb.GrpUserAddReq) (*pb.GrpUserAddRsp, error) {
	err := group.AddGroupUser(int64(req.Uid), int64(req.Gid), group.GROUP_NORMAL_USER, db.DBHandler)

	rsp := new(pb.GrpUserAddRsp)
	rsp.Err = new(pb.ErrorMsg)
	rsp.Err.Code = 0
	rsp.Err.Msg = ""
	if err != nil {
		rsp.Err.Code = -1
		rsp.Err.Msg = err.Error()
	}

	return rsp, err
}

func (serv *TalkCloudService) RemoveGrpUser(ctx context.Context, req *pb.GrpUserDelReq) (*pb.GrpUserDelRsp, error) {
	err := group.RemoveGroupUser(req.Uid, req.Gid, db.DBHandler)
	rsp := new(pb.GrpUserDelRsp)
	rsp.Err = new(pb.ErrorMsg)
	rsp.Err.Code = 0
	rsp.Err.Msg = ""
	if err != nil {
		rsp.Err.Code = -1
		rsp.Err.Msg = ""
	}

	return rsp, err
}

func (serv *TalkCloudService) ExitGrp(ctx context.Context, req *pb.UserExitReq) (*pb.UserExitRsp, error) {
	return nil, nil
}

func (serv *TalkCloudService) RemoveGrp(ctx context.Context, req *pb.GroupDelReq) (*pb.GroupDelRsp, error) {
	// clear group user first
	err := group.ClearGroupUser(req.Gid, db.DBHandler)
	rsp := new(pb.GroupDelRsp)
	rsp.Err = new(pb.ErrorMsg)

	if err != nil {
		rsp.Err.Code = -1
		rsp.Err.Msg = err.Error()
	}

	// then remove group
	err = group.RemoveGroup(req.Gid, db.DBHandler)
	if err != nil {
		rsp.Err.Code = -2
		rsp.Err.Msg = err.Error()
	}

	return rsp, err
}

func (serv *TalkCloudService) GetGroupList(ctx context.Context, req *pb.GrpListReq) (*pb.GroupListRsp, error) {
	rsp, err := group.GetGroupList(req.Uid, db.DBHandler)
	return rsp, err
}

func (serv *TalkCloudService) SearchGroup(ctx context.Context, req *pb.GrpSearchReq) (*pb.GroupListRsp, error) {
	return group.SearchGroup(req.Target, db.DBHandler)
}
