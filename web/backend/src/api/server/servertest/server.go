package main

import (
	gServer "api/server"
	"api/talk_cloud"
	"github.com/go-ini/ini"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)
var GrpcServerPort string

func init() {
	cfg, err := ini.Load("grpc_conf.ini") // 编译之后的执行文件所在位置的相对位置
	if err != nil {
		log.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	GrpcServerPort = cfg.Section("grpc_server").Key("port").String()
}

func main() {
	talkCloudServer := grpc.NewServer()
	talk_cloud.RegisterTalkCloudServer(talkCloudServer, &gServer.TalkCloudServiceImpl{})
	talk_cloud.RegisterTalkCloudLocationServer(talkCloudServer, &gServer.TalkCloudLocationServiceImpl{})
	talk_cloud.RegisterWebServiceServer(talkCloudServer, &gServer.WebServiceServerImpl{})

	lis, err := net.Listen("tcp", ":"+GrpcServerPort)
	if err != nil {
		log.Printf("group net listen err: %v", err)
	}
	log.Printf("listing %s", GrpcServerPort)
	if err := talkCloudServer.Serve(lis); err != nil {
		log.Printf("监听失败")
	} else {
		log.Println("listing")
	}
}

/*
func doWorker(req string, w *worker){
	log.Println(req, "working...")
	time.Sleep(time.Second*10)
	log.Println(req, "working done")
	w.done <- 1
}

type worker struct {
	done chan int
}

func createWorker(req string) *worker {
	w := worker{done:make(chan int)}
	go doWorker(req, &w)
	return &w
}

func main(){
	w := createWorker("papi")
	log.Println("create done")
	_ = <- w.done
	log.Println("w.done")
}*/
