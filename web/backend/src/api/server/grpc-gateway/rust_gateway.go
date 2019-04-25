/*
@Time : 2019/4/19 18:25 
@Author : yanKoo
@File : rest_gateway
@Software: GoLand
@Description:
*/
package main

import (
	pb "api/talk_cloud"
	"flag"
	"github.com/go-ini/ini"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
)

var (
	GrpcServerPort string
	GatewayPort    string
)

func init() {
	cfg, err := ini.Load("grpc_conf.ini") // 编译之后的执行文件所在位置的相对位置
	if err != nil {
		log.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	GrpcServerPort = cfg.Section("grpc_server").Key("port").String()
	GatewayPort = cfg.Section("grpc_gateway").Key("linsen_port").String()
}

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}))

	opts := []grpc.DialOption{grpc.WithInsecure()}
	echoEndpoint   := flag.String("echo_endpoint", "localhost:"+GrpcServerPort, "endpoint of TalkCloud")
	err := pb.RegisterTalkCloudLocationHandlerFromEndpoint(ctx, mux, *echoEndpoint, opts)
	if err != nil {
		return err
	}

	log.Println("listening " + GatewayPort + " ...")
	return http.ListenAndServe(":"+ GatewayPort, mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
