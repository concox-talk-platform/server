/*
@Time : 2019/4/9 17:28 
@Author : yanKoo
@File : session_test
@Software: GoLand
@Description:
*/
package session

import (
	"server/web-api/model"
	"testing"
)

func testUpdateSession(t *testing.T) {
	new := &model.SessionInfo{
		SessionID: "123456",
	}
	//old := &model.SessionInfo{
	//	SessionID: "777777",
	//}
	if err := InsertSession(new); err != nil {

	}

}


func TestGetRedisKey(t *testing.T) {
	t.Log(getRedisKey(int(2), "dsgasg"))
	t.Log(getRedisKey("123456", "dsgasg"))
}