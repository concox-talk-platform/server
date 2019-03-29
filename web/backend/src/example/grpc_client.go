package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	pb "api/talk_cloud"
	"time"
)

const (
	ADDR = "localhost:9000"
)
func main() {
	// Set up a connection to the server
	conn, err := grpc.Dial(ADDR, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewTalkCloudClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.CreateGroup(ctx, &pb.GroupNewReq{Uid: 10, GroupName: "test_grpc"})
	if err != nil {
		log.Fatal("could not create group: %v", err)
	}

	log.Printf("Create group ok: %v\n", r)
}