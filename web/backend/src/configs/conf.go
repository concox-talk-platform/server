/*
@Time : 2019/4/20 17:11 
@Author : yanKoo
@File : conf
@Software: GoLand
@Description:
*/
package configs

import (
	"github.com/go-ini/ini"
	"log"
	"os"
)

var (
	WebPort       string

	GrpcAddr      string

	FILE_BASE_DIR string
	FILE_BASE_URL string
)

func init() {
	cfg, err := ini.Load("web_conf.ini") // 编译之后的执行文件所在位置的相对位置
	if err != nil {
		log.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	WebPort = cfg.Section("web_api").Key("port").String()

	FILE_BASE_DIR = cfg.Section("upload_file").Key("save_path").String()
	FILE_BASE_URL = cfg.Section("upload_file").Key("save_path_url").String()

	GrpcAddr = cfg.Section("grpc").Key("addr").String()
}
