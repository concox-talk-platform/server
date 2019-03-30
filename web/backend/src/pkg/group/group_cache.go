/**
 * Copyright (c) 2019. All rights reserved.
 * some functions deal with cache data of the groups
 * Author: tesion
 * Date: March 29th 2019
 */
package group

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
)

// Before when you change these constants
const (
	GRP_MEM_KEY_FMT = "grp:%d:mem"
	GRP_DATA_KEY_FMT = "grp:%d:data"
	USR_DATA_KEY_FMT = "usr:%d:data"
)

type MemStat uint8

const (
	MEM_ONLINE MemStat = iota
	MEM_OFFLINE
)

type MemRecord struct {
	Uid int64
	Name string
	Role RoleType
	Stat MemStat
}

type GrpRecord struct {
	Gid int64
	Name string
}

type GrpMemData struct {
	Gid int64
	GroupName string
	MemberList []*MemRecord
}

func MakeGroupMemKey(gid int64) string {
	return fmt.Sprintf(GRP_MEM_KEY_FMT, gid)
}

func MakeGroupDataKey(gid int64) string {
	return fmt.Sprintf(GRP_DATA_KEY_FMT, gid)
}

func MakeUserDataKey(uid int64) string {
	return fmt.Sprintf(USR_DATA_KEY_FMT, uid)
}

// check the key whether exists or not
func IsKeyExists(key string, rd redis.Conn) (bool, error) {
	if rd == nil {
		return false, fmt.Errorf("rd is null")
	}

	reply, err := redis.Int(rd.Do("EXISTS", key))
	if err != nil {
		log.Printf("check key(%s) exists error: %s\n", key, err)
		return false, nil
	}

	if 0 == reply {
		return false, nil
	}

	return true, nil
}

// add new group data to cache
func AddGroupCache(uid int64, usrName string, gid int64, grpName string, rd redis.Conn) (*GrpMemData, error) {
	if rd == nil {
		return nil, fmt.Errorf("rd is null")
	}

	grpData, err := json.Marshal(GrpRecord{Gid: gid, Name: grpName})

	if err != nil {
		log.Printf("json marshal error: %s\n", err)
		return nil, err
	}

	grpKey := MakeGroupDataKey(gid)
	memKey := MakeGroupMemKey(gid)



	rd.Send("MULTI")
	rd.Send("SADD",  memKey, uid)
	rd.Send("SET", grpKey, grpData)

	replys, err := redis.Ints(rd.Do("EXEC"))
	if err != nil {
		log.Printf("add group to cache error: %s\n", err)
		return nil, err
	}

	// Todo: 是否有必要检测所有返回结果？后续可以考虑注释掉
	allOk := true
	for _, status := range replys {
		if 0 == status {
			allOk = false
			break
		}
	}

	if !allOk {
		return nil, fmt.Errorf("add new group cache with some errors\n")
	}

	grp := new(GrpMemData)
	grp.Gid = gid
	grp.GroupName = grpName
	mem := new(MemRecord)
	mem.Name = usrName
	mem.Uid = uid
	mem.Stat = MEM_ONLINE
	grp.MemberList = append(grp.MemberList, mem)

	return grp, nil
}

// add new member to the group
func AddGroupMemCache(gid, uid int64, rd redis.Conn) error {
	if rd == nil {
		return fmt.Errorf("rd is nil")
	}

	key := MakeGroupMemKey(gid)
	_, err := rd.Do("SADD", key, uid)
	if err != nil {
		return fmt.Errorf("add new group key(%s) error: %s", key, err)
	}

	return nil
}

// get group metadata
func GetGroupData(gid int64, rd redis.Conn) (*GrpRecord, error) {
	if rd == nil {
		return nil, fmt.Errorf("rd is null")
	}

	key := MakeGroupDataKey(gid)

	reply, err := redis.Bytes(rd.Do("GET", key))
	if err != nil {
		log.Printf("get cache key(%s) error: %s\n", key, err)
		return nil, err
	}

	grpData := new(GrpRecord)
	err = json.Unmarshal(reply, grpData)
	if err == nil {
		log.Printf("json parse group data(%s) error: %s", reply, err)
		return nil, err
	}

	return grpData, nil
}

// get user list from the group
func GetGroupMemData(gid int64, rd redis.Conn) (*GrpMemData, error) {
	if rd == nil {
		return nil, fmt.Errorf("rd is nil")
	}

	grpMemData := new(GrpMemData)
	grpData, err := GetGroupData(gid, rd)
	if err != nil {
		log.Printf("get group(%d) metadata fail")
		return nil, err
	}

	grpMemData.Gid = grpData.Gid
	grpMemData.GroupName = grpData.Name
	key := MakeGroupMemKey(gid)
	uids, err := redis.Int64s(rd.Do("SMEMBERS", key))
	if err != nil {
		return nil, fmt.Errorf("get members from %s error: %s", key, err)
	}

	sz := len(uids)
	if 0 == sz {
		log.Printf("no members for group(%s)\n", key)
		return grpMemData, nil
	}

	memKeys := make([]interface{}, 0)
	for i := 0; i < sz; i++ {
		memKeys = append(memKeys, MakeUserDataKey(uids[i]))
	}

	mems, err := redis.ByteSlices(rd.Do("MGET", memKeys...))

	if err != nil {
		log.Printf("mget users data error: %s\n", err)
		return grpMemData, nil
	}

	sz = len(mems)
	for i := 0; i < sz; i++ {
		userData := new(MemRecord)
		err = json.Unmarshal(mems[i], userData)
		if err != nil {
			log.Printf("json parse user data(%s) error: %s\n", string(mems[i]), err)
			continue
		}
		grpMemData.MemberList = append(grpMemData.MemberList, userData)
	}

	return grpMemData, nil
}

func SetUserStat(uid int64, stat MemStat, rd redis.Conn) error {
	if rd == nil {
		return fmt.Errorf("rd is nil")
	}

	key := MakeUserDataKey(uid)

	reply, err := redis.Bytes(rd.Do("GET", key))
	if err != nil {
		log.Printf("get user(%s) error: %s\n", key, err)
		return err
	}

	user := new(MemRecord)
	err = json.Unmarshal(reply, &user)
	if err != nil {
		log.Printf("json parse user data(%s) error: %s\n", string(reply), err)
		return err
	}

	user.Stat = stat
	data, err := json.Marshal(&user)
	if err != nil {
		log.Printf("json encode user data(%s) error: %s\n", user, err)
		return err
	}

	res, err := redis.String(rd.Do("SET", key, data))
	if err != nil {
		log.Printf("set user data(%s) error: %s", data, err)
		return nil
	}

	if "OK" != res {
		log.Printf("set user data not ok, reply: %s\n", res)
		return fmt.Errorf("reply (%s) not ok", res)
	}

	return nil
}