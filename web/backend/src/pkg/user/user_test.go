package user

import (
	"model"
	"strconv"
	"testing"
)

func TestAddDevice(t *testing.T) {
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
			AccountId: "1",
			}

		if err := AddUser(d); err != nil {
			t.Errorf("Error of AddDevice %v", err)
		}
	}
}

