/**
* @Author: yanKoo
* @Date: 2019/3/11 11:32
* @Description:
 */
package dbops

import (
	"github.com/concox-talk-platform/server/web/backend-api/defs"
	"strconv"
	"testing"
	//"strconv"
	//"time"
	//"fmt"
)

var tempvid string

func clearTables() {
	dbConn.Exec("truncate Accounts")
	dbConn.Exec("truncate video_info")
	dbConn.Exec("truncate comments")
	dbConn.Exec("truncate sessions")
}

//func TestMain(m *testing.M) {
//	clearTables()
//	m.Run()
//	clearTables()
//}

func TestAccountWorkFlow(t *testing.T) {
	//t.Run("Add", testAddAccount)
	//t.Run("Get", testGetAccount)
	//t.Run("Del", testDeleteAccount)
	//t.Run("Reget", testRegetAccount)
	//t.Run("upt", testUpdateAccount)
	//t.Run("updPwd", TestUpdateAccountPwd)
	//t.Run("ADD DEVICE", TestAddDevice)
	//t.Run("Select device id", TestSelectDeviceIdsByGroupId)
	//t.Run("groups", TestSelectGroupsByAccountId)
	//t.Run("devices", TestSelectDeviceByAccountId)
	//t.Run("create group",testGetSession)
	//t.Run("get Account Credential ", testGetAccountCredential)
	//t.Run("get group", testSelectGroupByGroupName)
	//t.Run("delete group", testDeleteGroup)
}

//func TestDBConn(t *testing.T) {
//
//}

func TestGetAccountByName(t *testing.T) {
	if s, err := GetAccountByName("tiger"); err != nil {
		t.Error("Get Session Error ", err)
	} else {
		t.Log(s)
	}
}

func testGetSessionValue(t *testing.T) {
	if s, err := GetSessionValue("bc191ccd-6208-43bf-aed0-6032d3dd347b"); err != nil {
		t.Error("Get Session Error ", err)
	} else {
		t.Log(s)
	}
}

func testInsertSession(t *testing.T) {
	var s = &defs.SessionInfo{
		SessionID: "23456789",
		UserName:  "tiger",
		AccountId: 2,
		TTL:       "123456789",
	}

	if err := InsertSession(s); err != nil {
		t.Errorf("insert Session error; %v", err)
	} else {
		t.Log("insert success")
	}
}

func testDeleteGroup(t *testing.T) {
	if err := DeleteGroup(&defs.GroupInfo{Id: 18, Status: "0"}); err != nil {
		t.Errorf("Delete group is wrong : %v", err)
	}
}

func testSelectGroupByGroupName(t *testing.T) {
	if pwd, err := SelectGroupByGroupName("重庆组"); err != nil {
		t.Errorf("Error of AccountCredential: %v", err)
	} else {
		t.Log("pwd :", pwd)
	}
}

func testGetAccountCredential(t *testing.T) {
	if pwd, err := GetAccountCredential(6); err != nil {
		t.Errorf("Error of AccountCredential: %v", err)
	} else {
		t.Log("pwd :", pwd)
	}
}

func testAddAccount(t *testing.T) {
	err := AddAccountCredential("avenssi", "123")
	if err != nil {
		t.Errorf("Error of AddAccount: %v", err)
	}
}

func testGetAccount(t *testing.T) {
	pwd, err := GetAccountCredential("avenssi")
	if pwd != "123" || err != nil {
		t.Errorf("Error of GetAccount")
	}
}

func testDeleteAccount(t *testing.T) {
	err := DeleteAccount("avenssi", "123")
	if err != nil {
		t.Errorf("Error of DeleteAccount: %v", err)
	}
}

func testUpdateAccount(t *testing.T) {
	if err := UpdateAccount(&defs.Account{Id: "2"}); err != nil {
		t.Errorf("Error of UpdateAccount %v", err)
	}
}

func testUpdateAccountPwd(t *testing.T) {
	if err := UpdateAccountPwd("123456789", 5); err != nil {
		t.Errorf("Error of UpdateAccount %v", err)
	}
}

func testAddDevice(t *testing.T) {
	base := int64(123456789123456)
	for i := int64(6); i < 3; i++ {
		imei := strconv.FormatInt(base-i, 10)
		d := &defs.Device{IMei: imei, AliasName: ("巡逻" + strconv.FormatInt(i, 10) + "号"), AccountId: 5, Status: "0", ActiveStatus: "0", BindStatus: "1"}
		if err := AddDevice(d); err != nil {
			t.Errorf("Error of AddDevice %v", err)
		}
	}
}

func testSelectDeviceIdsByGroupId(t *testing.T) {
	if ids, err := SelectDeviceIdsByGroupId(1); err != nil {
		t.Errorf("Error select devices id error : %v", err)
	} else {
		t.Log(ids)
	}
}

func testSelectGroupsByAccountId(t *testing.T) {
	if groups, err := SelectGroupsByAccountId(5); err != nil {
		t.Errorf("TestSelectGroupsByAccountId, error : %v", err)
	} else {
		t.Log(*groups[0])
		t.Log(len(groups))
	}
}

func testSelectDeviceByAccountId(t *testing.T) {
	if devices, err := SelectDeviceByAccountId(8); err != nil {
		t.Errorf("SelectDeviceByAccountId Test Error : %v", err)
	} else {
		t.Log(*devices[1])
	}
}

func testCreateGroup(t *testing.T) {
	nums := []int{9, 10, 11, 12, 13}
	g := &defs.GroupInfo{GroupName: "大连组", AccountId: 2, Id: 5}
	gl := &defs.GroupList{DeviceIds: nums, GroupInfo: g}
	if err := CreateGroup(gl); err != nil {
		t.Errorf("create group test error: %v", err)
	}
}

func testGetSession(t *testing.T) {
	if res, err := GetSessionValue("15ea3b6a-8e4f-4b30-94ad-ceef3f8a36f3"); err != nil {
		t.Errorf("create group test error: %v", err)
	} else {
		t.Log("res :", res)
	}
}
