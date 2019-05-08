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
	"api/talk_cloud"
	"testing"
)

func testPushDataExecutor(t *testing.T) {
	//re := []int32{333}
	//dc := &DataContent{DataChan:make(chan interface{}, 1)}
	//e := make(chan interface{}, 1)
	//
	//dc.DataChan <- "test"
	//dc.Receiver = re
	//if err := pushDataExecutor(dc, e); err != nil {
	//	t.Logf("gen error: %v", err)
	//}
}

func TestGetOfflineImMsgFromDB(t *testing.T) {
	res, err := GetOfflineImMsgFromDB(&talk_cloud.StreamRequest{
		Uid:1503,
	})
	log.Println(res, "-----------------------", err)
}