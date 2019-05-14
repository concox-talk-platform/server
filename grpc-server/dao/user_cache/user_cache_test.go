/*
@Time : 2019/5/14 11:49 
@Author : yanKoo
@File : user_cache_test
@Software: GoLand
@Description:
*/
package user_cache

import (
	"github.com/gomodule/redigo/redis"
	"server/grpc-server/cache"
	"testing"
)

func TestGetGroupMemDataFromCache(t *testing.T) {
	online, err := redis.Int(cache.GetRedisClient().Do("GET", MakeUserStatusKey(int32(7))))
	if err != nil {
		t.Error("get user online with err ", err.Error())
	}
	t.Log(online)
}
