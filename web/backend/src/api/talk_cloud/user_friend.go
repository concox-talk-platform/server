package talk_cloud

import (
	"context"
	"pkg/user_friend"
)
func (serv *TalkCloudService) AddFriend(ctx context.Context, req *FriendNewReq) (*FriendNewRsp, error) {

	_, err := user_friend.AddFriend(req.Uid, req.Fuid, )
	rsp := new(FriendNewRsp)
	rsp.Err.Code = 0
	rsp.Err.Msg = ""
	if err != nil {
		rsp.Err.Code = -1
		rsp.Err.Msg = err.Error()
	}

	return rsp, err
}

func (serv *TalkCloudService) DelFriend(ctx context.Context, req *FriendDelReq) (*FriendDelRsp, error) {
	_, err := user_friend.RemoveFriend(req.Uid, req.Fuid, )
	rsp := new(FriendDelRsp)
	rsp.Err.Code = 0
	rsp.Err.Msg = ""

	if err != nil {
		rsp.Err.Code = -1
		rsp.Err.Msg = err.Error()
	}

	return rsp, err
}

func (serv *TalkCloudService) GetFriendList(ctx context.Context, req *FriendsReq) (*FriendsRsp, error) {
	return user_friend.GetFriendReqList(req.Uid, )
}

func (serv *TalkCloudService) SearchUser(ctx context.Context, req *UserSearchReq) (*UserSearchRsp, error) {
	return user_friend.SearchUserByName(req.Uid, req.Target, )
}
