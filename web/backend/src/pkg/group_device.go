/**
* @Author: yanKoo
* @Date: 2019/3/18 11:34
* @Description:
 */
package pkg

import (
	"database/sql"
	"github.com/concox-talk-platform/server/web/backend/src/model"
	"log"
)

func SelectDevicesByGroupId(gid int) ([]*model.Device, error) {
	stmtOut, err := dbConn.Prepare(`SELECT id, imei, user_name, user_passwd, account_id, stat, active_stat, bind_stat, create_time, last_login_time, change_time 
									FROM device WHERE id IN (SELECT device_id FROM group_device WHERE group_id = ?);`)
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
		var iMei, userName, pwd, status, bindStatus, aStatus, cTime, llTime, changeTime sql.NullString
		if err := rows.Scan(&id, &iMei, &userName, &pwd, &accountId, &status, &aStatus, &bindStatus, &cTime, &llTime, &changeTime); err != nil {
			return res, err
		}

		d := &model.Device{
			Id: id, IMei: iMei,
			UserName: userName, PassWord:pwd,
			AccountId:accountId,
			Status:    status, ActiveStatus: aStatus, BindStatus: bindStatus,
			CrateTime: cTime, LLTime:llTime, ChangeTime:changeTime,
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