/**
* @Author: yanKoo
* @Date: 2019/3/11 11:16
* @Description:
 */
package customer

import (
	"database/sql"
	"server/common/src/db"
	"log"
	"model"
	"strconv"
	"time"
)

var dbConn = db.DBHandler
// 增加用户
func AddAccount(a *model.Account) error {
	tx, err := dbConn.Begin();
	if err != nil {
		log.Println("事物开启失败")
	}

	t := time.Now()
	ctime := t.Format("2006-1-2 15:04:05")
	u := &model.User{
		IMei:       "1",
		UserName:   a.Username,
		NickName:   a.NickName,
		PassWord:   a.Pwd,
		UserType:   a.RoleId,
		ParentId:   strconv.FormatInt(int64(a.Pid), 10),
		AccountId:  "0",
		LLTime:     sql.NullString{String: ctime},
		CreateTime: sql.NullString{String: ctime},
	}

	stmtQuery := "INSERT INTO user (imei, name, passwd, cid, pid, nick_name, user_type, last_login_time, create_time)	VALUES (?, ?, ?,?, ?, ?, ?, ?, ?)"
	userRes, err := tx.Exec(stmtQuery, u.IMei, u.UserName, u.PassWord, u.AccountId, u.ParentId, u.NickName, u.UserType, u.LLTime, u.CreateTime)
	if err != nil {
		return err
	}

	uid, err := userRes.LastInsertId()
	if err != nil {
		log.Println("get insert AddUser Fail")
		return err
	}

	// customer
	cusRes, err := tx.Exec("INSERT INTO customer (uid, pid, email, phone, address, remark) VALUES (?, ?, ?, ?, ?, ?)",
		uid, a.Pid, a.Email, a.Phone, a.Address, a.Remark)
	if err != nil {
		return err
	}
	var affUser, affCus int64

	if userRes != nil {
		affUser, _ = userRes.RowsAffected()
	}
	if cusRes != nil {
		affCus, _ = cusRes.RowsAffected();
	}

	if affUser == 1 && affCus == 1 {
		// 提交事务
		if err := tx.Commit(); err != nil {
			return err
		}
	} else {
		// 回滚
		if err := tx.Rollback(); err != nil {
			return err
		}

	}
	return nil
}

// 获取用户的密码
func GetAccountPwdByKey(key interface{}) (string, error) {
	var stmtOut *sql.Stmt
	var err error
	switch t := key.(type) {
	case int:
		stmtOut, err = dbConn.Prepare("SELECT passwd FROM user WHERE id = ?")
	case string:
		stmtOut, err = dbConn.Prepare("SELECT passwd FROM user WHERE user_name = ?")
	default:
		_ = t
		return "", err
	}
	if err != nil {
		log.Printf("%s", err)
		return "", err
	}

	var pwd string
	err = stmtOut.QueryRow(key).Scan(&pwd)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	if err == sql.ErrNoRows {
		return "", err
	}

	defer func() {
		if err := stmtOut.Close(); err != nil {
			log.Println("Statement close fail")
		}
	}()
	return pwd, nil
}

// 删除用户
func DeleteAccount(loginName string, pwd string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM user WHERE user_name = ? AND password = ?")
	if err != nil {
		log.Printf("DeleteAccount error: %s", err)
		return err
	}
	if _, err := stmtDel.Exec(loginName, pwd); err != nil {
		return err
	}

	defer func() {
		if err := stmtDel.Close(); err != nil {
			log.Println("Statement close fail")
		}
	}()
	return nil
}

