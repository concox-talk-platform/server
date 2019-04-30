/**
* @Author: yanKoo
* @Date: 2019/3/11 16:51
* @Description:
 */
package service

import (
	"log"
	"model"
	"net/http"
	s "pkg/session"
	"strconv"
	"time"
	"utils"
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

	ifExist, err := IsExistsSession(sid, value)
	if err != nil{
		log.Printf("validateAccountSession err: %v", err)
		return false
	}
	return ifExist
}

func nowInMilli() int64 {
	return time.Now().UnixNano() / 1000000
}

// 更新session 主要是登录的时候用
func InsertSessionInfo(aInfo *model.Account) (string, error) {
	newSid, _ := utils.NewUUID()
	ct := nowInMilli()
	ttl := ct + 30*60*1000 // Severside session valid time: 30 min

	ttlStr := strconv.FormatInt(ttl, 10)
	sInfo := &model.SessionInfo{
		SessionID: newSid,
		UserName:  aInfo.Username,
		AccountId: aInfo.Id,
		TTL:       ttlStr,
	}

	if err := s.InsertSession(sInfo); err != nil {
		return "", err
	}

	return newSid, nil
}

// 判断session 是否存在
func IsExistsSession(sid string, value interface{}) (bool, error) {
	ifExist, err := s.ExistsSession(sid, value)
	if err != nil {
		return false, err
	}
	return ifExist, nil
}

// 删除session
func DeleteSessionInfo(session string , aInfo *model.Account) error {
	if err := s.DeleteSession(session, aInfo); err != nil {
		return err
	}
	return nil
}
