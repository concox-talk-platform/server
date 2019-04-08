package user

import (
	"model"
	"server/common/src/cache"
	"strconv"
	"testing"
)

func testAddDevice(t *testing.T) {
	base := "12345678541"
	for i := int64(1350); i < 1450; i++ {
		imei := base + strconv.FormatInt(i, 10)
		d := &model.User{
			Id: int(i),
			IMei: imei,
			UserType: 1,
			PassWord: "123456",
			UserName: string([]byte(imei)[9:len(imei)]),
			NickName: ("开发" + strconv.FormatInt(i, 10) + "号"),
			AccountId: 1,
			}

		if err := AddUser(d); err != nil {
			t.Errorf("Error of AddDevice %v", err)
		}
	}
}

func testSelectUserByKey(t *testing.T) {
	if res, err := SelectUserByKey(1); err != nil {
		t.Logf("Test select user by key error: %v", err)
	} else {
		t.Log(res)
	}
}

func testAddUserCache(t *testing.T) {
	if err := AddUserCache(1000, 101, cache.GetRedisClient()); err != nil {
		t.Logf("Add user error: %v", err)
	} else {
	}
}
