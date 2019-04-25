/*
@Time : 2019/4/11 8:53 
@Author : yanKoo
@File : app_user_location_impl
@Software: GoLand
@Description: protoc -I. -I%GOPATH%/src -ID:\GoWorks\src\github.com\grpc-ecosystem\grpc-gateway\third_party\googleapis --go_out=plugins=grpc:. talk_cloud_location.proto
              protoc -I. -I%GOPATH%/src -ID:\GoWorks\src\github.com\grpc-ecosystem\grpc-gateway\third_party\googleapis --grpc-gateway_out=logtostderr=true:. talk_cloud_location.proto
*/
package server

import (
	pb "api/talk_cloud"
	"context"
	"log"
	"net/http"
	tl "pkg/location"
	"server/common/src/cache"
	"server/common/src/db"
)

const (
	MYSQL_STORE_SUCCESS = 1
	MYSQL_STORE_FAIL    = 0

	REDIS_UPDATE_SUCCESS = 1
	REDIS_UPDATE_FAIL    = 0
)

type TalkCloudLocationServiceImpl struct {
}

func (tcs *TalkCloudLocationServiceImpl) GetGpsData(ctx context.Context, req *pb.GPSHttpReq) (*pb.GPSHttpResp, error) {
	log.Printf("GetGpsData req : %v", req.Uid)
	res, err := tl.GetUserLocationInCache(req.Uid, cache.GetRedisClient())
	if err != nil {
		log.Printf("GetGpsData error: %+v", err)
	}
	return res, nil
}

type worker struct {
	dataStoreState   chan int
	updateRedisState chan int
}

// 处理上报数据
func (tcs *TalkCloudLocationServiceImpl) ReportGPSData(ctx context.Context, req *pb.ReportDataReq) (*pb.ReportDataResp, error) {
	// TODO 暂时认为存储到mysql之后，GPS数据一定更新成功, 更新不成，就去mysql查询最新一条数据出来
	gpsWorker := &worker{
		dataStoreState:   make(chan int, 1),
		updateRedisState: make(chan int, 1),
	}

	// 1. 首先对数据进行参数校验
	match, err := preCheckData(req)
	if !match {
		return &pb.ReportDataResp{Res: &pb.Result{Msg: "parmas is not correct", Code: http.StatusUnprocessableEntity}}, err
	}

	log.Printf("receiver data: %+v", req)
	// 2. 存储到mysql中
	go storeReportData(req, gpsWorker)

	// 3. 更新缓存中GPS数据
	go updateGPSDataByReq(req, gpsWorker)

	updateRedisState := <-gpsWorker.updateRedisState
	if updateRedisState == REDIS_UPDATE_FAIL {
		// 保证redis数据库里的数据一定是mysql中最新的那条记录
		updateGPSDataByMysql(req)
	}

	return &pb.ReportDataResp{Res: &pb.Result{Msg: "Receive data success", Code: 200}}, nil
}

// 存储到mysql上报数据
func storeReportData(req *pb.ReportDataReq, gw *worker) {
	if err := tl.InsertLocationData(req, db.DBHandler); err != nil {
		log.Printf("store data to mysql fail with error: %v", err)
		gw.dataStoreState <- MYSQL_STORE_FAIL
	} else {
		gw.dataStoreState <- MYSQL_STORE_SUCCESS
	}

}

// 更新缓存中的gps数据
func updateGPSDataByReq(req *pb.ReportDataReq, gw *worker) {
	mysqlState := <-gw.dataStoreState
	if mysqlState == MYSQL_STORE_SUCCESS {
		// 更新数据
		if err := tl.UpdateUserLocationInCache(req, cache.GetRedisClient()); err != nil {
			log.Printf("redis update data fail with error: %v", err)
			gw.updateRedisState <- REDIS_UPDATE_FAIL
		} else {
			gw.updateRedisState <- REDIS_UPDATE_SUCCESS
		}
	}

	// 暂时 如果mysql插入失败，就扔掉这个数据
	if mysqlState == MYSQL_STORE_FAIL {

	}
}

// 如果缓存更新失败，就去数据库里查询再来更新，估计不会出这样的问题
func updateGPSDataByMysql(req *pb.ReportDataReq) {
	// 去数据库查询数据

	// 更新缓存
}

// 校验数据合法性
func preCheckData(req *pb.ReportDataReq) (bool, error) {

	return true, nil
}
