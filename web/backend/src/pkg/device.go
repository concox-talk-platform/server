/**
* @Author: yanKoo
* @Date: 2019/3/16 15:29
* @Description:
 */
package pkg

import (
	"database/sql"
	"server/web/backend/src/model"
	"log"
)

// 增加设备
func AddDevice(d *model.Device) error {
	stmtIns, err := dbConn.Prepare(`INSERT INTO device (imei, alias_name, account_id, stat, active_stat, bind_stat) 
							VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	if _, err := stmtIns.Exec(d.IMei, d.UserName, d.AccountId, d.Status, d.ActiveStatus, d.BindStatus); err != nil {
		return err
	}

	defer func() {
		if err := stmtIns.Close(); err != nil {
			log.Println("Statement close fail")
		}
	}()
	return nil
}

// 查找设备
func SelectDeviceByAccountId(aid int) ([]*model.Device, error) {
	stmtOut, err := dbConn.Prepare("SELECT id, imei, user_name, user_passwd, account_id, stat, active_stat, bind_stat, create_time, last_login_time, change_time FROM device WHERE account_id = ?")
	if err != nil {
		return nil, err
	}

	var res []*model.Device

	rows, err := stmtOut.Query(aid)
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

// 更新设备
func Updatedevice(d *model.Device) error {
	stmtUpd, err := dbConn.Prepare(`UPDATE device SET alias_name = ?, account_id = ?, stat = ?, active_stat = ?, bind_stat = ? WHERE id = ?`)
	if err != nil {
		return err
	}

	if _, err := stmtUpd.Exec(d.UserName, d.AccountId, d.Status, d.ActiveStatus, d.BindStatus, d.Id); err != nil {
		return err
	}

	defer func() {
		if err := stmtUpd.Close(); err != nil {
			log.Println("Statement close fail")
		}
	}()
	return nil
}

// 删除设备
