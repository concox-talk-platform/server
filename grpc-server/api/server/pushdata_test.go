/*
@Time : 2019/4/16 13:57 
@Author : yanKoo
@File : pushdata_test
@Software: GoLand
@Description:
*/
package server

import (
	"log"
	"server/grpc-server/api/talk_cloud"
	"testing"
	"time"
)

func testPushDataExecutor(t *testing.T) {
	//re := []int32{333}
	//dc := &DataContext{Data:make(chan interface{}, 1)}
	//e := make(chan interface{}, 1)
	//
	//dc.Data <- "test"
	//dc.Receiver = re
	//if err := pushDataExecutor(dc, e); err != nil {
	//	t.Logf("gen error: %v", err)
	//}
}

func testGetOfflineImMsgFromDB(t *testing.T) {
	res, err := GetOfflineImMsgFromDB(&talk_cloud.StreamRequest{
		Uid: 1503,
	})
	log.Println(res, "-----------------------", err)
}

func testGetMsgFromRedis(t *testing.T) {
	pttImMsgImpl{}.Dispatcher(nil, nil)
}

func TestPushDataIntoRedis(t *testing.T) {

}

func testTIme(t *testing.T) {
	log.Println(time.Now().Unix())
}
