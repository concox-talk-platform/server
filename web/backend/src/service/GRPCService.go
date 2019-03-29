/*
@Time : 2019/3/29 17:09 
@Author : yanKoo
@File : GRPCService
@Software: GoLand
@Description:
*/
package service

import (
	pb "api/talk_cloud"
	"google.golang.org/grpc"
	"log"
)

const GROUP_PORT = "9000"

var WebCli pb.WebServiceClient

func init() {
	//host := "113.105.153.240:"
	host := "172.16.0.74:"
	Conn, err := grpc.Dial(host+GROUP_PORT, grpc.WithInsecure())
	if err != nil {
		log.Printf("grpc.Dial err : %v", err)
	}

	WebCli = pb.NewWebServiceClient(Conn)
	log.Println("Web Service Client running ", WebCli)
}


