/**
* @Author: yanKoo
* @Date: 2019/3/25 14:57
* @Description:
 */
package main

import (
	"context"
	"google.golang.org/grpc"
	pb "server/grpc-server/api/talk_cloud"
	"server/web-api/log"
	"sync"
	"time"
)

const GROUP_PORT = "9999"

var maps sync.Map

func main() {
	//host := "113.105.153.240"
	host := "127.0.0.1"

	conn, err := grpc.Dial(host+":9001", grpc.WithInsecure())
	if err != nil {
		log.Log.Printf("grpc.Dial err : %v", err)
	}
	defer conn.Close()
	userClient := pb.NewTalkCloudClient(conn)
	//res, err := webCli.ImportDeviceByRoot(context.Background(), &pb.ImportDeviceReq{
	//	DeviceImei:[]string{"1234567897777777"},
	//	AccountId: 1,
	//})

	// 调用调用GRPC接口，转发数据
	/*for i := 0; i <1000; i++ {
		go func() {
			res, err := webCli.ImMessagePublish(context.Background(), &pb.ImMsgReqData{
				Id:           8,
				SenderName:   "xiaoliu",
				ReceiverType: 2,
				ReceiverId:   229,
				ResourcePath: "SOSOS",
				MsgType:      6,
				ReceiverName: "xx group",
				SendTime:     strconv.FormatInt(time.Now().Unix(), 10),
			})

			log.Log.Printf("res:%+v err : %+v",res, err)
		}()
	}
	select {}*/
	/*webCli := pb.NewTalkCloudClient(conn)
	for i:= 0; i < 1500; i++ {
		time.Sleep(time.Microsecond*300)
		go func() {
			res, err := webCli.ImMessagePublish(context.Background(), &pb.ImMsgReqData{
				Id:           8,
				SenderName:   "iron man",
				ReceiverType: 2,
				ReceiverId:   229,
				ResourcePath: "SOS",
				MsgType:      3,
				ReceiverName: "boot",
				SendTime:     "hhhhh",
				MsgCode:      strconv.FormatInt(time.Now().Unix(), 10),
			})
			log.Log.Printf("%+v, err:%+v",res.Result.Msg, err)
		}()
	}*/
	/*res, err := userClient.AppRegister(context.Background(), &pb.AppRegReq{
		Name:     "355172100003878",
		Password: "123456",
	})
	log.Log.Printf(res.String())*/

	/*res, err := webCli.DeleteGroup(context.Background(), &pb.Group{
		Id: 102,
	})*/

	/*res, err := userClient.AddFriend(context.Background(), &pb.FriendNewReq{
		Uid:333,
		Fuid:500,
	})*/

	//log.Log.Println("---------------------------------------Login Start-------------------------------------------")
	//res, err := userClient.Login(context.Background(), &pb.LoginReq{
	//	Name:   "264333",
	//	Passwd: "123456",
	//})
	//log.Log.Println(res.GroupList)
	//
	//time.Sleep(time.Second * 3)
	/*
		for _, v := range res.GroupList {
			if v.GroupManager == -1 {
				log.Log.Printf("find a no mannage group: %d", v.Gid)
			}
		}*/

	//log.Log.Printf("this user groups:%d, all:%+v", len(res.GroupList), res)

	/*	res, err := userClient.InviteUserIntoGroup(context.Background(), &pb.InviteUserReq{
			Uids:"1536",
			Gid:208,
		})

		log.Log.Println(ress, "---------------++++++",  err)*/

	/*res, err := userClient.SearchUserByKey(context.Background(), &pb.UserSearchReq{
		Uid:uint64(333),
		Target:"",
	})*/
	/*res, err := userClient.GetFriendList(context.Background(), &pb.FriendsReq{
		Uid:333,
	})*/
	/*for i := 0; i < 3000; i++ {
		go func(i int) {
	//	*/
	//log.Log.Println("---------------------------------------GetGroupList Start-------------------------------------------")
	//resss, err := userClient.GetGroupList(context.Background(), &pb.GrpListReq{
	//	Uid: int32(333),
	//})
	//log.Log.Printf("%+v Get group list **************************** # %+v", err, resss)
	/*
		}(i)

		go func() {
			ress, err := userClient.InviteUserIntoGroup(context.Background(), &pb.InviteUserReq{
				Uids:"457",
				Gid:210,
			})
			log.Log.Println(ress, "*******---------++++++",  err)
			log.Log.Printf("InviteUserIntoGroup **************************** # %d", i)
		}()
	}*/

	//log.Log.Println("---------------------------------------Create group Start-------------------------------------------")
	//create, err := userClient.CreateGroup(context.Background(), &pb.CreateGroupReq{
	//	DeviceIds:   "1482,333,1003,1004",
	//	DeviceInfos: nil,
	//	GroupName:   "trsss组",
	//	AccountId:   333,
	//	GroupInfo: &pb.Group{
	//		Id:        0,
	//		GroupName: "trsss组",
	//		AccountId: 333,
	//		Status:    1,
	//	}})
	//log.Log.Printf("%+v>>>>>>>>>>>>>>>>>>>>>create:%+v", err, create)
	//
	//time.Sleep(time.Second * 2)
	//
	//log.Log.Println("---------------------------------------GetGroupList again Start-------------------------------------------")
	//glAgain, err := userClient.GetGroupList(context.Background(), &pb.GrpListReq{
	//	Uid: int32(333),
	//})
	//log.Log.Printf("%+v Get group list **************************** # %+v", err, glAgain)

	/*res, err := userClient.SearchGroup(context.Background(), &pb.GrpSearchReq{
		Uid:uint64(333),
		Target:"雷坤",
	})
	*/

	/*res , err := userClient.SearchUserByKey(context.Background(), &pb.UserSearchReq{
		Uid:333,
		Target:"121422",
	})*/

	/*res, err := userClient.JoinGroup(context.Background(), &pb.GrpUserAddReq{
		Uid: 1514,
		Gid: 151,
	})*/

	//log.Log.Println(res, "---------------",  err)
	/*res, err := userClient.SetLockGroupId(context.Background(), &pb.SetLockGroupIdReq{
		UId:333,
		GId:215,
	})*/

	// GPS 数据
	/*deviceCli := pb.NewTalkCloudLocationClient(conn)
	res, err := deviceCli.ReportGPSData(context.Background(), &pb.ReportDataReq{
		DataType: 2,
		DeviceInfo:&pb.Device{
			Id:1501,
		},
		LocationInfo: &pb.Location{
			GpsInfo: &pb.GPS{
				LocalTime: uint64(time.Now().Unix()),
				Longitude: float64( 113.13795866213755),
				Latitude: float64(22.480194593114472),
				Course:    123,
				Speed:     float32(123.456457),
			},
			BSInfo: &pb.BaseStation{
				Country: 460,
				Operator: 1,
				Lac:42705,
				Cid: 228408571,
				FirstBs: -49,
			},
			BtInfo: &pb.BlueTooth{
				FirstBt: -93,
				SecondBt: -89,
				ThirdBt:-98,
			},
			WifiInfo: &pb.Wifi{
				FirstWifi: 3,
				SecondWifi: 3,
				ThirdWifi: 3,
			},
		},
	})
	log.Log.Printf("%+v", res)*/

	//if err != nil {
	//	log.Log.Println(err)
	//} else {
	//	log.Log.Printf("%+v", len(res.GroupList))
	//}

	// TODO 服务端 客户端 双向流
	allStr, _ := userClient.DataPublish(context.Background())
	i := 8
	//for ; i < 1500 ; i++ {
	go func(i int) {
		log.Log.Printf("%d start send get offline msg", i)
		if err := allStr.Send(&pb.StreamRequest{
			Uid:      int32(i),
			DataType: 2,
		}); err != nil {
		}
	}(i)

	//}

	go func() {
		for {
			log.Log.Println("start send heartbeat")
			if err := allStr.Send(&pb.StreamRequest{
				Uid:      8,
				DataType: 4,
			}); err != nil {
			}

			time.Sleep(time.Second * 5)
		}
	}()

	const (
		FIRST_LOGIN_DATA                = 1 // 初次登录返回的数据。比如用户列表，群组列表，该用户的个人信息
		OFFLINE_IM_MSG                  = 2 // 用户离线时的IM数据
		IM_MSG_FROM_UPLOAD_OR_WS_OR_APP = 3 // APP和web通过httpClient上传的文件信息、在线时通信的im数据
		KEEP_ALIVE_MSG                  = 4 // 用户登录后，每隔interval秒向stream发送一个消息，测试能不能连通
		LOGOUT_NOTIFY_MSG               = 5 // 用户掉线之后，通知和他在一个组的其他成员
		LOGIN_NOTIFY_MSG                = 6 // 用户上线之后，通知和他在一个组的其他成员

		IM_MSG_FROM_UPLOAD_RECEIVER_IS_USER  = 1 // APP和web通过httpClient上传的IM信息是发给个人
		IM_MSG_FROM_UPLOAD_RECEIVER_IS_GROUP = 2 // APP和web通过httpClient上传的IM信息是发给群组

		USER_OFFLINE = 1 // 用户离线
		USER_ONLINE  = 2 // 用户在线

		UNREAD_OFFLINE_IM_MSG = 1 // 用户离线消息未读
		READ_OFFLINE_IM_MSG   = 2 // 用户离线消息已读

		CLIENT_EXCEPTION_EXIT = -1 // 客户端异常终止

		NOTIFY = 1 // 通知完一个
	)
	/*go func(allStr *pb.TalkCloud_DataPublishClient) {
		for {
			_, _ = (*allStr).Recv()
			//if data.DataType == KEEP_ALIVE_MSG {
			//	log.Log.Println("client receive: 5 ", data.KeepAlive)
			//} else if data.DataType == OFFLINE_IM_MSG {
			//	log.Log.Println("client receive: 2", data.OfflineImMsgResp)
			//} else if data.DataType == IM_MSG_FROM_UPLOAD_OR_WS_OR_APP {
			//	log.Log.Println("client receive: 2", data.ImMsgData)
			//}
		}
	}(&allStr)*/
	select {}
}
