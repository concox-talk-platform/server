package group

import (
	"cache"
	"db"
	"model"
	"testing"
)



func testDeleteGroup(t *testing.T)  {
	if err := DeleteGroup(&model.GroupInfo{
		Id:101,
	}); err != nil {
		t.Logf("Test Delete Group error : %v", err)
	} else {
		t.Logf("sccess delete")
	}
}

func testCreateGroupByWeb(t *testing.T) {
	//nums := []int{9000, 10000, 11000, 12000, 13000}
	ds := make([]interface{}, 0)
	ds = append(ds, &model.User{
		Id: 333,
		UserName:  "123466",
		PassWord:  "123456",
		AccountId: 0, // TODO 默认给谁  普通用户默认是0
		IMei:      "123467894512365",
		UserType: 1,
	})

	//g := &model.GroupInfo{GroupName: "天津组1", AccountId: 2, Id: 333} // 用户id
	//gl := &model.GroupList{DeviceIds: nums, GroupInfo: g, DeviceInfo:ds}
	//if _, err := CreateGroupByWeb(gl, 1); err != nil {
	//	t.Errorf("create group test error: %v", err)
	//}
}

//func testGetGroupList(t *testing.T) {
//	if res, err := GetGroupList(uint64(333),db.DBHandler); err != nil {
//		t.Log("Get GroupList Error : ", err )
//	} else {
//		t.Log(*res)
//	}
//}

func testAddGroupCache(t *testing.T) {
	nums := []int{9000, 10000, 11000, 12000, 13000}
	ds := make([]interface{}, 0)

	for i := 1; i < 101; i++ {
		ds = append(ds, map[string]interface{}{
			"id": i,
			"UserName":  "123466",
			"PassWord":  "123456",
			//AccountId"": 0, // TODO 默认给谁  普通用户默认是0
			//IMei:      "123467894512365",
			//UserType: 1,
		})
	}

	g := &model.GroupInfo{GroupName: "天津组2", AccountId: 2, Id: 101} // 用户id
	gl := &model.GroupList{DeviceIds: nums, GroupInfo: g, DeviceInfo:ds}
	if err := AddGroupInCache(gl, cache.GetRedisClient()); err != nil {
		t.Logf("Add GroupCache error: %v", err)
	} else {
		//t.Logf("res:%v", res)
	}
}

//func testGetGroupList(t *testing.T) {
//	res, err := user.GetGroupList(uint64(1000), cache.GetRedisClient())
//	if err != nil && err != cache.NofindInCacheError {
//		t.Logf("Add GroupCache error: %v", err)
//	} else if err == cache.NofindInCacheError {
//		t.Logf("Can't find in cache")
//	} else {
//		t.Logf("find group list:%-v", res)
//	}
//}


func TestSearchGroup(t *testing.T) {
	res, err := SearchGroup("雷坤", db.DBHandler)
	if err != nil {
		t.Logf("Add GroupCache error: %v", err)
	} else {
		t.Logf("find group list:%-v", res)
	}

}