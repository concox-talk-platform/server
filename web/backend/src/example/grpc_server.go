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
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTalkCloudServer(s, &api.TalkCloudService{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
