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
	"fmt"
	"github.com/smartwalle/dbs"
	"log"
	"model"
	"server/common/src/db"
	"strconv"
)

var dbConn = db.DBHandler

func CreateGroup(uid int64, groupName string, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	res, err := db.Exec("INSERT INTO group(group_name) VALUES(?)", groupName)
	if err != nil {
		log.Printf("query error(%s)\n", err)
		return err
	}

	group_id, err := res.LastInsertId()
	if err != nil {
		log.Printf("get last insert id error: %s", err)
		return err
	}

	return AddGroupUser(uid, group_id, GROUP_MANAGER, db)
}

func AddGroupUser(uid, gid int64, userType RoleType, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	sql := fmt.Sprintf("INSERT INTO group_device(group_id, device_id, role_type) VALUES(%d, %d, %d)", gid, uid, userType)
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("query(%s), error(%s)", sql, err)
		return err
	}

	defer rows.Close()

	return nil
}

func RemoveGroupUser(uid, gid uint64, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	sql := fmt.Sprintf("DELETE FROM group_device WHERE device_id=%d AND group_id=%d", uid, gid)
	_, err := db.Query(sql)
	if err != nil {
		log.Printf("query(%s), error(%s)", sql, err)
		return err
	}

	return nil
}

func RemoveGroup(gid uint64, db *sql.DB) error {
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

func ClearGroupUser(gid uint64, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	sql := fmt.Sprintf("DELETE FROM group_device WHERE group_id=%d", gid)
	_, err := db.Query(sql)
	if err != nil {
		log.Printf("clear gruop(%d) user error: %s\n", gid, err)
		return err
	}

	return nil
}

func GetGroupList(uid uint64, db *sql.DB) (*pb.GroupListRsp, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	sql := fmt.Sprintf("SELECT g.id, g.name "+
		"FROM group AS g RIGHT LEFT JOIN group_device AS gd "+
		"ON g.id=gd.group_id WHERE gd.device_id=%d", uid)

	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("query(%s), error(%s)", sql, err)
		return nil, err
	}

	defer rows.Close()

	groups := &pb.GroupListRsp{Uid: uid, GroupList: nil}

	for rows.Next() {
		group := new(pb.GroupRecord)
		err = rows.Scan(&group.Gid, &group.GroupName)
		if err != nil {
			return nil, err
		}

		groups.GroupList = append(groups.GroupList, group)
	}

	return groups, nil
}

func SearchGroup(target string, db *sql.DB) (*pb.GroupListRsp, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	sql := fmt.Sprintf("SELECT id, group_name FROM group WHERE group_name LIKE '%%s%'", target)
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("query(%s), error(%s)\n", sql, err)
	}

	defer rows.Close()

	groups := &pb.GroupListRsp{GroupList: nil}

	for rows.Next() {
		group := new(pb.GroupRecord)
		err = rows.Scan(&group.Gid, &group.GroupName)
		if err != nil {
			return nil, err
		}

		groups.GroupList = append(groups.GroupList, group)
	}

	return groups, nil
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

		g := &model.GroupInfo{Id: gid, GroupName: groupName, Status: status, CTime: cTime}
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
func CreateGroupByWeb(gl *model.GroupList) error {
	stmtInsG, err := dbConn.Prepare("INSERT INTO user_group (group_name, account_id) VALUES (?, ?)")
	if err != nil {
		return err
	}

	if _, err := stmtInsG.Exec(gl.GroupInfo.GroupName, gl.GroupInfo.AccountId); err != nil {
		return err
	}
	defer func() {
		if err := stmtInsG.Close(); err != nil {
			log.Println("Statement close fail.")
		}
	}()

	tx, err := dbConn.Begin()
	log.Printf("group name :%s", gl.GroupInfo.GroupName)

	g, err := SelectGroupByGroupName(gl.GroupInfo.GroupName)
	if err != nil {
		log.Printf("查找群组失败")
		return err
	}

	var ib = dbs.NewInsertBuilder()
	ib.Table("group_device")
	ib.Columns("group_id", "device_id")
	for i := 0; i < len(gl.DeviceIds); i++ {
		ib.Values(g.Id, gl.DeviceIds[i])
	}

	if _, err := ib.Exec(dbConn); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// 查找群组
func SelectGroupByGroupName(GroupName string) (*model.GroupInfo, error) {
	stmtOut, err := dbConn.Prepare("SELECT id, stat, create_time FROM user_group WHERE group_name = ?")
	if err != nil {
		return nil, err
	}

	var gid int
	var status, cTime string
	if err := stmtOut.QueryRow(GroupName).Scan(&gid, &status, &cTime); err != nil {
		return nil, err
	}

	g := &model.GroupInfo{Id: gid, GroupName: GroupName, Status: status, CTime: cTime}

	defer func() {
		if err := stmtOut.Close(); err != nil {
			log.Println("Statement close fail.")
		}
	}()
	return g, nil
}

// 更新群组 TODO
func UpdateGroup() {
	// 首先更新组，有更新group表，然后就去群里有几个设备，更新准备使用第三方库
	var ub = dbs.NewUpdateBuilder()
	ub.Table("update")

	ub.SET("name", "新的名字")

	ub.Where("id = ? ", 1)
	ub.Limit(1)
	fmt.Println(ub.ToSQL())
}

// 删除群组
func DeleteGroup(g *model.GroupInfo) error {
	tx, err := dbConn.Begin()
	if err != nil {
		return err
	}
	stat, statErr := strconv.Atoi(g.Status)
	if statErr != nil {
		return statErr
	}

	stmtUpd, err := tx.Prepare("UPDATE user_group SET stat = ? WHERE id = ?")
	if err != nil {
		return err
	}

	if _, err := stmtUpd.Exec(stat, g.Id); err != nil {
		return err
	}

	stmtUpdDG, err := tx.Prepare("UPDATE group_device SET stat = ? WHERE group_id = ?")
	if err != nil {
		return err
	}

	if _, err := stmtUpdDG.Exec(stat, g.Id); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
