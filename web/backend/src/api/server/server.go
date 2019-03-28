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

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("group net listen err: %v", err)
	}

	if err := groupServer.Serve(lis); err != nil {
		log.Fatalf("监听失败")
	} else {
		log.Println("listing")
	}
}
