/**
* @Author: yanKoo
* @Date: 2019/3/15 14:11
* @Description: 调试阶段redis存储 session
 */
package session

import (
	"encoding/json"
	"model"
	"github.com/gomodule/redigo/redis"
	"log"
)

// 添加一个session
func InsertSession(sInfo *model.SessionInfo) error {
	// 将session转换成json数据，注意：转换后的value是一个byte数组
	value, err := json.Marshal(*sInfo)
	if err != nil {
		return err
	}
	if _, err := rdsConn.Do("SET", sInfo.SessionID, value, "ex", 60*60*3); err != nil {
		return err
	}
	log.Printf("Insert generate new sessionid : %s", sInfo.SessionID)

	//defer rdsConn.Close() // redis原子操作不需要close
	return nil
}

// 删除缓存中的session
func DeleteSession(sid string) error {
	log.Printf("[REDIS OPS] delete sesion : %s", err)
	if _, err := rdsConn.Do("DEL", sid); err != nil {
		return err
	}

	//defer rdsConn.Close()
	return nil
}

// 更新缓存中的session
func UpdateSession(oldSInfo, newSInfo *model.SessionInfo) error {
	if _, err := rdsConn.Do("DEL", oldSInfo.SessionID); err != nil {
		return err
	}

	log.Println("delete session is done")
	value, err := json.Marshal(*newSInfo)
	if err != nil {
		return err
	}
	if _, err := rdsConn.Do("SET", newSInfo.SessionID, value, "ex", 60*60*3); err != nil {
		return err
	}

	//defer rdsConn.Close()
	return nil
}

// 判断是否存在session
func ExistsSession(sid string) (bool, error) {
	ifExist, err := redis.Bool(rdsConn.Do("EXISTS", sid))
	if err != nil {
		return false, err
	}
	return ifExist, nil
}

// 获取session
func GetSessionValue(key string) (*model.SessionInfo, error) {
	if resBytes, err := redis.Bytes(rdsConn.Do("GET", key)); err != nil {
		return nil, err
	} else {
		res := &model.SessionInfo{}
		if err := json.Unmarshal(resBytes, res); err != nil {
			log.Println("json err")
			return nil, err
		}
		return res, nil
	}
}
