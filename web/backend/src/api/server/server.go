package main

import (
	"api"
	"api/talk_cloud"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	groupServer := grpc.NewServer()
	talk_cloud.RegisterTalkCloudServer(groupServer, &api.TalkCloudService{})
	talk_cloud.RegisterWebServiceServer(groupServer, &api.WebServiceServerImpl{})

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Printf("group net listen err: %v", err)
	}
	log.Println("listing")
	if err := groupServer.Serve(lis); err != nil {
		log.Printf("监听失败")
	} else {
		log.Println("listing")
	}
}
