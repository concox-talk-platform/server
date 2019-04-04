package server

import (
	pb "api/talk_cloud"
	"context"
	"db"
	"log"
	"pkg/user_friend"
)

func (serv *TalkCloudService) AddFriend(ctx context.Context, req *pb.FriendNewReq) (*pb.FriendNewRsp, error) {
	log.Printf("Add friend: uid: %d, friend_id:%d", req.Uid, req.Fuid)
	_, err := user_friend.AddFriend(req.Uid, req.Fuid, db.DBHandler)
	if err != nil {
		return &pb.FriendNewRsp{Res:&pb.Result{Msg:"Add friend error, please try again later",Code: 500}}, err
	}
	return &pb.FriendNewRsp{Res:&pb.Result{Msg:"Add friend success",Code: 200}}, nil
}

func (serv *TalkCloudService) DelFriend(ctx context.Context, req *pb.FriendDelReq) (*pb.FriendDelRsp, error) {
	_, err := user_friend.RemoveFriend(req.Uid, req.Fuid, db.DBHandler)
	rsp := new(pb.FriendDelRsp)
	rsp.Err = new(pb.Result)
	rsp.Err.Code = 0
	rsp.Err.Msg = ""

	if err != nil {
		rsp.Err.Code = -1
		rsp.Err.Msg = err.Error()
	}

	return rsp, err
}

func (serv *TalkCloudService) GetFriendList(ctx context.Context, req *pb.FriendsReq) (*pb.FriendsRsp, error) {
	fList, _, err := user_friend.GetFriendReqList(req.Uid, db.DBHandler)
	if err != nil {
		return &pb.FriendsRsp{Res: &pb.Result{Code: 500, Msg: "process error, please try again"}}, err
	}
	fList.Res = &pb.Result{Msg:"Get friend list success", Code:200}
	return fList, nil
}

// 根据关键字查询用户,携带是否好友字段
func (serv *TalkCloudService) SearchUserByKey(ctx context.Context, req *pb.UserSearchReq) (*pb.UserSearchRsp, error) {
	// SearchUserByName
	uSResp, err := user_friend.SearchUserByName(req.Uid, req.Target, db.DBHandler)
	if err != nil {
		return &pb.UserSearchRsp{Res: &pb.Result{Code: 500, Msg: "process error, please try again"}}, err
	}
	// GetFriendReqList
	_, fMap, err := user_friend.GetFriendReqList(req.Uid, db.DBHandler)
	if err != nil {
		return &pb.UserSearchRsp{Res: &pb.Result{Code: 500, Msg: "process error, please try again"}}, err
	}

	for _, v := range uSResp.UserList {
		if _, ok := (*fMap)[v.Uid]; ok {
			v.IsFriend = true
		}
	}
	uSResp.Res = &pb.Result{Msg:"Search User success", Code:200}
	return uSResp, nil
}
