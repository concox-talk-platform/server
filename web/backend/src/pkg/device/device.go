/**
* @Author: yanKoo
* @Date: 2019/3/16 15:29
* @Description:
 */
package device

import (
	"database/sql"
	"github.com/smartwalle/dbs"
	"log"
	"server/web/backend/src/model"
	"time"
)

// 增加设备
func AddDevice(d *model.Device) error {
	var stmtInsB = dbs.NewInsertBuilder()

	stmtInsB.Table("device")
	stmtInsB.SET("imei", d.IMei)
	stmtInsB.SET("user_name", d.UserName)
	stmtInsB.SET("user_passwd", d.PassWord)
	stmtInsB.SET("account_id", d.AccountId)

	t := time.Now()
	ctime := t.Format("2006-1-2 15:04:05")
	stmtInsB.SET("last_login_time", ctime)
	stmtInsB.SET("change_time", ctime)

	if d.Status.Valid {
		stmtInsB.SET("stat", d.Status)
	}
	if d.ActiveStatus.Valid {
		stmtInsB.SET("active_stat", d.ActiveStatus)
	}
	if d.BindStatus.Valid {
		stmtInsB.SET("bind_stat", d.BindStatus)
	}

	if _, err := stmtInsB.Exec(dbConn); err != nil {
		return err
	}

	return nil
}

// 查找设备
func SelectDeviceByAccountId(aid int) (interface{}, error) {
	var stmtOut *sql.Stmt

	stmtOut, err = dbConn.Prepare("SELECT id, imei, user_name, user_passwd, account_id, stat, active_stat, bind_stat, create_time, last_login_time, change_time FROM device WHERE account_id = ?")

	if err != nil {
		log.Printf("%s", err)
		return "", err
	}
	var res []*model.Device

	rows, err := stmtOut.Query(aid)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id, accountId                                          int
			userName, pwd, iMei                                    string
			status, bindStatus, aStatus, cTime, llTime, changeTime sql.NullString
		)
		if err := rows.Scan(&id, &iMei, &userName, &pwd, &accountId, &status, &aStatus, &bindStatus, &cTime, &llTime, &changeTime); err != nil {
			return res, err
		}

		d := &model.Device{
			Id: id, IMei: iMei,
			UserName: userName, PassWord: pwd,
			AccountId: accountId,
			Status:    status, ActiveStatus: aStatus, BindStatus: bindStatus,
			CrateTime: cTime, LLTime: llTime, ChangeTime: changeTime,
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

// 通过关键词查找用户名
func SelectDeviceByKey(key interface{}) (*model.Device, error) {
	var stmtOut *sql.Stmt
	switch t := key.(type) {
	case int:
		stmtOut, err = dbConn.Prepare("SELECT id, imei, user_name, user_passwd, account_id, stat, active_stat, bind_stat, create_time, last_login_time, change_time FROM device WHERE account_id = ?")
	case string:
		stmtOut, err = dbConn.Prepare("SELECT id, imei, user_name, user_passwd, account_id, stat, active_stat, bind_stat, create_time, last_login_time, change_time FROM device WHERE user_name = ?")
	default:
		_ = t
		return nil, err
	}
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	var (
		id, accountId                                          int
		userName, pwd, iMei                                    string
		status, bindStatus, aStatus, cTime, llTime, changeTime sql.NullString
	)
	err := stmtOut.QueryRow(key).Scan(&id, &iMei, &userName, &pwd, &accountId, &status, &aStatus, &bindStatus, &cTime, &llTime, &changeTime)
	if err != nil{
		return nil, err
	}

	res := &model.Device{
		Id:id,
		IMei:iMei,
		UserName:userName,
		PassWord:pwd,
		AccountId:accountId,
		Status:status,
		BindStatus:bindStatus,
		ActiveStatus:aStatus,
		CrateTime:cTime,
		LLTime:llTime,
		ChangeTime:changeTime,
	}

	defer func() {
		if err := stmtOut.Close(); err != nil {
			log.Println("Statement close fail")
		}
	}()
	return res, nil
}

func GetAdviceByName(userName string) (int, error) {
	stmtOut, err := dbConn.Prepare("SELECT count(id) FROM device WHERE user_name = ?")
	if err != nil {
		log.Printf("DB error :%s", err)
		return -1, err
	}

	var res int
	if err := stmtOut.QueryRow(userName).Scan(&res); err != nil {
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
