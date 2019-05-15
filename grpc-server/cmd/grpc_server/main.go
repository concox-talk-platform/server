package main

import (
	"flag"
	"google.golang.org/grpc"
	"net"
	"net/http"
	_ "net/http/pprof"
	gServer "server/grpc-server/api/server"
	"server/grpc-server/api/talk_cloud"
	cfgGs "server/grpc-server/configs/grpc_server"
	"server/grpc-server/log"
)

func init() {
	// 加载数据库中所有的数据到缓存
	//gServer.DataInit()
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
		log.Log.Errorf("group net listen err: %v", err)
	}

	//log.Log.Printf("listing %s", *p)
	log.Log.Infof("listing %s", cfgGs.GrpcSPort)
	if err := talkCloudServer.Serve(lis); err != nil {
		log.Log.Errorln("监听失败")
	} else {
		log.Log.Println("listing")
	}
}
