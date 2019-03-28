/**
* @Author: yanKoo
* @Date: 2019/3/25 14:57
* @Description: 
*/
package main


import (
	"context"
	"google.golang.org/grpc"
	"log"
	pb "server/web/backend/src/api/talk_cloud"
)

const GROUP_PORT = "9000"


func main() {
	//host := "113.105.153.240:"
	host := "172.16.0.74:"
	conn, err := grpc.Dial(host+GROUP_PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err : %v", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("connection err : %v", err)
		}
	}()

	groupCli := pb.NewTalkCloudClient(conn)
	resp, err := groupCli.AppRegister(context.Background(),&pb.AppRegReq{
		Name:"雷坤",
		Password:"123d",
	})
	//resp, err := groupCli.SelectGroups(context.Background(), &pb.TestReq{
	//	TestStr:"123456",
	//})
	if err != nil {
		log.Fatalf("group client err : %v", err)
	}

	log.Println(resp)

}
