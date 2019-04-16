/*
@Time : 2019/4/3 14:35 
@Author : yanKoo
@File : user_cache
@Software: GoLand
@Description:
*/
package user

import (
	pb "api/talk_cloud"
	"database/sql"
	"encoding/json"
	"pkg/group"

	//"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
)

// Before when you change these constants
const (
	GRP_MEM_KEY_FMT   = "grp:%d:mem"
	GRP_DATA_KEY_FMT  = "grp:%d:data"
	USR_DATA_KEY_FMT  = "usr:%d:data"
	USR_GROUP_KEY_FMT = "usr:%d:grps"
)

func MakeGroupMemKey(gid int64) string {
	return fmt.Sprintf(GRP_MEM_KEY_FMT, gid)
}

func MakeGroupDataKey(gid int64) string {
	return fmt.Sprintf(GRP_DATA_KEY_FMT, gid)
}

func MakeUserDataKey(uid int32) string {
	return fmt.Sprintf(USR_DATA_KEY_FMT, uid)
}

func MakeUserGroupKey(uid int64) string {
	return fmt.Sprintf(USR_GROUP_KEY_FMT, uid)
}

// 获取群组
func GetGroupList(uId int32, rd redis.Conn) (*pb.GroupListRsp, error) {
	if rd == nil {
		return nil, errors.New("redis conn is nil")
	}
	// How many groups is this user in
	gIds, err := redis.Int64s(rd.Do("SMEMBERS", MakeUserGroupKey(int64(uId))))
	if err != nil {
		return nil, err
	}
	log.Println(gIds)

	sz := len(gIds)
	if 0 == sz {
		log.Printf("user is not in any group\n")
		return nil, sql.ErrNoRows
	}

	gList := make([]*pb.GroupInfo, 0)
	keys := make([]interface{}, 0)
	for _, v := range gIds {
		keys = append(keys, MakeGroupDataKey(v))
	}

	groups, err := redis.Strings(rd.Do("MGET", keys...))
	for _, v := range groups {
		gInfo := &pb.GroupInfo{}
		err = json.Unmarshal([]byte(v), gInfo)
		if err != nil {
			log.Printf("json parse user data(%s) error: %s\n", string(v), err)
			g, err := group.SelectGroupByKey(v)
			if err != nil {
				return nil, err
			}
			gInfo.GroupName = g.GroupName
			gInfo.Gid = int32(g.Id)
		}
		gList = append(gList, gInfo)
	}

	return &pb.GroupListRsp{Uid: uId, GroupList: gList, Res: &pb.Result{Msg: "get group in cache success", Code: 200}}, nil
}

// 一个用户添加进组 可以在加载数据的时候用
func AddUserForSingleGroupCache(uId, gid int32, rd redis.Conn) error {
	if rd == nil {
		return errors.New("redis conn is nil")
	}

	res, err := rd.Do("SADD", MakeUserGroupKey(int64(uId)), gid)
	if err != nil {
		return err
	}
	log.Println(res)
	return nil
}

// 一个用户在多个组， 用来更新，获取群组列表之后，去缓存中获取群组列表
func AddUserInGroupToCache(gl *pb.GroupListRsp, rd redis.Conn) error {
	if rd == nil {
		return errors.New("redis conn is nil")
	}
	_ = rd.Send("MULTI")
	for _, v := range gl.GroupList {
		_ = rd.Send("SADD", MakeUserGroupKey(int64(gl.Uid)), v.Gid)
	}
	_, err := rd.Do("EXEC")
	if err != nil {
		return err
	}
	return nil
}

// 向缓存添加用户信息数据
func AddUserDataInCache(m *pb.Member, redisCli redis.Conn) error {
	if redisCli == nil {
		return errors.New("redis conn is nil")
	}

	if _, err := redisCli.Do("HMSET", MakeUserDataKey(m.Id),
		"id", m.Id, "imei", m.IMei, "username", m.UserName, "online", m.Online, "lock_gid", m.LockGroupId); err != nil {
		return errors.New("hSet failed with error: " + err.Error())
	}

	return nil
}

// 更新用户所在默认用户组，就是更新用户data
func UpdateLockGroupIdInCache(req *pb.SetLockGroupIdReq, redisCli redis.Conn) error {
	if redisCli == nil {
		return errors.New("redis conn is nil")
	}

	if _, err := redisCli.Do("HSET", MakeUserDataKey(req.UId), "lock_gid", req.GId); err != nil {
		return errors.New("UpdateLockGroupIdInCache hSet failed with error:" + err.Error())
	}
	return nil
}

// 更新用户在线状态
func UpdateOnlineInCache(m *pb.Member, redisCli redis.Conn) error {
	if redisCli == nil {
		return errors.New("redis conn is nil")
	}

	if _, err := redisCli.Do("HSET", MakeUserDataKey(m.Id), "online", m.Online); err != nil {
		return errors.New("UpdateLockGroupIdInCache hSet failed with error:" + err.Error())
	}
	return nil
}