/*
@Time : 2019/4/29 10:35 
@Author : yanKoo
@File : grpc
@Software: GoLand
@Description:
*/
package grpc_server

import (
	"flag"
	"github.com/go-ini/ini"
	"os"
	"server/grpc-server/log"
)

var (
	GrpcSPort     string // grpc服务监听端口
	RedisCoMax    int    // 启动grpc服务的时候，读取redis工作的最大连接数
	FILE_BASE_URL string // 保存文件到fastdfs服务器之后的访问前缀（ip、域名）
	Interval      int    // im 心跳检测时间间隔
	PttMsgKey     string // ptt音视频数据在redis中的key
	PttWaitTime   int    // ptt音视频获取时阻塞等待的时间
	PprofAddr     string // pprof监听的地址
	ExpireTime    int
)

func init() {
	iniFilePath := flag.String("d", "grpc_conf.ini", "grpc server conf file path")
	flag.Parse()
	cfg, err := ini.Load(*iniFilePath) // 编译之后的执行文件所在位置的相对位置
	if err != nil {
		log.Log.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	GrpcSPort = cfg.Section("grpc_server").Key("port").String()
	RedisCoMax, _ = cfg.Section("grpc_server").Key("redisCoMax").Int()

	Interval, _ = cfg.Section("im").Key("heartbeat_interval").Int()

	PprofAddr = cfg.Section("pprof").Key("pprof_addr").String()

	FILE_BASE_URL = cfg.Section("upload_file").Key("save_path_url").String()

	PttMsgKey = cfg.Section("ptt").Key("ptt_msg_key").String()

	PttWaitTime, _ = cfg.Section("ptt").Key("wait_time").Int()

	ExpireTime, _ = cfg.Section("im").Key("expire_time").Int()
}
