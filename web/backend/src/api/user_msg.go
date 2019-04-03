package api

import (
	pb "api/talk_cloud"
	"context"
	"db"
	"log"
	"pkg/msg"
)

func (serv *TalkCloudService) AddMsg(ctx context.Context, req *pb.MsgNewReq) (*pb.MsgNewRsp, error) {
	err := msg.AddMultiMsg(req, db.DBHandler)

	rsp := new(pb.MsgNewRsp)
	rsp.Err = new(pb.ErrorMsg)

	if err != nil {
		rsp.Err.Code = -1
		rsp.Err.Msg = err.Error()
		log.Printf("add multi message fail\n")
	}

	return rsp, err
}

func (serv *TalkCloudService) GetMsg(ctx context.Context, req *pb.MsgReq) (*pb.MsgRsp, error) {
	rsp := new(pb.MsgRsp)
	rsp.Uid = req.Uid
	rsp.Stat = req.Stat
	var err error
	rsp.MsgList, err = msg.GetMsg(req.Uid, req.Stat, db.DBHandler)

	if err != nil {
		log.Printf("get user(%d) msg fail\n", req.Uid)
	}

	return rsp, err
}

func (serv *TalkCloudService) SetMsgStat(ctx context.Context, req *pb.MsgStatReq) (*pb.MsgStatRsp, error) {
	err := msg.SetMultiMsgStat(req.MsgIds, req.Stat, db.DBHandler)

	rsp := new(pb.MsgStatRsp)
	rsp.Err = new(pb.ErrorMsg)

	if err != nil {
		rsp.Err.Code = -1
		rsp.Err.Msg = err.Error()
		log.Printf("set multi msg stat fail")
	}

	return rsp, err
}

func (serv *TalkCloudService) DelMsg(ctx context.Context, req *pb.MsgDelReq) (*pb.MsgDelRsp, error) {
	rsp := new(pb.MsgDelRsp)
	rsp.Err = new(pb.ErrorMsg)

	err := msg.DeleteMsg(req.MsgIds, db.DBHandler)
	if err != nil {
		rsp.Err.Code = -1
		rsp.Err.Msg = err.Error()
		log.Printf("delete msg list fail")
	}

	return rsp, nil
}