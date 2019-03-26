package talk_cloud

import (
	"context"
	"pkg/group"
)

func (serv *TalkCloudService) CreateGroup(ctx context.Context, req *GroupNewReq) (*GroupNewRsp, error) {
	err := group.CreateGroup(req.Uid, req.GroupName, )
	rsp := new(GroupNewRsp)
	rsp.Err.Code = 0
	rsp.Err.Msg = ""

	if err != nil {
		rsp.Err.Code = -1
		rsp.Err.Msg = err.Error()
	}

	return rsp, err
}

func (serv *TalkCloudService) JoinGroup(ctx context.Context, req *GrpUserAddReq) (*GrpUserAddRsp, error) {
	err := group.AddGroupUser(req.Uid, req.Gid, group.GROUP_NORMAL_USER, )

	rsp := new(GrpUserAddRsp)
	rsp.Err.Code = 0
	rsp.Err.Msg = ""
	if err != nil {
		rsp.Err.Code = -1
		rsp.Err.Code = err.Error()
	}

	return rsp, err
}

func (serv *TalkCloudService) RemoveGrpUser(ctx context.Context, req *GrpUserDelReq) (*GrpUserDelRsp, error) {
	err := group.RemoveGroupUser(req.Uid, req.Gid, )
	rsp := new(GrpUserDelRsp)
	rsp.Err.Code = 0
	rsp.Err.Msg = ""
	if err != nil {
		rsp.Err.Code = -1
		rsp.Err.Msg = ""
	}

	return rsp, err
}

func (serv *TalkCloudService) ExitGrp(ctx context.Context, req *UserExitReq) (*UserExitRsp, error) {
	return nil, nil
}

func (serv *TalkCloudService) RemoveGrp(ctx context.Context, req *GroupDelReq) (*GroupDelRsp, error) {
	// clear group user first
	err := group.ClearGroupUser(req.Gid, )
	rsp := new(GroupDelRsp)

	if err != nil {
		rsp.Err.Code = -1
		rsp.Err.Msg = err.Error()
	}

	// then remove group
	err = group.RemoveGroup(req.Gid, )
	if err != nil {
		rsp.Err.Code = -2
		rsp.Err.Msg = err.Error()
	}

	return rsp, err
}

func (serv *TalkCloudService) GetGroupList(ctx context.Context, req *GrpListReq) (*GroupListRsp, error) {
	rsp, err := group.GetGroupList(req.Uid, )
	return rsp, err
}

func (serv *TalkCloudService) SearchGroup(ctx context.Context, req *GrpSearchReq) (*GroupListRsp, error) {
	return group.SearchGroup(req.Target, )
}
