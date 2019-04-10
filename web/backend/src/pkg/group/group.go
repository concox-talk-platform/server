/**
 * Copyrights (c) 2019. All rights reserved.
 * Group handlers
 * Author: tesion
 * Date: March 26 2019
 */
package group

import (
	pb "api/talk_cloud"
	"database/sql"
	"db"
	"errors"
	"fmt"
	"github.com/smartwalle/dbs"
	"log"
	"model"
)

var dbConn = db.DBHandler

func AddGroupMember(uid, gid int64, userType RoleType, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	sql := fmt.Sprintf("INSERT INTO group_member(gid, uid, role_type) VALUES(%d, %d, %d)", gid, uid, userType)
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("query(%s), error(%s)", sql, err)
		return err
	}

	defer rows.Close()

	return nil
}

func RemoveGroupMember(uid, gid int32, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	sql := fmt.Sprintf("DELETE FROM group_member WHERE uid=%d AND gid=%d", uid, gid)
	_, err := db.Query(sql)
	if err != nil {
		log.Printf("query(%s), error(%s)", sql, err)
		return err
	}

	return nil
}

func RemoveGroupUser(uid, gid uint64, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	sql := fmt.Sprintf("DELETE FROM group_member WHERE uid=%d AND gid=%d", uid, gid)
	_, err := db.Query(sql)
	if err != nil {
		log.Printf("query(%s), error(%s)", sql, err)
		return err
	}

	return nil
}

func RemoveGroup(gid int32, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	sql := fmt.Sprintf("DELETE FROM group WHERE id=%d", gid)
	_, err := db.Query(sql)
	if err != nil {
		log.Printf("remove group(%d) error: %s\n", gid, err)
		return err
	}

	return nil
}

func ClearGroupMember(gid int32, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	sql := fmt.Sprintf("DELETE FROM group_member WHERE gid=%d", gid)
	_, err := db.Query(sql)
	if err != nil {
		log.Printf("clear gruop(%d) user error: %s\n", gid, err)
		return err
	}

	return nil
}

// 获取该用户在哪几个群组
func GetGroupList(uid int32, db *sql.DB) (*pb.GroupListRsp, *map[int32]string, error) {
	if db == nil {
		return nil, nil, errors.New("db is nil")
	}

	sql := fmt.Sprintf("SELECT g.id, g.group_name "+
		"FROM user_group AS g INNER JOIN group_member AS gm "+
		"ON g.id=gm.gid WHERE gm.uid=%d", uid)

	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("query(%s), error(%s)", sql, err)
		return nil, nil, err
	}

	defer rows.Close()

	groups := &pb.GroupListRsp{Uid: uid, GroupList: nil}

	gMap := make(map[int32]string, 0)
	for rows.Next() {
		group := &pb.GroupInfo{}
		err = rows.Scan(&group.Gid, &group.GroupName)
		if err != nil {
			return nil, nil, err
		}
		group.IsExist = true
		gMap[group.Gid] = group.GroupName
		groups.GroupList = append(groups.GroupList, group)
	}

	return groups, &gMap, nil
}

func SearchGroup(target string, db *sql.DB) (*pb.GroupListRsp, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	rows, err := db.Query("SELECT id, group_name FROM user_group WHERE group_name LIKE ?", "%"+target+"%")
	if err != nil {
		log.Printf("query error : %v\n", err)
	}

	defer rows.Close()

	groups := &pb.GroupListRsp{GroupList: nil}

	for rows.Next() {
		group := new(pb.GroupInfo)
		err = rows.Scan(&group.Gid, &group.GroupName)
		if err != nil {
			return nil, err
		}
    	group.IsExist = false
		groups.GroupList = append(groups.GroupList, group)
	}

	return groups, nil
}

// 查找当前管理员能管理的群组
func GetGruopMembers(gid int32, db *sql.DB) (*pb.GrpMemberList, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	sql := fmt.Sprintf("SELECT u.id, u.name, u.user_type, gm.role_type "+
		"FROM user AS u RIGHT JOIN group_member AS gm ON gm.uid=u.id WHERE gm.gid=%d AND gm.stat=1", gid)

	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("query(%s), error(%s)\n", sql, err)
		return nil, err
	}

	defer rows.Close()

	grpMems := new(pb.GrpMemberList)
	grpMems.Gid = gid

	for rows.Next() {
		gm := new(pb.UserRecord)
		err = rows.Scan(&gm.Uid, &gm.Name, &gm.UserType, &gm.GrpRole)
		if err != nil {
			return nil, err
		}

		grpMems.UsrList = append(grpMems.UsrList, gm)
	}

	return grpMems, nil
}

