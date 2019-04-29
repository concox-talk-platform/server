package main

import (
	gServer "api/server"
	"api/talk_cloud"
	cfgGs "configs/grpc_server"
	"google.golang.org/grpc"
	"log"
	"net"
)

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
