/**
* @Author: yanKoo
* @Date: 2019/3/16 15:29
* @Description:
 */
package dbops

import (
	"github.com/concox-talk-platform/server/web/backend-api/defs"
	"log"
)

// 增加设备
func AddDevice(d *defs.Device) error {
	stmtIns, err := dbConn.Prepare(`INSERT INTO device (imei, alias_name, account_id, stat, active_stat, bind_stat) 
							VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	if _, err := stmtIns.Exec(d.IMei, d.AliasName, d.AccountId, d.Status, d.ActiveStatus, d.BindStatus); err != nil {
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
func SelectDeviceByAccountId(aid int) ([]*defs.Device, error) {
	stmtOut, err := dbConn.Prepare("SELECT id, imei, alias_name, stat, bind_stat, active_stat, create_time FROM device WHERE account_id = ?")
	if err != nil {
		return nil, err
	}

	var res []*defs.Device

	rows, err := stmtOut.Query(aid)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id int
		var iMei, aliasName, status, bindStatus, aStatus, cTime string
		if err := rows.Scan(&id, &iMei, &aliasName, &status, &bindStatus, &aStatus, &cTime); err != nil {
			return res, err
		}

		d := &defs.Device{
			Id: id, IMei: iMei,
			AliasName: aliasName,
			Status:    status, BindStatus: bindStatus,
			ActiveStatus: aStatus, CrateTime: cTime,
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
func Updatedevice(d *defs.Device) error {
	stmtUpd, err := dbConn.Prepare(`UPDATE device SET alias_name = ?, account_id = ?, stat = ?, active_stat = ?, bind_stat = ? WHERE id = ?`)
	if err != nil {
		return err
	}

	if _, err := stmtUpd.Exec(d.AliasName, d.AccountId, d.Status, d.ActiveStatus, d.BindStatus, d.Id); err != nil {
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
