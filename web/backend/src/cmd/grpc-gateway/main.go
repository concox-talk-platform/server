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
	cfgGgw "configs/grpc_gateway"
	"flag"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}))

	opts := []grpc.DialOption{grpc.WithInsecure()}
	echoEndpoint := flag.String("echo_endpoint", "localhost:"+cfgGgw.GrpcServerPort, "endpoint of TalkCloud")
	err := pb.RegisterTalkCloudLocationHandlerFromEndpoint(ctx, mux, *echoEndpoint, opts)
	if err != nil {
		return err
	}

	log.Println("listening " + cfgGgw.GatewayPort + " ...")
	return http.ListenAndServe(":"+cfgGgw.GatewayPort, mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
