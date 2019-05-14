/**
* @Author: yanKoo
* @Date: 2019/3/18 11:34
* @Description:
 */
package group_member

import (
	"server/web-api/log"
	"server/web-api/db"
	"server/web-api/model"
)

var dbConn = db.DBHandler

func SelectDevicesByGroupId(gid int) ([]*model.User, error) {
	stmtOut, err := dbConn.Prepare(`SELECT id, imei, name, passwd, user_type, cid, create_time, last_login_time, change_time 
									FROM user WHERE id IN (SELECT uid FROM group_member WHERE gid = ?) AND user_type = 1`)
	if err != nil {
		return nil, err
	}

	var res []*model.User
	rows, err := stmtOut.Query(gid)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id, accountId, userType int
		var iMei, userName, pwd string
		var cTime, llTime, changeTime string
		if err := rows.Scan(&id, &iMei, &userName, &pwd, &userType, &accountId, &cTime, &llTime, &changeTime); err != nil {
			return res, err
		}

		d := &model.User{
			Id: id, IMei: iMei,
			UserName:  userName, //PassWord: pwd,
			AccountId: accountId,
			UserType:  userType,
			//Status:    status, ActiveStatus: aStatus, BindStatus: bindStatus,
			CreateTime: cTime, LLTime: llTime, ChangeTime: changeTime,
		}
		res = append(res, d)
	}

	defer func() {
		if err := stmtOut.Close(); err != nil {
			log.Log.Println("Statement close fail")
		}
	}()
	return res, nil
}

func SelectDeviceIdsByGroupId(gid int) ([]int, error) {
	stmtOut, err := dbConn.Prepare("SELECT uid FROM group_member WHERE gid = ?")
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
			log.Log.Println("statement close fail.")
		}
	}()
	return res, nil
}
