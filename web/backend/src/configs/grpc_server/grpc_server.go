/*
@Time : 2019/4/29 10:35 
@Author : yanKoo
@File : grpc
@Software: GoLand
@Description:
*/
package grpc_server

import (
	"github.com/go-ini/ini"
	"log"
	"os"
)

var GrpcSPort string

func init() {
	cfg, err := ini.Load("grpc_conf.ini") // 编译之后的执行文件所在位置的相对位置
	if err != nil {
		log.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	GrpcSPort = cfg.Section("grpc_server").Key("port").String()
}
