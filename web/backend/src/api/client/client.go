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
	"strconv"
	"sync"
	"time"
)

const GROUP_PORT = "9000"

var maps sync.Map

func main1()  {
	go func() {
		i := 1
		for {
			maps.Store(i, "one" + strconv.Itoa(i))
			time.Sleep(time.Second * 5)
			v1, ok1 := maps.Load(1)
			log.Println("get" ,v1, ok1)
			time.Sleep(time.Second * 500)
		}
	}()
	go func() {
		time.Sleep(time.Second * 2)
		for {
			v, ok := maps.Load(1)
			log.Println(v, ok)
			time.Sleep(time.Second * 5)
			maps.Delete(1)
			v1, ok1 := maps.Load(1)
			log.Println(v1, ok1)
			time.Sleep(time.Second * 500)
		}
	}()
	time.Sleep(time.Hour)
}

func main8(){
	c := make(chan int, 0)
	maps.Store(1, "one")
	go timerTask(c)
	go timerTasks(c)
    time.Sleep(time.Hour*12)
}


func timerTask(c chan int){
	timerTask := time.NewTicker(time.Second * time.Duration(1))
forLoop:
	for {
		select {
		case <-timerTask.C:
			loop:
			for {
				select {
				case <- c:
					break forLoop
				default:
					log.Println("get maps")
					break loop
				}
			}
		}
	}
}

func timerTasks(c chan int){
	time.Sleep(time.Second*5)
	log.Println("set c")
	c <- 3
}

func main() {
	//host := "113.105.153.240"
	host := "172.16.0.7"

	conn, err := grpc.Dial(host+":9000", grpc.WithInsecure())
	if err != nil {
		log.Printf("grpc.Dial err : %v", err)
	}

	//defer func() {
	//	if err := conn.Close(); err != nil {
	//		log.Printf("connection err : %v", err)
	//	}
	//}()
	//webCli := pb.NewWebServiceClient(conn)
	//res, err := webCli.ImportDeviceByRoot(context.Background(), &pb.ImportDeviceReq{
	//	DeviceImei:[]string{"1234567897777777"},
	//	AccountId: 1,
	//})

	userClient := pb.NewTalkCloudClient(conn)
	/*res, err := userClient.AppRegister(context.Background(), &pb.AppRegReq{
		Name:     "姚明6666",
		Password: "123456",
	})*/
	/*res, err := userClient.CreateGroup(context.Background(), &pb.CreateGroupReq{
		DeviceIds: "1482,333,1003,1004",
		DeviceInfos: nil,
		GroupName:"papi组",
		AccountId: 333,
		GroupInfo: &pb.Group{
			Id:0,
			GroupName:"papi酱40组",
			AccountId: 333,   // 如果写自己，岂不是普通用户也变成调度员了
			Status: 1,
		}})
*/
	/*res, err := webCli.DeleteGroup(context.Background(), &pb.Group{
	//	Id: 102,
	//})*/

	/*res, err := userClient.AddFriend(context.Background(), &pb.FriendNewReq{
		Uid:333,
		Fuid:500,
	})*/

	/*res, err := userClient.Login(context.Background(), &pb.LoginReq{
		Name:"264333",
		Passwd:"123456",
	})*/

	/*res, err := userClient.SearchUserByKey(context.Background(), &pb.UserSearchReq{
		Uid:uint64(333),
		Target:"",
	})*/
	/*res, err := userClient.GetFriendList(context.Background(), &pb.FriendsReq{
		Uid:333,
	})*/
	/*
		res, err := userClient.GetGroupList(context.Background(), &pb.GrpListReq{
			Uid:uint64(333),
		})*/

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
		Uid:333,
		Gid:152,
	})*/
	/*res, err := userClient.SetLockGroupId(context.Background(), &pb.SetLockGroupIdReq{
		UId:333,
		GId:215,
	})*/

	//服务端 客户端 双向流
	allStr, _ := userClient.DataPublish(context.Background())
	/*go func() {
			if err := allStr.Send(&pb.StreamRequest{
				Uid:      336,
				DataType: 1,
				Name:     "264336",
				Passwd:   "123456",
			}); err != nil {
			}
			time.Sleep(time.Second * 60)
	}()
	if err != nil {
		log.Println("error : ", err)
		return
	}*/
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

	go func() {
		for {
			log.Println("start send get offline msg")
			if err := allStr.Send(&pb.StreamRequest{
				Uid:      334,
				DataType: 2,
				ACK:      334,
			}); err != nil {
			}

			time.Sleep(time.Second * 600)
		}
	}()


	go func() {
		for {
			data, _ := allStr.Recv()
			log.Println("client receive :", data)
		}
	}()

	select {}
	/*time.Sleep(time.Second*30)

	ress, err := userClient.GetGroupList(context.Background(), &pb.GrpListReq{
		Uid:int32(333),
	})

	log.Println("ress:", ress)
*/
}
