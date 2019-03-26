/**
 * Copyrights 2019. All rights reserved.
 * Author: tesion
 * Date: March 22th 2019
 */

package user_friend

import (
	pb "api/talk_cloud"
	"database/sql"
	"fmt"
)

func CheckFriendExist(uid, fuid uint64, db *sql.DB) (bool, error) {
	if db == nil {
		return false, fmt.Errorf("db is nil")
	}

	sql := fmt.Sprintf("SELECT COUNT(*) FROM user_friend WHERE (uid=%[1]d AND friend_uid=%[2]d AND stat!=1) OR (uid=%[2]d AND friend_uid=%[1]d AND stat!=1)", uid, fuid)

	rows, err := db.Query(sql)
	if err != nil {
		return false, fmt.Errorf("query(%s) error(%s)", sql, err)
	}

	defer rows.Close()

	for rows.Next() {
		var cnt int
		err = rows.Scan(&cnt)
		if err != nil {
			return false, fmt.Errorf("query scan error(%s)", err)
		}

		if 0 >= cnt {
			return false, nil
		}
	}

	return true, nil
}

// 添加好友
func AddFriend(uid, fuid uint64, db *sql.DB) (bool, error) {
	ret, err := CheckFriendExist(uid, fuid, db)
	if err != nil {
		return false, err
	}

	if ret {
		return false, fmt.Errorf("user friend exists")
	}

	sql := fmt.Sprintf("INSERT INTO user_friend(uid, fuid) VALUES(%d, %d)", uid, fuid)
	_, err = db.Query(sql)
	if err != nil {
		return false, fmt.Errorf("query(%s), error(%s)", sql, err)
	}

	return true, nil
}

// 删除好友
func RemoveFriend(uid, fuid uint64, db *sql.DB) (bool, error) {
	if db == nil {
		return false, fmt.Errorf("db is nil")
	}

	sql := fmt.Sprintf("DELETE FROM user_friend WHERE uid=%d AND friend_uid=%d", uid, fuid)
	_, err := db.Query(sql)

	if err != nil {
		return false, err
	}

	return true, nil
}

// 获取好友请求列表
func GetFriendReqList(uid uint64, db *sql.DB) (*pb.FriendsRsp, error) {
    if db == nil {
        return nil, fmt.Errorf("db is nil")
    }
    
    sql := fmt.Sprintf("SELECT d.id, d.user_name, d.imei FROM device AS d RIGHT JOIN user_friend AS u ON u.friend_uid=d.id WHERE u.uid=%d", uid)
    rows, err := db.Query(sql)
    if err != nil {
        return nil, err
    }
    
    defer rows.Close()

    friends := &pb.FriendsRsp{Uid: uid, FriendList:nil}

    for rows.Next() {
        friend := new(pb.FriendRecord)
        err = rows.Scan(&friend.Uid, &friend.Name, &friend.Imei)
        if err != nil {
            return nil, err
        }
        
        friends.FriendList = append(friends.FriendList, friend)
    }
    
    return friends,nil
}

// 查找好友
func SearchFriend(fname string, db *sql.DB) (error) {

}
