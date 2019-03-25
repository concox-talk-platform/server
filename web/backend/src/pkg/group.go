/**
* @Author: yanKoo
* @Date: 2019/3/18 10:31
* @Description: 对群组表格的操作
 */
package pkg

import (
	"fmt"
	"server/web/backend/src/model"
	"github.com/smartwalle/dbs"
	"log"
	"strconv"
)

// 查看
func SelectGroupsByAccountId(aid int) ([]*model.GroupInfo, error) {
	stmtOut, err := dbConn.Prepare("SELECT id, group_name, stat, create_time FROM `group` WHERE account_id = ? AND stat = 1")
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
func CreateGroup(gl *model.GroupList) error {
	stmtInsG, err := dbConn.Prepare("INSERT INTO `group` (group_name, account_id) VALUES (?, ?)")
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
	stmtOut, err := dbConn.Prepare("SELECT id, stat, create_time FROM `group` WHERE group_name = ?")
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

	stmtUpd, err := tx.Prepare("UPDATE `group` SET stat = ? WHERE id = ?")
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
