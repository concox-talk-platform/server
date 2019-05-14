package main

import (
	gServer "api/server"
	"api/talk_cloud"
	cfgGs "configs/grpc_server"
	"flag"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
)

func init() {
	// 加载数据库中所有的数据到缓存
	//server.DataInit()
	//server.ConcurrentEngine{
	//	Scheduler: &server.SimpleScheduler{},
	//	WorkerCount: 10,   // 加载redis数据的协程数
	//}.Run()
}

func main() {
	//p := flag.String("p", "9000", "grpc listening port")
	flag.Parse()

	go func() {
		_ = http.ListenAndServe(cfgGs.PprofAddr, nil)
	}()


	talkCloudServer := grpc.NewServer()
	talk_cloud.RegisterTalkCloudServer(talkCloudServer, &gServer.TalkCloudServiceImpl{})
	talk_cloud.RegisterTalkCloudLocationServer(talkCloudServer, &gServer.TalkCloudLocationServiceImpl{})
	talk_cloud.RegisterWebServiceServer(talkCloudServer, &gServer.WebServiceServerImpl{})

	lis, err := net.Listen("tcp", ":"+cfgGs.GrpcSPort)
	//lis, err := net.Listen("tcp", ":"+ *p)
	if err != nil {
		log.Printf("group net listen err: %v", err)
	}

	//log.Printf("listing %s", *p)
	log.Printf("listing %s", cfgGs.GrpcSPort)
	if err := talkCloudServer.Serve(lis); err != nil {
		log.Printf("监听失败")
	} else {
		log.Println("listing")
	}
}
