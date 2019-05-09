/*
@Time : 2019/4/16 13:57 
@Author : yanKoo
@File : pushdata_test
@Software: GoLand
@Description:
*/
package server

import (
	"api/talk_cloud"
	"encoding/json"
	"log"
	"server/common/src/cache"
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
	m := &interphoneMsg{
		Uid:       9,
		MsgType:   "ptt",
		Md5:       "555555555555",
		GId:    229,
		FilePath: "123456789",
		Timestamp: time.Now().String(),
	}
	v, e := json.Marshal(m)
	log.Printf("%s", v)
	if e != nil {

	}
	_, err := cache.GetRedisClient().Do("lpush", "mylist", v)
	if err != nil {
		log.Printf("push redis data with error: %+v", err)
		return
	}
}
