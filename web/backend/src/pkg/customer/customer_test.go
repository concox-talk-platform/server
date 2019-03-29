package customer

import (
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
		Email:       "948162@qq.com",
		Phone:       "123456789",
		Address:     "株洲",
		Remark:      "熊猫",
		PrivilegeId: 1,
		RoleId:      2,
	}); err != nil {
		t.Log("Test add account error : ", err)
	}
}

func TestGetAccount(t *testing.T) {
	if res, err := GetAccount(1); err != nil {
		t.Log("Test error : ", err)
	} else {
		log.Println(res)
	}
}