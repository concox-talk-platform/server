/**
 * Copyright (c) 2019. All rights reserved.
 * some functions deal with cache data of the groups
 * Author: tesion
 * Date: March 29th 2019
 */
package group_cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	pb "server/grpc-server/api/talk_cloud"
	"server/web-api/dao/user_cache"
	"server/web-api/log"
	"server/web-api/model"
)

// Before when you change these constants
const (
	GRP_MEM_KEY_FMT   = "grp:%d:mem"
	GRP_DATA_KEY_FMT  = "grp:%d:data"
	USR_STATUS_KEY_FMT  = "usr:%d:stat"
	USR_GROUP_KEY_FMT = "usr:%d:grps"

	USER_OFFLINE = 1 // 用户离线
	USER_ONLINE  = 2 // 用户在线

	GROUP_MEMBER  = 1
	GROUP_MANAGER = 2
)

type MemStat uint8

const (
	MEM_ONLINE MemStat = iota
	MEM_OFFLINE
)

func MakeGroupMemKey(gid int32) string {
	return fmt.Sprintf(GRP_MEM_KEY_FMT, gid)
}

func MakeGroupDataKey(gid int32) string {
	return fmt.Sprintf(GRP_DATA_KEY_FMT, gid)
}

func MakeUserStatusKey(uid int32) string {
	return fmt.Sprintf(USR_STATUS_KEY_FMT, uid)
}

func MakeUserGroupKey(uid int32) string {
	return fmt.Sprintf(USR_GROUP_KEY_FMT, uid)
}

// check the key whether exists or not
func IsKeyExists(key string, rd redis.Conn) (bool, error) {
	if rd == nil {
		return false, fmt.Errorf("rd is null")
	}

	reply, err := redis.Int(rd.Do("EXISTS", key))
	if err != nil {
		log.Log.Printf("check key(%s) exists error: %s\n", key, err)
		return false, nil
	}

	if 0 == reply {
		return false, nil
	}

	return true, nil
}

// add new group data to cache
func AddGroupCache(ur []*pb.UserRecord, gInfo *pb.GroupInfo, rd redis.Conn) error {
	if rd == nil {
		return errors.New("rd is null")
	}
	defer rd.Close()

	grpData, err := json.Marshal(pb.GroupInfo{Gid: gInfo.Gid, GroupName: gInfo.GroupName, GroupManager: gInfo.GroupManager})
	if err != nil {
		log.Log.Printf("json marshal error: %s\n", err)
		return err
	}

	grpKey := MakeGroupDataKey(gInfo.Gid)
	memKey := MakeGroupMemKey(gInfo.Gid)

	_ = rd.Send("MULTI")
	for _, v := range ur {
		_ = rd.Send("SADD", memKey, v.Uid)
	}
	_ = rd.Send("SET", grpKey, grpData)

	_, err = rd.Do("EXEC")
	if err != nil {
		log.Log.Printf("add group to cache error: %s\n", err)
		return err
	}
	//log.Log.Printf("reply:%T", reply)

	return nil
}

// add new member to the group
func AddGroupSingleMemCache(gid, uid int32, rd redis.Conn) error {
	if rd == nil {
		return fmt.Errorf("rd is nil")
	}
	defer rd.Close()

	key := MakeGroupMemKey(gid)
	_, err := rd.Do("SADD", key, uid)
	if err != nil {
		return fmt.Errorf("add new group key(%s) error: %s", key, err)
	}

	return nil
}

// remove new member to the group
func RemoveGroupSingleMemCache(gid, uid int32, rd redis.Conn) error {
	if rd == nil {
		return fmt.Errorf("rd is nil")
	}
	defer rd.Close()

	key := MakeGroupMemKey(gid)
	_, err := rd.Do("SREM", key, uid)
	if err != nil {
		return fmt.Errorf("add new group key(%s) error: %s", key, err)
	}

	return nil
}

// 往群组成员里加成员
func AddGroupMemsInCache(gid int32, uids []int32, rd redis.Conn) error {
	if rd == nil {
		return fmt.Errorf("rd is nil")
	}
	defer rd.Close()

	key := MakeGroupMemKey(gid)
	_ = rd.Send("MULTI")
	for _, v := range uids {
		_ = rd.Send("SADD", key, v)
	}
	_, err := rd.Do("EXEC")
	if err != nil {
		return fmt.Errorf("add new group key(%s) error: %s", key, err)
	}

	return nil
}

