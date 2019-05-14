package server

import (
	"context"
	pb "server/grpc-server/api/talk_cloud"
	"server/grpc-server/dao/user_friend"
	"server/grpc-server/db"
	"server/grpc-server/log"
)

// 添加好友 TODO 暂时不等对方确认加不加好友，直接给你加上
func (serv *TalkCloudServiceImpl) AddFriend(ctx context.Context, req *pb.FriendNewReq) (*pb.FriendNewRsp, error) {
	log.Log.Printf("Add friend: uid: %d, friend_id:%d", req.Uid, req.Fuid)
	resp := &pb.FriendNewRsp{Res: &pb.Result{Msg: "Add friend error, please try again later", Code: 500}}
	_, err := user_friend.AddFriend(req.Fuid, req.Uid, db.DBHandler)
	if err != nil {
		log.Log.Printf("AddFriend friend add self error: %v", err)
		return resp, err
	}
	_, err = user_friend.AddFriend(req.Uid, req.Fuid, db.DBHandler)
	if err != nil {
		log.Log.Printf("AddFriend self add friend error: %v", err)
		return resp, err
	}
	return &pb.FriendNewRsp{Res: &pb.Result{Msg: "Add friend success", Code: 200}}, nil
}

func (serv *TalkCloudServiceImpl) GetFriendList(ctx context.Context, req *pb.FriendsReq) (*pb.FriendsRsp, error) {
	fList, _, err := user_friend.GetFriendReqList(req.Uid, db.DBHandler)
	if err != nil {
		return &pb.FriendsRsp{Res: &pb.Result{Code: 500, Msg: "process error, please try again"}}, err
	}
	fList.Res = &pb.Result{Msg: "Get friend list success", Code: 200}
	return fList, nil
}

// 根据关键字查询用户,携带是否好友字段
func (serv *TalkCloudServiceImpl) SearchUserByKey(ctx context.Context, req *pb.UserSearchReq) (*pb.UserSearchRsp, error) {
	if req.Target == "" {
		return &pb.UserSearchRsp{Res: &pb.Result{Code: 422, Msg: "process error, please input target"}}, nil
	}

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
	uSResp.Res = &pb.Result{Msg: "Search User success", Code: 200}
	return uSResp, nil
}

func (serv *TalkCloudServiceImpl) DelFriend(ctx context.Context, req *pb.FriendDelReq) (*pb.FriendDelRsp, error) {
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