// 查看
func SelectGroupsByAccountId(aid int) ([]*model.GroupInfo, error) {
	stmtOut, err := dbConn.Prepare("SELECT id, group_name, stat, create_time FROM user_group WHERE account_id = ? AND stat = 1")
	if err != nil {
		return nil, err
	}

	var res []*model.GroupInfo

	rows, err := stmtOut.Query(aid)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var gid int
		var groupName, status, cTime string
		if err := rows.Scan(&gid, &groupName, &status, &cTime); err != nil {
			return res, err
		}

		g := &model.GroupInfo{Id: gid, GroupName: groupName, AccountId: aid, Status: status, CTime: cTime}
		res = append(res, g)
	}

	defer func() {
		if err := stmtOut.Close(); err != nil {
			log.Println("statement close fail.")
		}
	}()
	return res, nil
}

// 创建组
func CreateGroup(gl *model.GroupList, userType int) (int64, error) {
	tx, err := dbConn.Begin()
	if err != nil {
		return -1, err
	}
	stmtInsG := "INSERT INTO user_group (group_name, account_id, user_type) VALUES (?, ?, ?)"

	insGroupRes, err := tx.Exec(stmtInsG, gl.GroupInfo.GroupName, gl.GroupInfo.AccountId, userType)
	if err != nil {
		log.Println("Insert Group error : ", err)
		return -1, err
	}

	var (
		groupAff, groupId, groupDeviceAff int64
	)

	if insGroupRes != nil {
		groupAff, _ = insGroupRes.RowsAffected()
		groupId, _ = insGroupRes.LastInsertId()
	}

	var ib = dbs.NewInsertBuilder()
	ib.Table("group_member")
	ib.Columns("gid", "uid", "role_type")
	// 如果是1就是web用户 range 每个设备的id
	if userType == 1 {
		for _, v := range gl.DeviceInfo {
			ib.Values(groupId, (v.(map[string]interface{}))["id"], 0)
		}
		ib.Values(groupId, gl.GroupInfo.AccountId, 1)
	} else {
		for _, v := range gl.DeviceIds {
			if v !=  gl.GroupInfo.AccountId {
				ib.Values(groupId, v, 0)
			}
		}
		ib.Values(groupId, gl.GroupInfo.AccountId, 1)  // 默认accountId属性作为group_member的群主，TODO 会有歧义，就是app用户创建的群组，调度员能否可见。
	}

	stmtInsGD, value, err := ib.ToSQL()
	if err != nil {
		log.Println("Error in ib ToSQL", err)
		return -1, err
	}

	insGroupDeviceRes, err := tx.Exec(stmtInsGD, value...)
	if err != nil {
		log.Println("Error in insert group device", err)
		return -1, err
	}

	if insGroupDeviceRes != nil {
		groupDeviceAff, _ = insGroupDeviceRes.RowsAffected()
	}

	log.Println(groupAff, groupDeviceAff, len(gl.DeviceIds), len(gl.DeviceInfo)+1)
	if (groupDeviceAff == int64(len(gl.DeviceInfo)+1) || groupDeviceAff == int64(len(gl.DeviceIds))) && groupAff == 1 {
		log.Println("commit")
		if err := tx.Commit(); err != nil {
			log.Println("tx commit")
			return -1, err
		}
	} else {
		log.Println("rollback")
		if err := tx.Rollback(); err != nil {
			return -1, err
		}
		return -1, errors.New("rollback")
	}
	return groupId, nil
}

// 查找群组
func SelectGroupByKey(key interface{}) (*model.GroupInfo, error) {
	var stmtOut *sql.Stmt
	var err error
	switch t := key.(type) {
	case int:
		stmtOut, err = dbConn.Prepare("SELECT id, group_name, account_id, stat, create_time FROM user_group WHERE id = ?")
	case string:
		stmtOut, err = dbConn.Prepare("SELECT id, group_name, account_id, stat, create_time FROM user_group WHERE group_name = ?")
	default:
		_ = t
	}

	if err != nil {
		return nil, err
	}

	var gid, accountId int
	var status, cTime, gName string
	if err := stmtOut.QueryRow(key).Scan(&gid, &gName, &accountId, &status, &cTime); err != nil {
		return nil, err
	}

	g := &model.GroupInfo{Id: gid, AccountId: accountId, GroupName: gName, Status: status, CTime: cTime}

	defer func() {
		if err := stmtOut.Close(); err != nil {
			log.Println("Statement close fail.")
		}
	}()
	return g, nil
}

