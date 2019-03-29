/**
* @Author: yanKoo
* @Date: 2019/3/18 11:34
* @Description:
 */
package group_device

import (
	"database/sql"
	"db"
	"log"
	"server/web/backend/src/model"
)

var dbConn = db.DBHandler

func SelectDevicesByGroupId(gid int) ([]*model.Device, error) {
	stmtOut, err := dbConn.Prepare(`SELECT id, imei, name, passwd, cid, create_time, last_login_time, change_time 
									FROM user WHERE id IN (SELECT device_id FROM group_device WHERE group_id = ?);`)
	if err != nil {
		return nil, err
	}

	var res []*model.Device
	rows, err := stmtOut.Query(gid)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id, accountId int
		var iMei, userName, pwd string
		var cTime, llTime, changeTime sql.NullString
		if err := rows.Scan(&id, &iMei, &userName, &pwd, &accountId, &cTime, &llTime, &changeTime); err != nil {
			return res, err
		}

		d := &model.Device{
			Id: id, IMei: iMei,
			UserName: userName, PassWord: pwd,
			AccountId: accountId,
			//Status:    status, ActiveStatus: aStatus, BindStatus: bindStatus,
			CreateTime: cTime, LLTime: llTime, ChangeTime: changeTime,
		}
		res = append(res, d)
	}

	defer func() {
		if err := stmtOut.Close(); err != nil {
			log.Println("Statement close fail")
		}
	}()
	return res, nil
}

func SelectDeviceIdsByGroupId(gid int) ([]int, error) {
	stmtOut, err := dbConn.Prepare("SELECT device_id FROM group_device WHERE group_id = ?")
	if err != nil {
		return nil, err
	}

	var res []int
	rows, err := stmtOut.Query(gid)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return res, err
		}
		res = append(res, id)
	}

	defer func() {
		if err := stmtOut.Close(); err != nil {
			log.Println("statement close fail.")
		}
	}()
	return res, nil
}
