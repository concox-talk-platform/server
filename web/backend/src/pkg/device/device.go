/**
* @Author: yanKoo
* @Date: 2019/3/16 15:29
* @Description:
 */
package device

import (
	"db"
	"github.com/smartwalle/dbs"
	"log"
	"model"
	"time"
)

var dbConn = db.DBHandler

// 批量导入设备
func ImportDevice(u []*model.User) error {
	tx, err := dbConn.Begin()
	if err != nil {
		return err
	}
	var ib = dbs.NewInsertBuilder()

	ib.Table("user")
	ib.Columns("imei", "name", "passwd", "cid", "pid", "nick_name", "user_type", "last_login_time", "create_time")
	for _, v := range u {
		t := time.Now()
		ctime := t.Format("2006-1-2 15:04:05")
		ib.Values(v.IMei, v.UserName, v.PassWord, v.AccountId, v.ParentId, v.UserName, 1, ctime, ctime)
	}

	stmtIns, values, err := ib.ToSQL()
	res, err := tx.Exec(stmtIns, values...)
	if err != nil {
		return err
	}
	var resAff int64
	if res != nil {
		resAff, _ = res.RowsAffected()
	}

	if resAff == int64(len(u)) {
		_ = tx.Commit()
	} else {
		log.Println("rollback")
		if err := tx.Rollback(); err != nil {
			return err
		}
	}
	return nil
}

// 批量转移设备
func MultiUpdateDevice(accountDevices *model.AccountDeviceTransReq) error {
	var ub = dbs.NewUpdateBuilder()
	ub.Table("user")
	ub.SET("cid", accountDevices.Receiver.AccountId)
	dImeiArr := make([]string, 0)
	log.Println(accountDevices.Devices)
	for _, v := range accountDevices.Devices {
		dImeiArr = append(dImeiArr, v.IMei)
	}
	ub.Where(dbs.IN("imei", dImeiArr))

	if _, err := ub.Exec(dbConn); err != nil {
		log.Println("multi update device error :", err)
		return err
	}
	return nil
}
