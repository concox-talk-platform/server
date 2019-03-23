/**
* @Author: yanKoo
* @Date: 2019/3/11 16:51
* @Description:
 */
package service

import (
	"github.com/concox-talk-platform/server/web/backend/src/model"
	"github.com/concox-talk-platform/server/web/backend/src/pkg"
	"github.com/concox-talk-platform/server/web/backend/src/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)
var HEADER_FIELD_SESSION = "Authorization"

func GetSessionId(r *http.Request) string {
	return r.Header.Get(HEADER_FIELD_SESSION)
}

// Check if the current user has the permission
// Use session id to do the check
func ValidateAccountSession(r *http.Request, value interface{}) bool {
	sid := GetSessionId(r)
	if len(sid) == 0 {
		return false
	}

	// 1. 是否存在
	ifExist, err := IsExistsSession(sid)
	if err != nil || !ifExist {
		log.Printf("validateAccountSession err: %s", err)
		return false
	}

	// 2. 判断session和发送过来的用户名或者账户id是否匹配
	log.Printf("sid : %s", sid)
	log.Printf("value: %s", value)
	sObj, err := pkg.GetSessionValue(sid)
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
	//	_ = pkg.UpdateSession(nil, nil)
	//}
	res := ifExist
	return res
}

func nowInMilli() int64 {
	return time.Now().UnixNano() / 1000000
}

func DeleteExpiredSession(sid string) error {
	if err := DeleteSessionInfo(sid); err != nil {
		return err
	}
	return nil
}

// 生成Session 放进redis数据库，主要就是注册的时候用
func GenerateNewSessionId(un string) string {
	sid, _ := utils.NewUUID()
	ct := nowInMilli()
	ttl := ct + 30*60*1000 // Severside session valid time: 30 min

	ttlStr := strconv.FormatInt(ttl, 10)

	sInfo := &model.SessionInfo{SessionID: sid, UserName: un, TTL: ttlStr}

	// 把sessionId放进redis缓存
	if err := pkg.InsertSession(sInfo); err != nil {
		return ""
	}
	return sid
}

// 更新session 主要是登录的时候用
func UpdateUserSessionId(oldSInfo *model.SessionInfo) (string, error) {
	newSid, _ := utils.NewUUID()
	ct := nowInMilli()
	ttl := ct + 30*60*1000 // Severside session valid time: 30 min

	ttlStr := strconv.FormatInt(ttl, 10)
	newSInfo := &model.SessionInfo{
		SessionID: newSid,
		UserName:  oldSInfo.UserName,
		AccountId: oldSInfo.AccountId,
		TTL:       ttlStr,
	}
	//log.Printf("old session: %s", oldSInfo.SessionID)
	//log.Printf("new session: %s", newSInfo.SessionID)
	if err := pkg.UpdateSession(oldSInfo, newSInfo); err != nil {
		return "", err
	}

	return newSid, nil
}

// 判断session 是否存在
func IsExistsSession(sid string) (bool, error) {
	ifExist, err := pkg.ExistsSession(sid)
	if err != nil {
		return false, err
	}
	return ifExist, nil
}

// 删除session
func DeleteSessionInfo(sid string) error {
	if err := pkg.DeleteSession(sid); err != nil {
		return err
	}
	return nil
}

// TODO 判断session是否过期，如果过期，就返回空string和true
func IsSessionExpired(sid string) (string, bool) {
	return "", false
}
