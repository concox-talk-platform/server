/*
@Time : 2019/5/4 13:48 
@Author : yanKoo
@File : redis_data_sync
@Software: GoLand
@Description:
*/
package server

import (
	pb "api/talk_cloud"
	"cache"
	"database/sql"
	"db"
	"log"
	tg "pkg/group"
	tgc "pkg/group_cache"
	tu "pkg/user"
	tuc "pkg/user_cache"
	"sync"
)

type ConcurrentEngine struct {
	Scheduler   Scheduler
	WorkerCount int
}

type Scheduler interface {
	Submit(int32)
	ConfigureMasterWorkerChan(chan int32)
}

func (e ConcurrentEngine) Run() {
	in := make(chan int32)
	var wg sync.WaitGroup
	e.Scheduler.ConfigureMasterWorkerChan(in)

	for i := 0; i < e.WorkerCount; i++ {
		createWorker(in, &wg)
	}

	// 查找所有的用户id
	uIds, _ := tu.SelectAllUserId()
	for _, v := range uIds {
		log.Printf("# uid %d", v)
		wg.Add(1)
		go func() {e.Scheduler.Submit(v)}()
	}
	wg.Wait()
	log.Printf("**********************redis data sync done*****************************")
}

func createWorker(in chan int32, wg *sync.WaitGroup)  {
	go func() {
		for {
			uId := <- in
			err := UserData(uId)
			if err != nil {
				continue
			}
			err = GroupData(uId)
			if err != nil {
				continue
			}
			wg.Done()
		}
	}()
}

func DataInit() {
	uIds, _ := tu.SelectAllUserId()
	for _, v := range uIds {
		_ = UserData(v)
		_ = GroupData(v)

	}
}

func UserData(uId int32) error {
	// 根据用户id去获取每一位的信息，放进缓存
	res, err := tu.SelectUserByKey(int(uId))
	if err != nil && err != sql.ErrNoRows {
		log.Printf("UserData SelectUserByKey error : %s", err)
		return err
	}

	userInfo := &pb.Member{
		Id:          int32(res.Id),
		IMei:        res.IMei,
		UserName:    res.UserName,
		NickName:    res.NickName,
		UserType:    int32(res.UserType),
		LockGroupId: int32(res.LockGroupId),
		Online:      tuc.USER_OFFLINE, // 加载数据默认全部离线
	}
	log.Println("Add User Info into cache start")

	if err := tuc.AddUserDataInCache(userInfo, cache.GetRedisClient()); err != nil {
		log.Println("Add user information to cache with error: ", err)
	}
	log.Println("Add User Info into cache done")
	return nil
}

func GroupData(uid int32) error {
	gl, _, err := tg.GetGroupListFromDB(int32(uid), db.DBHandler)
	if err != nil {
		return err
	}
	log.Println("GroupData GetGroupListFromDB start update redis")
	// 新增到缓存 更新两个地方，首先，每个组的信息要更新，就是group data，记录了群组的id和名字
	if err := tgc.AddGroupInCache(gl, cache.GetRedisClient()); err != nil {
		return err
	}

	// 其次更新一个userSet  就是一个组里有哪些用户
	if err := tuc.AddUserInGroupToCache(gl, cache.GetRedisClient()); err != nil {
		return err
	}

	// 每个用户的信息
	for _, g := range gl.GroupList {
		for _, u := range g.UsrList {
			if err := tuc.AddUserDataInCache(&pb.Member{
				Id:          u.Uid,
				IMei:        u.Imei,
				NickName:    u.Name,
				Online:      u.Online,
				LockGroupId: u.LockGroupId,
			}, cache.GetRedisClient()); err != nil {
				log.Println("Add user information to cache with error: ", err)
			}
		}
	}

	// 每一个群组拥有的成员
	for _, v := range gl.GroupList {
		if err := tgc.AddGroupCache(v.UsrList, v, cache.GetRedisClient()); err != nil {
			return err
		}
	}
	return nil
}

