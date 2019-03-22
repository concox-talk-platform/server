/**
* @Author: yanKoo
* @Date: 2019/3/11 17:08
* @Description: Use session Id to check if the current user has the permission.
 */
package main

import (
	"github.com/concox-talk-platform/server/web/backend-api/dbops"
	"github.com/concox-talk-platform/server/web/backend-api/session"
	"log"
	"net/http"
)

var HEADER_FIELD_SESSION = "Authorization"

func GetSessionId(r *http.Request) string {
	return r.Header.Get(HEADER_FIELD_SESSION)
}

// Check if the current user has the permission
// Use session id to do the check
func validateAccountSession(r *http.Request, value interface{}) bool {
	sid := GetSessionId(r)
	if len(sid) == 0 {
		return false
	}

	// 1. 是否存在
	ifExist, err := session.IsExistsSession(sid)
	if err != nil || !ifExist {
		log.Printf("validateAccountSession err: %s", err)
		return false
	}

	// 2. 判断session和发送过来的用户名或者账户id是否匹配
	log.Printf("sid : %s", sid)
	log.Printf("value: %s", value)
	sObj, err := dbops.GetSessionValue(sid)
	log.Println(sObj)
	switch v := value.(type) {
	case int:
		if value != sObj.AccountId {
			log.Printf("validateAccountSession err: id is no match session id ")
			return false
		}
	case string:
		if value != sObj.UserName {
			log.Printf("validateAccountSession err: name is no match session id ")
			return false
		}
	default:
		_ = v
		return false
	}

	// 3. TODO 是否过期，过期就更新
	//if ifExist {
	//	_ = dbops.UpdateSession(nil, nil)
	//}
	res := ifExist
	return res
}
