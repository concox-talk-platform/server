package api

import (
	"context"
	"db"
	"pkg/user_friend"
	pb "api/talk_cloud"
)
func (serv *TalkCloudService) AddFriend(ctx context.Context, req *pb.FriendNewReq) (*pb.FriendNewRsp, error) {

	_, err := user_friend.AddFriend(req.Uid, req.Fuid, db.DBHandler)
	rsp := new(pb.FriendNewRsp)
	rsp.Err = new(pb.ErrorMsg)
	rsp.Err.Code = 0
	rsp.Err.Msg = ""
	if err != nil {
		rsp.Err.Code = -1
		rsp.Err.Msg = err.Error()
	}

	return rsp, err
}

func (serv *TalkCloudService) DelFriend(ctx context.Context, req *pb.FriendDelReq) (*pb.FriendDelRsp, error) {
	_, err := user_friend.RemoveFriend(req.Uid, req.Fuid, db.DBHandler)
	rsp := new(pb.FriendDelRsp)
	rsp.Err = new(pb.ErrorMsg)
	rsp.Err.Code = 0
	rsp.Err.Msg = ""

	if err != nil {
		rsp.Err.Code = -1
		rsp.Err.Msg = err.Error()
	}

	return rsp, err
}

func (serv *TalkCloudService) GetFriendList(ctx context.Context, req *pb.FriendsReq) (*pb.FriendsRsp, error) {
	return user_friend.GetFriendReqList(req.Uid, db.DBHandler)
}

func (serv *TalkCloudService) SearchUser(ctx context.Context, req *pb.UserSearchReq) (*pb.UserSearchRsp, error) {
	return user_friend.SearchUserByName(req.Uid, req.Target, db.DBHandler)
}
