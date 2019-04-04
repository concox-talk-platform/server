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
)

const GROUP_PORT = "9000"

func main() {
	host := "113.105.153.240"
	//host := "172.16.0.74"
	//conn, err := grpc.Dial(host+GROUP_PORT, grpc.WithInsecure())
	conn, err := grpc.Dial(host+":9000", grpc.WithInsecure())
	if err != nil {
		log.Printf("grpc.Dial err : %v", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("connection err : %v", err)
		}
	}()

	//groupCli := pb.NewTalkCloudClient(conn)
	//resp, err := groupCli.AppRegister(context.Background(), &pb.AppRegReq{
	//	Name:     "雷坤",
	//	Password: "123d",
	//})
	//if err != nil {
	//	log.Printf("group client err : %v", err)
	//}
	//
	//log.Println(resp)
	//
	//if err != nil {
	//	log.Printf("grpc.Dial err : %v", err)
	//}

	//webCli := pb.NewWebServiceClient(conn)
	//res, err := webCli.ImportDeviceByRoot(context.Background(), &pb.ImportDeviceReq{
	//	DeviceImei:[]string{"1234567897777777"},
	//	AccountId: 1,
	//})
/*
	res, err := webCli.CreateGroup(context.Background(), &pb.CreateGroupReq{
		DeviceIds: []int64{1000,1001,1002},
		DeviceInfos: nil,
		GroupInfo: &pb.Group{
			Id:0,
			GroupName:"雷坤6组",
			AccountId: 0,   // 如果写自己，岂不是普通用户也变成调度员了
			Status: 1,
		}})*/

	//res, err := webCli.DeleteGroup(context.Background(), &pb.Group{
	//	Id: 102,
	//})


	userFriendCli := pb.NewTalkCloudClient(conn)
	/*res, err := userFriendCli.AddFriend(context.Background(), &pb.FriendNewReq{
		Uid:333,
		Fuid:500,
	})*/

	/*res, err := userFriendCli.Login(context.Background(), &pb.LoginReq{
		Name:"264333",
		Passwd:"123456",
	})*/

	res, err := userFriendCli.SearchUserByKey(context.Background(), &pb.UserSearchReq{
		Uid:uint64(333),
		Target:"",
	})
	/*res, err := userFriendCli.GetGroupList(context.Background(), &pb.GrpListReq{
		Uid:uint64(333),
	})*/
	/*
	res, err := userFriendCli.SearchGroup(context.Background(), &pb.GrpSearchReq{
		Uid:uint64(333),
		Target:"雷坤",
	})
	*/
/*
	res , err := userFriendCli.SearchUserByKey(context.Background(), &pb.UserSearchReq{
		Uid:333,
		Target:"121422",
	})*/
	if err != nil {
		log.Println("error : ", err)
		return
	}
	log.Println(res.UserList)
	//for _, v := range res.UserList {
	//	log.Println(v.Uid, v.Name, v.IsFriend)
	//}
}
