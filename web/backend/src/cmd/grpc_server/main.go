package main

import (
	gServer "api/server"
	"api/talk_cloud"
	cfgGs "configs/grpc_server"
	"google.golang.org/grpc"
	"log"
	"net"
)

func init()  {
	// 获取数据里所有的数据
	//server.DataInit()
}

func main() {
	talkCloudServer := grpc.NewServer()
	talk_cloud.RegisterTalkCloudServer(talkCloudServer, &gServer.TalkCloudServiceImpl{})
	talk_cloud.RegisterTalkCloudLocationServer(talkCloudServer, &gServer.TalkCloudLocationServiceImpl{})
	talk_cloud.RegisterWebServiceServer(talkCloudServer, &gServer.WebServiceServerImpl{})

	lis, err := net.Listen("tcp", ":" + cfgGs.GrpcSPort)
	if err != nil {
		log.Printf("group net listen err: %v", err)
	}
	log.Printf("listing %s", cfgGs.GrpcSPort)
	if err := talkCloudServer.Serve(lis); err != nil {
		log.Printf("监听失败")
	} else {
		log.Println("listing")
	}
}
