package main

import (
	"api"
	pb "api/talk_cloud"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	PORT = ":9000"
)
func main() {
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Printf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTalkCloudServer(s, &api.TalkCloudServiceImpl{})

	if err := s.Serve(lis); err != nil {
		log.Printf("failed to serve: %v", err)
	}
}
