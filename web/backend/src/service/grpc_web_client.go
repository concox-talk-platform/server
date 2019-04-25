/*
@Time : 2019/3/29 17:09 
@Author : yanKoo
@File : GRPCService
@Software: GoLand
@Description:
*/
package service

import (
	"google.golang.org/grpc"
	"server/web/backend/src/configs"
	cp "service/client_pool"
)

var Addr string

var CilentPool *cp.ConnectionTracker

func init() {
	//cfg, err := ini.Load("web_conf.ini")  // 编译之后的执行文件所在位置的相对位置
	//if err != nil {
	//	log.Printf("Fail to read file: %v", err)
	//	os.Exit(1)
	//}

	//Addr = cfg.Section("grpc").Key("addr").String()
	CilentPool = cp.New(func(addr string) (*grpc.ClientConn, error) {
		return grpc.Dial(configs.GrpcAddr,grpc.WithInsecure())
	})
}


