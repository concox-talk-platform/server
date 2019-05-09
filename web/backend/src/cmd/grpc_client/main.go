/**
* @Author: yanKoo
* @Date: 2019/3/25 14:57
* @Description: 
*/
package main

import (
	pb "api/talk_cloud"
	"context"
	"google.golang.org/grpc"
	"log"
	"sync"
)

const GROUP_PORT = "9999"

var maps sync.Map

func main() {
	//host := "113.105.153.240"
	host := "127.0.0.1"

	conn, err := grpc.Dial(host+":9000", grpc.WithInsecure())
	if err != nil {
		log.Printf("grpc.Dial err : %v", err)
	}
	defer conn.Close()
	//webCli := pb.NewWebServiceClient(conn)
	//res, err := webCli.ImportDeviceByRoot(context.Background(), &pb.ImportDeviceReq{
	//	DeviceImei:[]string{"1234567897777777"},
	//	AccountId: 1,
	//})

	// 调用调用GRPC接口，转发数据

	if err != nil {
		log.Printf("grpc.Dial err : %v", err)
	}
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
			log.Printf("%+v, err:%+v",res.Result.Msg, err)
		}()
	}*/
	//streamCli := pb.NewStreamServiceClient(conn)
	userClient := pb.NewTalkCloudClient(conn)
	/*res, err := userClient.AppRegister(context.Background(), &pb.AppRegReq{
		Name:     "355172100003878",
		Password: "123456",
	})
	log.Printf(res.String())*/
	/*res, err := userClient.CreateGroup(context.Background(), &pb.CreateGroupReq{
		DeviceIds: "1482,333,1003,1004",
		DeviceInfos: nil,
		GroupName:"papi333组",
		AccountId: 333,
		GroupInfo: &pb.Group{
			Id:0,
			GroupName:"papi333组",
			AccountId: 333,
			Status: 1,
		}})
*/
	/*res, err := webCli.DeleteGroup(context.Background(), &pb.Group{
		Id: 102,
	})*/

	/*res, err := userClient.AddFriend(context.Background(), &pb.FriendNewReq{
		Uid:333,
		Fuid:500,
	})*/

	res, err := userClient.Login(context.Background(), &pb.LoginReq{
		Name:     "zebra3",
		Passwd: "123456",
	})
	log.Println(res.GroupList)
/*
	for _, v := range res.GroupList {
		if v.GroupManager == -1 {
			log.Printf("find a no mannage group: %d", v.Gid)
		}
	}

	log.Println("this user groups:", len(res.GroupList))
*/
/*	res, err := userClient.InviteUserIntoGroup(context.Background(), &pb.InviteUserReq{
		Uids:"1536",
		Gid:208,
	})

	log.Println(ress, "---------------++++++",  err)*/


	/*res, err := userClient.SearchUserByKey(context.Background(), &pb.UserSearchReq{
		Uid:uint64(333),
		Target:"",
	})*/
	/*res, err := userClient.GetFriendList(context.Background(), &pb.FriendsReq{
		Uid:333,
	})*/
	/*for i := 0; i < 3000; i++ {
		go func(i int) {
			resss, err := userClient.GetGroupList(context.Background(), &pb.GrpListReq{
				Uid:int32(333),
			})
			log.Println(resss, "-----------++++----",  err)
			log.Printf("Get group list **************************** # %d", i)
		}(i)

		go func() {
			ress, err := userClient.InviteUserIntoGroup(context.Background(), &pb.InviteUserReq{
				Uids:"457",
				Gid:210,
			})
			log.Println(ress, "*******---------++++++",  err)
			log.Printf("InviteUserIntoGroup **************************** # %d", i)
		}()
	}*/

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

	//log.Println(res, "---------------",  err)
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
	log.Printf("%+v", res)*/

	//if err != nil {
	//	log.Println(err)
	//} else {
	//	log.Printf("%+v", len(res.GroupList))
	//}

	// TODO 服务端 客户端 双向流
	//allStr, _ := userClient.DataPublish(context.Background())
	/*
		go func() {
			for {
				log.Println("start send heartbeat")
				if err := allStr.Send(&pb.StreamRequest{
					Uid:      334,
					DataType: 3,
					ACK:      334,
				}); err != nil {
				}

				time.Sleep(time.Second * 3)
			}
		}()*/
		//i := 335
	//for ; i < 1500 ; i++ {
		/*go func(i int) {
			log.Printf("%d start send get offline msg", i)
			if err := allStr.Send(&pb.StreamRequest{
				Uid:      int32(i),
				DataType: 2,
			}); err != nil {
			}
		}(i)
		time.Sleep(time.Microsecond*300)*/
	//}

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
			//	log.Println("client receive: 5 ", data.KeepAlive)
			//} else if data.DataType == OFFLINE_IM_MSG {
			//	log.Println("client receive: 2", data.OfflineImMsgResp)
			//} else if data.DataType == IM_MSG_FROM_UPLOAD_OR_WS_OR_APP {
			//	log.Println("client receive: 2", data.ImMsgData)
			//}
		}
	}(&allStr)*/
}
