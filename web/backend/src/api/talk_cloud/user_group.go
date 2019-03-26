package talk_cloud

import (
	"context"
)

func (serv *TalkCloudService) CreateGroup(context.Context, *GroupNewReq) (*GroupNewRsp, error) {
	return nil, nil
}

func (serv *TalkCloudService) JoinGroup(context.Context, *GrpUserAddReq) (*GrpUserAddRsp, error) {
	return nil, nil
}

func (serv *TalkCloudService) RemoveGrpUser(context.Context, *GrpUserDelReq) (*GrpUserDelRsp, error) {
	return nil, nil
}

func (serv *TalkCloudService) ExitGrp(context.Context, *UserExitReq) (*UserExitRsp, error) {
	return nil, nil
}

func (serv *TalkCloudService) RemoveGrp(context.Context, *GroupDelReq) (*GroupDelRsp, error) {
	return nil, nil
}

func (serv *TalkCloudService) GetGroupList(context.Context, *GrpListReq) (*GroupListRsp, error) {
	return nil, nil
}
