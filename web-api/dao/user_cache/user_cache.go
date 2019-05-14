/*
@Time : 2019/4/3 14:35 
@Author : yanKoo
@File : user_cache
@Software: GoLand
@Description:
*/
package user_cache

import (
	pb "server/grpc-server/api/talk_cloud"
	"database/sql"
	"encoding/json"
	"server/web-api/model"
	"server/web-api/cache"
	"server/web-api/db"
	"strconv"
	//"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"server/web-api/log"
)

// Before when you change these constants
const (
	GRP_MEM_KEY_FMT   = "grp:%d:mem"
	GRP_DATA_KEY_FMT  = "grp:%d:data"
	USR_DATA_KEY_FMT  = "usr:%d:data"
	USR_STATUS_KEY_FMT  = "usr:%d:stat"
	USR_GROUP_KEY_FMT = "usr:%d:grps"

	USER_OFFLINE = 1 // 用户离线
	USER_ONLINE  = 2 // 用户在线
)

func MakeGroupMemKey(gid int32) string {
	return fmt.Sprintf(GRP_MEM_KEY_FMT, gid)
}

func MakeGroupDataKey(gid int32) string {
	return fmt.Sprintf(GRP_DATA_KEY_FMT, gid)
}

func MakeUserDataKey(uid int32) string {
	return fmt.Sprintf(USR_DATA_KEY_FMT, uid)
}

func MakeUserStatusKey(uid int32) string {
	return fmt.Sprintf(USR_STATUS_KEY_FMT, uid)
}


func MakeUserGroupKey(uid int32) string {
	return fmt.Sprintf(USR_GROUP_KEY_FMT, uid)
}

// 从缓存中获取群组列表
func GetGroupListFromRedis(uId int32, rd redis.Conn) (*pb.GroupListRsp, error) {
	if rd == nil {
		return nil, errors.New("redis conn is nil")
	}
	defer rd.Close()

	log.Log.Println("start get group list from redis")
	// How many groups is this user in
	gIds, err := redis.Int64s(rd.Do("SMEMBERS", MakeUserGroupKey(int32(uId))))
	if err != nil {
		return nil, err
	}
	log.Log.Println(gIds)

	sz := len(gIds)
	if 0 == sz {
		log.Log.Printf("user is not in any group\n")
		return nil, sql.ErrNoRows
	}

	gList := make([]*pb.GroupInfo, 0)
	keys := make([]interface{}, 0)
	for _, v := range gIds {
		keys = append(keys, MakeGroupDataKey(int32(v)))
	}

	groups, err := redis.Strings(rd.Do("MGET", keys...))
	log.Log.Printf("%+v", groups)
	for i, v := range groups {
		gInfo := &pb.GroupInfo{}
		log.Log.Printf("find group name %+v", v)
		if v == "" {
			// TODO 会报空针
			continue
		}
		err = json.Unmarshal([]byte(v), gInfo)
		log.Log.Printf("Get Group info from cache: %+v", gInfo)
		if err != nil {
			log.Log.Printf("json parse user data(%s) error: %s\n", string(v), err)
			return nil, err
		}

		// 获取每个群组中的userList
		gInfo.UsrList, err = GetGroupMemDataFromCache(int32(gIds[i]), cache.GetRedisClient())
		if err != nil {
			log.Log.Printf("get user from group (%s) error: %s\n", string(gIds[i]), err)
			return nil, err
		}
		gList = append(gList, gInfo)
	}

	return &pb.GroupListRsp{Uid: uId, GroupList: gList, Res: &pb.Result{Msg: "get group in cache success", Code: 200}}, nil
}