func GetGroupMem(gid int32, rd redis.Conn) ([]int64, error) {
	if rd == nil {
		return nil, fmt.Errorf("rd is nil")
	}
	key := MakeGroupMemKey(gid)
	uids, err := redis.Int64s(rd.Do("SMEMBERS", key))
	if err != nil {
		return nil, fmt.Errorf("get members from %s error: %s", key, err)
	}

	return uids, nil
}

// 更新群组的信息
func AddGroupInCache(gl *pb.GroupListRsp, rd redis.Conn) error {
	if rd == nil {
		return errors.New("rd is null")
	}
	defer rd.Close()

	_ = rd.Send("MULTI")
	for _, v := range gl.GroupList {
		grpData, err := json.Marshal(&pb.GroupInfo{Gid: v.Gid, GroupName: v.GroupName, GroupManager: v.GroupManager})
		if err != nil {
			log.Log.Printf("json marshal error: %s\n", err)
			return err
		}
		grpKey := MakeGroupDataKey(v.Gid)
		_ = rd.Send("SET", grpKey, grpData)
	}

	if _, err := rd.Do("EXEC"); err != nil {
		log.Log.Printf("Add group data to cache error: %s\n", err)
		return err
	}

	return nil
}

func AddGroupAndUserInCache(gl *model.GroupList, rd redis.Conn) error {
	if rd == nil {
		return errors.New("rd is null")
	}
	defer rd.Close()
	log.Log.Println("start add create group to cache")

	grpData, err := json.Marshal(&pb.GroupInfo{Gid: int32(gl.GroupInfo.Id), GroupName: gl.GroupInfo.GroupName, GroupManager: int32(gl.GroupInfo.AccountId)})
	if err != nil {
		log.Log.Printf("json marshal error: %s\n", err)
		return err
	}

	grpKey := MakeGroupDataKey(int32(gl.GroupInfo.Id))
	memKey := MakeGroupMemKey(int32(gl.GroupInfo.Id))

	// TODO redis 错误处理
	_ = rd.Send("MULTI")

	// 1. 新建一个组，涉及到的是一个组加入了很多个成员，就有一个groupDataKey值和一个memberKey，
	for _, v := range gl.DeviceInfo {
		_ = rd.Send("SADD", memKey, v.(map[string]interface{})["id"])

	}
	// 2.更新每一个userGroups的key里面的组数
	for _, v := range gl.DeviceIds {
		_ = rd.Send("SADD", MakeUserGroupKey(int32(v)), gl.GroupInfo.Id)
	}
	_ = rd.Send("SET", grpKey, grpData)

	_, err = rd.Do("EXEC")
	if err != nil {
		log.Log.Printf("Add group to cache error: %s\n", err)
		return err
	}

	return nil
}

func GetGroupInfoFromCache(gId int32, rd redis.Conn) (*pb.GroupInfo, error) {
	key := MakeGroupDataKey(int32(gId))
	defer rd.Close()

	g, err := redis.String(rd.Do("GET", key))
	log.Log.Printf("get group info%+v", g)
	gInfo := &pb.GroupInfo{}
	err = json.Unmarshal([]byte(g), gInfo)
	log.Log.Printf("**********************%+v\n", gInfo)
	if err != nil {
		log.Log.Printf("json parse user data(%s) error: %s\n", string(g), err)
		return nil, err
	}
	// 获取每个群组中的userList
	gInfo.UsrList, err = user_cache.GetGroupMemDataFromCache(gId, rd)
	if err != nil {
		log.Log.Printf("get user from group (%s) error: %s\n", string(gId), err)
		return nil, err
	}
	return gInfo, nil
}

// 获取用户状态
func GetUserStatusFromCache(uId int32, redisCli redis.Conn) (int32, error) {
	if redisCli == nil {
		return -1, errors.New("redis conn is nil")
	}

	value, err := redis.Int(redisCli.Do("GET", MakeUserStatusKey(uId)))
	if err != nil {
		log.Log.Println("Get user online status fail with err", err.Error())
		return USER_OFFLINE, err
	}

	if value == 0 {
		return USER_OFFLINE, nil
	} else {
		return int32(value), nil
	}
}