// 更新群组 TODO
func UpdateGroup(info *model.GroupInfo, db *sql.DB) error {
	// 首先更新组，有更新group表，然后就去群里有几个设备，更新准备使用第三方库 目前web只用更新群的名字
	var ub = dbs.NewUpdateBuilder()
	ub.Table("user_group")
	ub.SET("group_name", info.GroupName)
	ub.Where("id = ? ", info.Id)
	if _, err := ub.Exec(db); err != nil {
		log.Println("update group name error：", err)
		return err
	}
	return nil
}

// 删除群组
func DeleteGroup(g *model.GroupInfo) error {
	tx, err := dbConn.Begin()
	if err != nil {
		return err
	}
	stmtUpd, err := tx.Prepare("DELETE FROM user_group WHERE id = ?")
	if err != nil {
		return err
	}

	if _, err := stmtUpd.Exec(g.Id); err != nil {
		return err
	}

	stmtUpdDG, err := tx.Prepare("DELETE FROM group_member WHERE gid = ?")
	if err != nil {
		return err
	}

	if _, err := stmtUpdDG.Exec(g.Id); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func CheckDuplicateGName(g *model.GroupInfo) (int, error) {
	stmtOut, err := dbConn.Prepare("SELECT count(id) FROM user_group WHERE group_name = ? AND account_id = ?")
	if err != nil {
		log.Printf("DB error :%s", err)
		return -1, err
	}

	var res int
	if err := stmtOut.QueryRow(g.GroupName, g.AccountId).Scan(&res); err != nil {
		log.Printf("err: %s", err)
		return -1, err
	}

	defer func() {
		if err := stmtOut.Close(); err != nil {
			log.Println("Statement close fail")
		}
	}()
	return res, nil
}

// 更新组成员
func UpdateGroupMember(gl *model.GroupList, userType int) (int64, error) {
	tx, err := dbConn.Begin()
	if err != nil {
		return -1, err
	}
	// 如果是更新，opsType就是true
	var updGMAff, groupDeviceAff int64
	stmtUpdDG := "DELETE FROM group_member WHERE gid = ?"
	if err != nil {
		return -1, err
	}

	//log.Println(gl.GroupInfo.Id)
	if updDGRes, err := tx.Exec(stmtUpdDG, gl.GroupInfo.Id); err != nil {
		return -1, err
	} else {
		if updDGRes != nil {
			updGMAff, err = updDGRes.RowsAffected()
		}
	}

	var ib = dbs.NewInsertBuilder()
	ib.Table("group_member")
	ib.Columns("gid", "uid", "role_type")
	// 如果是1就是web用户 range 每个设备的id
	if userType == 1 {
		for _, v := range gl.DeviceInfo {
			log.Printf("%T", v)
			log.Println("web test", (v.(map[string]interface{}))["id"])
			ib.Values(gl.GroupInfo.Id, (v.(map[string]interface{}))["id"], 0)
		}
		ib.Values(gl.GroupInfo.Id, gl.GroupInfo.AccountId, 1)
	} else {
		for index, v := range gl.DeviceIds {
			if index == 0 { // 默认把创建群组的切片第一个作为管理员
				ib.Values(gl.GroupInfo.Id, v, 1)
			}
			ib.Values(gl.GroupInfo.Id, v, 0)
		}
	}

	stmtInsGD, value, err := ib.ToSQL()
	if err != nil {
		log.Println("Error in ib ToSQL", err)
		return -1, err
	}

	insGroupDeviceRes, err := tx.Exec(stmtInsGD, value...)
	if err != nil {
		log.Println("Error in insert group device", err)
		return -1, err
	}

	if insGroupDeviceRes != nil {
		groupDeviceAff, _ = insGroupDeviceRes.RowsAffected()
	}

	log.Println(updGMAff, groupDeviceAff, len(gl.DeviceIds)+1, len(gl.DeviceInfo)+1)

	if (updGMAff == updGMAff) && (groupDeviceAff == int64(len(gl.DeviceInfo)+1) || groupDeviceAff == int64(len(gl.DeviceIds)+1)) {
		if err := tx.Commit(); err != nil {
			log.Println("tx commit")
			return -1, err
		}
	} else {
		log.Println("rollback")
		if err := tx.Rollback(); err != nil {
			return -1, err
		}
		return -1, errors.New("rollback")
	}

	return int64(gl.GroupInfo.Id), nil
}