// 获取单个群组中的用户列表信息
func GetGroupMemDataFromCache(gid int32, rd redis.Conn) ([]*pb.UserRecord, error) {
	if rd == nil {
		return nil, fmt.Errorf("rd is nil")
	}
	defer rd.Close()

	res := make([]*pb.UserRecord, 0)
	resOffline := make([]*pb.UserRecord, 0)
	key := MakeGroupMemKey(gid)
	uids, err := redis.Int64s(rd.Do("SMEMBERS", key))
	if err != nil {
		return nil, fmt.Errorf("get members from %s error: %s", key, err)
	}

	log.Log.Printf("group %d have user id : %+v in cache", gid, uids)

	sz := len(uids)
	if 0 == sz {
		log.Log.Printf("group is not has any user\n")
		return nil, nil
	}

	memKeys := make([]interface{}, 0)
	for i := 0; i < sz; i++ {
		memKeys = append(memKeys, )
	}

	// 获取缓存中某个群成员信息
	for _, v := range uids {
		user := &pb.UserRecord{}
		value, err := redis.Values(rd.Do("HMGET", MakeUserDataKey(int32(v)),
			"id", "imei", "nickname", "lock_gid"))
		if err != nil {
			log.Log.Errorf("hmget failed", err.Error())
		}
		//log.Log.Printf("Get group %d user info value string  : %s from cache ", gid, value)

		var valueStr string
		resStr := make([]string, 0)
		for _, v := range value {
			if v != nil {
				valueStr = string(v.([]byte))
				resStr = append(resStr, valueStr)
			} else {
				break // redis找不到，去数据库加载
			}
		}
		online, err := redis.Int(rd.Do("GET", MakeUserStatusKey(int32(v))))
		if err != nil {
			log.Log.Errorf("GetGroupMemDataFromCache get user online with err ", err.Error())
			online = USER_OFFLINE
		}
		log.Log.Printf("Get group %d user info : %v from cache", gid, resStr)
		if value != nil {
			if value[0] != nil { // 只要任意一个字段为空就是没有这个数据
				uid, _ := strconv.Atoi(resStr[0])
				lockGId, _ := strconv.Atoi(resStr[3])

				user.Uid = int32(uid)
				user.Imei = resStr[1]
				user.Name = resStr[2]
				user.Online = int32(online)
				user.LockGroupId = int32(lockGId)

			} else {
				log.Log.Printf("can't find user %d from redis", int(v))
				UpdateUserFromDBToRedis(user, int(v))
			}

			// 在线离线顺序
			if user.Online == USER_ONLINE {
				res = append(res, user)
			} else {
				resOffline = append(resOffline, user)
			}
		}
	}
	res = append(res, resOffline...)
	return res, nil
}

func UpdateUserFromDBToRedis(user *pb.UserRecord, v int) {
	res, err := selectUserByKey(v)
	if err != nil {
		log.Log.Printf("GetGroupMemDataFromCache UpdateUserFromDBToRedis selectUserByKey has error: %v", err)
		return
	}
	user.Uid = int32(res.Id)
	user.Imei = res.IMei
	user.Name = res.NickName
	user.Online = USER_OFFLINE
	user.LockGroupId = int32(res.LockGroupId)
	// 增加到缓存
	if err := AddUserDataInCache(&pb.Member{
		Id:          user.Uid,
		IMei:        user.Imei,
		NickName:    user.Name,
		Online:      user.Online,
		LockGroupId: user.LockGroupId,
	}, cache.GetRedisClient()); err != nil {
		log.Log.Println("UpdateUserFromDBToRedis Add user information to cache with error: ", err)
	}
}

// 通过关键词查找用户名
func selectUserByKey(key interface{}) (*model.User, error) {
	var stmtOut *sql.Stmt
	var err error
	switch t := key.(type) {
	case int:
		stmtOut, err = db.DBHandler.Prepare("SELECT id, name, nick_name, passwd, imei, user_type, pid, cid, lock_gid, create_time, last_login_time, change_time FROM `user` WHERE id = ?")
	case string:
		stmtOut, err = db.DBHandler.Prepare("SELECT id, name, nick_name, passwd, imei, user_type, pid, cid, lock_gid, create_time, last_login_time, change_time  FROM `user` WHERE name = ?")
	default:
		_ = t
		return nil, err
	}
	if err != nil {
		log.Log.Printf("%s", err)
		return nil, err
	}

	var (
		id, userType, cId, lockGId                                    int
		pId, userName, nickName, pwd, iMei, cTime, llTime, changeTime string
	)
	err = stmtOut.QueryRow(key).Scan(&id, &userName, &nickName, &pwd, &iMei, &userType, &pId, &cId, &lockGId, &cTime, &llTime, &changeTime)
	if err != nil {
		return nil, err
	}

	res := &model.User{
		Id:          id,
		IMei:        iMei,
		UserName:    userName,
		PassWord:    pwd,
		NickName:    nickName,
		UserType:    userType,
		ParentId:    pId,
		AccountId:   cId,
		LockGroupId: lockGId,
		CreateTime:  cTime,
		LLTime:      llTime,
		ChangeTime:  changeTime,
	}

	defer func() {
		if err := stmtOut.Close(); err != nil {
			log.Log.Println("Statement close fail")
		}
	}()
	return res, nil
}

// 一个用户添加进组 可以在加载数据的时候用
func AddUserForSingleGroupCache(uId, gid int32, rd redis.Conn) error {
	if rd == nil {
		return errors.New("redis conn is nil")
	}
	defer rd.Close()

	res, err := rd.Do("SADD", MakeUserGroupKey(int32(uId)), gid)
	if err != nil {
		return err
	}
	log.Log.Println(res)
	return nil
}

// 一个用户在多个组， 用来更新，获取群组列表之后，去缓存中获取群组列表
func AddUsersGroupInCache(uid []int32, gId int32, rd redis.Conn) error {
	if rd == nil {
		return errors.New("redis conn is nil")
	}
	defer rd.Close()
	_ = rd.Send("MULTI")
	for _, v := range uid {
		_ = rd.Send("SADD", MakeUserGroupKey(int32(v)), gId)
	}
	_, err := rd.Do("EXEC")
	if err != nil {
		return err
	}
	return nil
}

