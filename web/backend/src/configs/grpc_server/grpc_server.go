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

var (
	GrpcSPort     string
	RedisCoMax    int
	FILE_BASE_URL string // 保存文件到fastdfs服务器之后的访问前缀（ip、域名）
	Interval      int  // im 心跳检测时间间隔
)

func init() {
	cfg, err := ini.Load("grpc_conf.ini") // 编译之后的执行文件所在位置的相对位置
	if err != nil {
		log.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	GrpcSPort = cfg.Section("grpc_server").Key("port").String()
	RedisCoMax, _ = cfg.Section("grpc_server").Key("redisCoMax").Int()

	Interval, _ = cfg.Section("im").Key("heartbeat_interval").Int()

	FILE_BASE_URL = cfg.Section("upload_file").Key("save_path_url").String()

}
