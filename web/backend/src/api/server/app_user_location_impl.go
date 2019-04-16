/*
@Time : 2019/4/11 8:53 
@Author : yanKoo
@File : app_user_location_impl
@Software: GoLand
@Description:
*/
package server

import (
	pb "api/talk_cloud"
	"context"
)

type TalkCloudLocationServiceImpl struct {
}

type worker struct {
	done chan int
}

func createWorker(req *pb.ReportDataReq) worker {
	w := worker{done: make(chan int)}
	go storeReportData(req, w)
	return w
}

// 存储上报数据
func storeReportData(req *pb.ReportDataReq, w worker) {
	// ...

	w.done <- 1
}

// 直接用一个channel就行，没有必要用waitGroup TODO
func (tcs *TalkCloudLocationServiceImpl) ReportGPSData(ctx context.Context, req *pb.ReportDataReq) (*pb.ReportDataResp, error) {
	w := createWorker(req)
	_ = <- w.done


	return &pb.ReportDataResp{Res:&pb.Result{Msg:"Receive data success", Code:200}}, nil
}

