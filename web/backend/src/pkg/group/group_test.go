package group

import (
	"model"
	"testing"
)

func TestCreateGroupByWeb(t *testing.T) {
	nums := []int{9000, 10000, 11000, 12000, 13000}
	g := &model.GroupInfo{GroupName: "拉萨组", AccountId: 2, Id: 5}
	gl := &model.GroupList{DeviceIds: nums, GroupInfo: g}
	if err := CreateGroupByWeb(gl); err != nil {
		t.Errorf("create group test error: %v", err)
	}
}
