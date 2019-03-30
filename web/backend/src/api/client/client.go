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
	//host := "113.105.153.240:"
	//host := "172.16.0.74:"
	//conn, err := grpc.Dial(host+GROUP_PORT, grpc.WithInsecure())
	conn, err := grpc.Dial("172.16.0.74:9000", grpc.WithInsecure())
	if err != nil {
		log.Printf("grpc.Dial err : %v", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("connection err : %v", err)
		}
	}()

	groupCli := pb.NewTalkCloudClient(conn)
	resp, err := groupCli.AppRegister(context.Background(), &pb.AppRegReq{
		Name:     "雷坤",
		Password: "123d",
	})
	if err != nil {
		log.Printf("group client err : %v", err)
	}

	log.Println(resp)

	if err != nil {
		log.Printf("grpc.Dial err : %v", err)
	}

	webCli := pb.NewWebServiceClient(conn)
	res, err := webCli.ImportDeviceByRoot(context.Background(), &pb.ImportDeviceReq{
		DeviceImei:[]string{"1234567897777777"},
		AccountId: 1,
	})
	if err != nil {
		log.Println("import device error : ", err)
		return
	}
	log.Println(res)
}
