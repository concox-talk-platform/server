/**
* @Author: yanKoo
* @Date: 2019/3/11 16:51
* @Description:
 */
package session

import (
	"github.com/bittiger/server/web/interphone-api/dbops"
	"github.com/bittiger/server/web/interphone-api/defs"
	"github.com/bittiger/server/web/interphone-api/utils"
	"strconv"
	"time"
)

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

	sInfo := &defs.SessionInfo{SessionID: sid, UserName: un, TTL: ttlStr}

	// 把sessionId放进redis缓存
	if err := dbops.InsertSession(sInfo); err != nil {
		return ""
	}
	return sid
}

// 更新session 主要是登录的时候用
func UpdateUserSessionId(oldSInfo *defs.SessionInfo) (string, error) {
	newSid, _ := utils.NewUUID()
	ct := nowInMilli()
	ttl := ct + 30*60*1000 // Severside session valid time: 30 min

	ttlStr := strconv.FormatInt(ttl, 10)
	newSInfo := &defs.SessionInfo{
		SessionID: newSid,
		UserName:  oldSInfo.UserName,
		AccountId: oldSInfo.AccountId,
		TTL:       ttlStr,
	}
	//log.Printf("old session: %s", oldSInfo.SessionID)
	//log.Printf("new session: %s", newSInfo.SessionID)
	if err := dbops.UpdateSession(oldSInfo, newSInfo); err != nil {
		return "", err
	}

	return newSid, nil
}

// 判断session 是否存在
func IsExistsSession(sid string) (bool, error) {
	ifExist, err := dbops.ExistsSession(sid)
	if err != nil {
		return false, err
	}
	return ifExist, nil
}

// 删除session
func DeleteSessionInfo(sid string) error {
	if err := dbops.DeleteSession(sid); err != nil {
		return err
	}
	return nil
}

// TODO 判断session是否过期，如果过期，就返回空string和true
func IsSessionExpired(sid string) (string, bool) {
	return "", false
}
