/**
* @Author: yanKoo
* @Date: 2019/3/15 14:11
* @Description: 调试阶段redis存储 session
 */
package session

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"server/web-api/cache"
	"server/web-api/log"
	"server/web-api/model"
	"strconv"
)

func getRedisKey(key interface{}, sId string) string {
	switch key.(type) {
	case int:
		return strconv.FormatInt(int64(key.(int)), 10) + ":" + sId
	case string:
		return key.(string) + ":" + sId
	}
	return ""
}

// 删除缓存中的session
func DeleteSession(s string, aInfo *model.Account) error {
	rdsCli := cache.GetRedisClient()
	if rdsCli == nil {
		return errors.New("redis connection is nil")
	}

	sInfo, err := GetSessionValue(s, aInfo.Username)
	if err != nil {
		return errors.New("delete session error: " + err.Error())
	}

	_ = rdsCli.Send("MULTI")
	_ = rdsCli.Send("DEL", getRedisKey(aInfo.Username, s))
	_ = rdsCli.Send("DEL", getRedisKey(sInfo.AccountId, s))
	_, err = rdsCli.Do("EXEC")
	if err != nil {
		log.Log.Println("delete session error: ", err)
		return err
	}

	return nil
}

// 更新缓存中的session
func InsertSession(sInfo *model.SessionInfo) error {
	rdsCli := cache.GetRedisClient()
	if rdsCli == nil {
		return errors.New("redis connection is nil")
	}

	value, err := json.Marshal(*sInfo)
	if err != nil {
		return err
	}
	_ = rdsCli.Send("MULTI")
	_ = rdsCli.Send("SET", getRedisKey(sInfo.UserName, sInfo.SessionID), value, "ex", 60*60*3)
	_ = rdsCli.Send("SET", getRedisKey(sInfo.AccountId, sInfo.SessionID), value, "ex", 60*60*3)
	_, err = rdsCli.Do("EXEC")
	if err != nil {
		return err
	}
	return nil
}

// 判断是否存在session
func ExistsSession(sid string, value interface{}) (bool, error) {
	log.Log.Println("sid, value :", sid, value)
	ifExist, err := redis.Bool(cache.GetRedisClient().Do("EXISTS", getRedisKey(value, sid)))
	if err != nil {
		return false, err
	}
	return ifExist, nil
}

// 获取session
func GetSessionValue(sId string, value interface{}) (*model.SessionInfo, error) {
	log.Log.Println("getRedisKey(value, key):", getRedisKey(value, sId))
	if resBytes, err := redis.Bytes(cache.GetRedisClient().Do("GET", getRedisKey(value, sId))); err != nil {
		return nil, err
	} else {
		res := &model.SessionInfo{}
		if err := json.Unmarshal(resBytes, res); err != nil {
			log.Log.Println("json err")
			return nil, err
		}
		return res, nil
	}
}
