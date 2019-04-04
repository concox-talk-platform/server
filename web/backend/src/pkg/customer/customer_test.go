package customer

import (
	"database/sql"
	"log"
	"model"
	"testing"
)

func testAddAccount(t *testing.T) {
	if err := AddAccount(&model.Account{

		Pid:         9,
		Username:    "panda4",
		NickName:    "winna",
		Pwd:         "123456789",
		Email:       sql.NullString{String:"948162@qq.com"},
		Phone:       sql.NullString{String:"123456789"},
		Address:     sql.NullString{String:"株洲"},
		Remark:      sql.NullString{String:"熊猫"},
		PrivilegeId: 1,
		RoleId:      2,
	}); err != nil {
		t.Log("Test add account error : ", err)
	}
}

func testGetAccount(t *testing.T) {
	if res, err := GetAccount(1); err != nil {
		t.Log("Test error : ", err)
	} else {
		log.Println(res)
	}
}

func TestUpdateAccount(t *testing.T) {
	if err := UpdateAccount(&model.AccountUpdate{
		LoginId:     "9",
		Id: "11",
		NickName:    "ZZZZZZZZ",
		Email:       "948162@qq.com",
		Phone:       "123456789",
		Address:     "株洲",
		Remark:      "",
	}); err != nil {
		t.Log("Test add account error : ", err)
	}
}