// 移除单个用户
func RemoveUserForSingleGroupCache(uId, gid int32, rd redis.Conn) error {
	if rd == nil {
		return errors.New("redis conn is nil")
	}
	defer rd.Close()

	res, err := rd.Do("SREM", MakeUserGroupKey(int32(uId)), gid)
	if err != nil {
		return err
	}
	log.Log.Println(res)
	return nil
}

// 一个用户在多个组， 用来更新，获取群组列表之后，去缓存中获取群组列表
func AddUserInGroupToCache(gl *pb.GroupListRsp, rd redis.Conn) error {
	if rd == nil {
		return errors.New("redis conn is nil")
	}
	defer rd.Close()

	_ = rd.Send("MULTI")
	for _, v := range gl.GroupList {
		_ = rd.Send("SADD", MakeUserGroupKey(int32(gl.Uid)), v.Gid)
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
	defer redisCli.Close()
	//log.Log.Printf(">>>>> start AddUserDataInCache")
	if _, err := redisCli.Do("HMSET", MakeUserDataKey(m.Id),
		"id", m.Id, "imei", m.IMei, "username", m.UserName, "nickname", m.NickName, "online", m.Online, "lock_gid", m.LockGroupId); err != nil {
			//log.Log.Printf("AddUserDataInCache HMSET error: %+v",err)
		return errors.New("hSet failed with error: " + err.Error())
	}
	//log.Log.Printf(">>>>> done AddUserDataInCache")
	return nil
}

// 更新用户所在默认用户组，就是更新用户data
func UpdateLockGroupIdInCache(req *pb.SetLockGroupIdReq, redisCli redis.Conn) error {
	if redisCli == nil {
		return errors.New("redis conn is nil")
	}
	defer redisCli.Close()

	if _, err := redisCli.Do("HSET", MakeUserDataKey(req.UId), "lock_gid", req.GId); err != nil {
		return errors.New("UpdateLockGroupIdInCache hSet failed with error:" + err.Error())
	}
	return nil
}

// 更新用户在线状态
func UpdateOnlineInCache(m *pb.Member, redisCli redis.Conn) error {
	log.Log.Printf("start update user online state")
	if redisCli == nil {
		return errors.New("redis conn is nil")
	}
	defer redisCli.Close()

	if _, err := redisCli.Do("SET", MakeUserStatusKey(m.Id), m.Online, "ex", 5); err != nil {
		return errors.New("UpdateOnlineInCache hSet failed with error:" + err.Error())
	}
	return nil
}

// 获取用户状态
func GetUserStatusFromCache(uId int32, redisCli redis.Conn) (int32, error) {
	if redisCli == nil {
		return USER_OFFLINE, errors.New("redis conn is nil")
	}
	defer redisCli.Close()

	value, err := redis.Int(redisCli.Do("GET", MakeUserStatusKey(uId)))
	if err != nil {
		fmt.Println("get failed", err.Error())
		return USER_OFFLINE, err
	}

	log.Log.Printf("online value :%d", value)
	if value == 0 {
		return USER_OFFLINE, errors.New("no find")
	} else {
		return int32(value), nil
	}
}

// 获取单个成员信息
func GetUserFromCache(uId int32) (*pb.UserRecord, error) {
	rd := cache.GetRedisClient()
	defer rd.Close()

	user := &pb.UserRecord{}
	value, err := redis.Values(rd.Do("HMGET", MakeUserDataKey(uId),
		"id", "imei", "nickname", "lock_gid"))
	if err != nil {
		log.Log.Println("hmget failed", err.Error())
	}
	//log.Log.Printf("Get group %d user info value string  : %s from cache ", gid, value)

	var valueStr string
	resStr := make([]string, 0)
	for _, v := range value {
		if v != nil {
			valueStr = string(v.([]byte))
			resStr = append(resStr, valueStr)
		} else {
			break // redis找不到，去数据库加载
		}
	}
	online, err := redis.Int(rd.Do("GET", MakeUserStatusKey(int32(uId))))
	if err != nil {
		log.Log.Infof("get user online with err ", err.Error())
		online = USER_OFFLINE
	}
	log.Log.Printf("Get %d user info : %v from cache", uId, resStr)
	if value != nil && value[0] != nil { // 只要任意一个字段为空就是没有这个数据
		uid, _ := strconv.Atoi(resStr[0])
		lockGId, _ := strconv.Atoi(resStr[3])

		user.Uid = int32(uid)
		user.Imei = resStr[1]
		user.Name = resStr[2]
		user.Online = int32(online)
		user.LockGroupId = int32(lockGId)

	} else {
		log.Log.Printf("can't find user %d from redis", int(uId))
		UpdateUserFromDBToRedis(user, int(uId))
	}

	return user, nil
}
