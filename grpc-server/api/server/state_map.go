/*
@Time : 2019/4/16 17:47 
@Author : yanKoo
@File : redisMap
@Software: GoLand
@Description:
*/
package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"server/common/cache"
	pb "server/grpc-server/api/talk_cloud"
	"strconv"
)

func GetUserState(uIdKey []interface{}, rd redis.Conn) (map[int32]string, error) {
	if rd == nil {
		return nil, fmt.Errorf("rd is null")
	}
	/*
		states, err := redis.Strings(rd.Do("MGET", uIdKey...))
		for _, v := range states {
			gInfo := &pb.GroupInfo{}
			err = json.Unmarshal([]byte(v), gInfo)
			if err != nil {
				log.Printf("json parse user data(%s) error: %s\n", string(v), err)
				g, err := group.SelectGroupByKey(v)
				if err != nil {
					return nil, err
				}
				gInfo.GroupName = g.GroupName
				gInfo.Gid = int32(g.Id)
			}
			gList = append(gList, gInfo)
		}*/
	return nil, nil
}

// 保存连接信息
func AddUserStream(uId int32, srv pb.TalkCloud_DataPublishServer) error {
	rd := cache.GetRedisClient()
	if rd == nil {
		return fmt.Errorf("rd is nil")
	}

	// 先序列化后设置，同时设置在线时间
	srvStream, err := json.Marshal(srv)
	if err != nil {
		log.Printf("json marshal error: %s\n", err)
		return err
	}

	_, err = rd.Do("SET", "state:"+strconv.Itoa(int(uId)), srvStream, "ex", 5)
	if err != nil {
		return errors.New("SetUserStream error: " + err.Error())
	}
	return nil
}
