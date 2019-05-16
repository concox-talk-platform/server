/*
@Time : 2019/5/11 14:40 
@Author : yanKoo
@File : mian.go
@Software: GoLand
@Description:
*/
package main

import (
	"encoding/json"
	"server/grpc-server/log"
	"server/web-api/cache"
	"time"
)

type interphoneMsg struct {
	Uid       string `json:"uid"`
	MsgType   string `json:"m_type"`
	Md5       string `json:"md5"`
	GId       string `json:"grp_id"`
	FilePath  string `json:"file_path"`
	Timestamp string `json:"timestamp"`
}

func main() {
	r := cache.GetRedisClient()
	for i := 0; i < 1000; i++ {
		//go func(i int) {
		m := &interphoneMsg{
			Uid:       "8", //strconv.FormatInt(int64(i),10),
			MsgType:   "ptt",
			Md5:       "555555555555",
			GId:       "229",
			FilePath:  "123456789",
			Timestamp: time.Now().String(),
		}
		v, e := json.Marshal(m)
		log.Log.Printf("%s", v)
		if e != nil {

		}
		_, err := r.Do("lpush", "janusImUrlList", v)
		if err != nil {
			log.Log.Printf("push redis data with error: %+v", err)
			return
		}
		//}(i)
	}
	//select {}
}
