package main

import (
	gServer "api/server"
	"api/talk_cloud"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	groupServer := grpc.NewServer()
	talk_cloud.RegisterTalkCloudServer(groupServer, &gServer.TalkCloudService{})
	talk_cloud.RegisterWebServiceServer(groupServer, &gServer.WebServiceServerImpl{})

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