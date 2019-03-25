/**
* @Author: yanKoo
* @Date: 2019/3/20 16:15
* @Description:
 */
package api

import (
	"context"
	pb "server/web/backend/src/api/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

type AccountService struct{}

func (a *AccountService) SignOut(ctx context.Context, req *pb.SignOutRequest) (*pb.SignOutResponse, error) {
	var res bool
	if req.GetSessionId() == "123" {
		res = true
	} else {
		res = false
	}
	return &pb.SignOutResponse{Result: res}, nil
}
func (a *AccountService) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {
	//var res bool
	//if req.GetSessionId() == "123" {
	//	res = true
	//}
	return nil, nil
}
func (a *AccountService) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	//var res bool
	//if req.GetSessionId() == "123" {
	//	res = true
	//}
	return nil, nil
}

const PORT = "9000"

func main() {
	server := grpc.NewServer()
	pb.RegisterAccountServiceServer(server, &AccountService{})

	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	if err := server.Serve(lis); err != nil {
		log.Println("监听失败")
	}
}
