/*
@Time : 2019/4/29 10:36
@Author : yanKoo
@File : grpc_gateway
@Software: GoLand
@Description:
*/
package grpc_gateway

import (
	"github.com/go-ini/ini"
	"os"
	"server/grpc-server/log"
)

var (
	GrpcServerPort string
	GatewayPort    string
)

func init() {
	cfg, err := ini.Load("grpc_conf.ini") // 编译之后的执行文件所在位置的相对位置
	if err != nil {
		log.Log.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	GrpcServerPort = cfg.Section("grpc_server").Key("port").String()
	GatewayPort = cfg.Section("grpc_gateway").Key("linsen_port").String()
}