// 通过账户名获取账户数 注册查重
func GetAccountByName(userName string) (int, error) {
	stmtOut, err := dbConn.Prepare("SELECT count(id) FROM user WHERE name = ?")
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

// 获取用户
func GetAccount(key interface{}) (*model.Account, error) {
	var stmtOut *sql.Stmt
	var stmtErr error
	switch t := key.(type) {
	case int:
		stmtOut, stmtErr = dbConn.Prepare(`SELECT user.id, user.pid, user.name, user.nick_name, user.passwd, user_type, user.last_login_time, user.create_time, user.change_time, 
																email, phone, remark, address FROM user LEFT JOIN customer ON user.id = customer.uid WHERE user.id = ?`)
	case string:
		stmtOut, stmtErr = dbConn.Prepare(`SELECT user.id, user.pid, user.name, user.nick_name, user.passwd, user_type, user.last_login_time, user.create_time, user.change_time, 
																email, phone, remark, address FROM user LEFT JOIN customer ON user.id = customer.uid WHERE user.name = ?`)
	default:
		_ = t
	}

	if stmtErr != nil {
		log.Printf("%s", stmtErr)
		return nil, stmtErr
	}
	var (
		id          int
		pid         int
		username    string
		nickname    string
		pwd         string
		email       string
		phone       string
		remark      string
		address     string
		privilegeId int
		roleId      int
		stat        string
		llTime      string
		cTime       string
		changeTime  string
	)
	// 查询数据
	err := stmtOut.QueryRow(key).
		Scan(&id, &pid, &username, &nickname, &pwd, &roleId, &llTime, &cTime, &changeTime, &email, &phone, &remark, &address)

	if err != nil {
		log.Printf("err: %s", err)
		return nil, err
	}

	if err == sql.ErrNoRows {
		log.Println("no rows")
		return nil, nil
	}

	// 赋值给返回的结构体
	log.Println("get account : ", id, "  ", username, " ", privilegeId, " ", pwd, " ", cTime)
	res := &model.Account{
		Id:          id,
		Pid:         pid,
		Username:    username,
		NickName:    nickname,
		Pwd:         pwd,
		Email:       email,
		PrivilegeId: privilegeId,
		RoleId:      roleId,
		State:       stat,
		LlTime:      llTime,
		ChangeTime:  changeTime,
		CTime:       cTime,
		Phone:       phone,
		Address:     address,
		Remark:      address,
	}

	defer func() {
		if err := stmtOut.Close(); err != nil {
			log.Println("Statement close fail")
		}
	}()
	return res, nil
}

// 更新用户
func UpdateAccount(a *model.AccountUpdate) error {
	tx, err := dbConn.Begin()
	if err != nil {
		log.Println("事务开启失败")
		return err
	}

	userUpdStmt := "UPDATE user SET nick_name = ? WHERE id = ?"
	cusUpdStmt := "UPDATE customer SET remark = ?, address = ?, email = ?, phone = ? WHERE uid = ?"

	userRes, err := tx.Exec(userUpdStmt, a.NickName, a.Id)
	if err != nil {
		log.Println("update user error : ", err)
		return err
	}

	cusRes, err := tx.Exec(cusUpdStmt, a.Remark, a.Address, a.Email, a.Phone, a.Id)
	if err != nil {
		log.Println("update customer error : ", err)
		return err
	}

	var userAff, cusAff int64

	if userRes != nil {
		userAff, _ = userRes.RowsAffected()
	}

	if cusRes != nil {
		cusAff, _ = cusRes.RowsAffected()
	}

	if userAff == 1 && cusAff == 1 {
		 _ = tx.Commit()
	} else {
		_ = tx.Rollback()
	}

	return nil
}

// 更新密码
func UpdateAccountPwd(pwd string, id int) error {
	stmtUpd, err := dbConn.Prepare("UPDATE user SET passwd = ? WHERE id = ?")
	if err != nil {
		return err
	}

	if _, err := stmtUpd.Exec(pwd, id); err != nil {
		return err
	}

	return nil
}

// 查找下级目录
func SelectChildByPId(pId int) ([]*model.Account, error) {
	stmtOut, err := dbConn.Prepare("SELECT id, name FROM user WHERE pid = ?")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := stmtOut.Close(); err != nil {
			log.Println("statement close fail.")
		}
	}()

	rows, err := stmtOut.Query(pId)
	if err != nil {
		return nil, err
	}

	var res []*model.Account

	for rows.Next() {
		var id int
		var userName string
		if err := rows.Scan(&id, &userName); err != nil {
			return res, err
		}

		acc := &model.Account{Id: id, Pid: pId, Username: userName}
		res = append(res, acc)
	}

	return res, nil
}
