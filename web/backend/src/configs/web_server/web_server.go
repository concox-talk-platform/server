/*
@Time : 2019/4/20 17:11 
@Author : yanKoo
@File : conf
@Software: GoLand
@Description:
*/
package web_server

import (
	"github.com/go-ini/ini"
	"log"
	"os"
)

var (
	WebPort           string // web api的监听端口
	GrpcAddr          string // grpc服务端的地址（包含端口）， 主要用于web模块调用grpc服务
	FILE_BASE_URL     string // 保存文件到fastdfs服务器之后的访问前缀（ip、域名）
	TrackerServerAddr string // fastdfs的tracker服务器的地址（包含ip）
	MaxConn           int    // fastdfs上传文件的最大连接数
	CertFile          string // https pem文件名
	KeyFile           string // https key文件
)

func init() {
	cfg, err := ini.Load("web_conf.ini") // 编译之后的执行文件所在位置的相对位置
	if err != nil {
		log.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	WebPort = cfg.Section("web_api").Key("port").String()

	FILE_BASE_URL = cfg.Section("upload_file").Key("save_path_url").String()

	GrpcAddr = cfg.Section("grpc").Key("addr").String()

	//fastdfs
	TrackerServerAddr = cfg.Section("fastdfs").Key("tracker_server").String()
	MaxConn, _ = cfg.Section("fastdfs").Key("maxConns").Int()

	// https
	CertFile = cfg.Section("https").Key("cert_file").String()
	KeyFile = cfg.Section("https").Key("key_file").String()
}
