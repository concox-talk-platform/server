/**
* @Author: yanKoo
* @Date: 2019/3/11 11:16
* @Description:
 */
package account

import (
	"database/sql"
	"db"
	"github.com/smartwalle/dbs"
	"log"
	"model"
)

var dbConn = db.DBHandler
// 增加用户
func AddAccount(a *model.Account) error {
	stmtIns, err := dbConn.Prepare("INSERT INTO account (user_name, passwd) VALUES (?, ?)")
	if err != nil {
		return err
	}

	if _, err := stmtIns.Exec(a.Username, a.Pwd); err != nil {
		log.Println("AddAccountCredential Fail")
		return err
	}

	defer func() {
		if err := stmtIns.Close(); err != nil {
			log.Println("Statement close fail")
		}
	}()
	return nil
}

// 获取用户的密码
func GetAccountPwdByKey(key interface{}) (string, error) {
	var stmtOut *sql.Stmt
	var err error
	switch t := key.(type) {
	case int:
		stmtOut, err = dbConn.Prepare("SELECT passwd FROM account WHERE id = ?")
	case string:
		stmtOut, err = dbConn.Prepare("SELECT passwd FROM account WHERE user_name = ?")
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
	stmtDel, err := dbConn.Prepare("DELETE FROM account WHERE user_name = ? AND password = ?")
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
	stmtOut, err := dbConn.Prepare("SELECT count(id) FROM account WHERE user_name = ?")
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
		stmtOut, stmtErr = dbConn.Prepare(`SELECT id, pid, user_name, nick_name, passwd, email, phone, remark, address, privilege_id, 
       									  role_type, stat, last_login_time, create_time, change_time 
								FROM account WHERE id = ?`)
	case string:
		stmtOut, stmtErr = dbConn.Prepare(`SELECT id, pid, user_name, nick_name, passwd, email, phone, remark, address, privilege_id, 
       									  role_type, stat, last_login_time, create_time, change_time 
								FROM account WHERE user_name = ?`)
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
	err := stmtOut.QueryRow(key).Scan(&id, &pid, &username, &nickname, &pwd, &email, &phone, &remark, &address, &privilegeId, &roleId, &stat, &llTime, &cTime, &changeTime)
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
	var ub = dbs.NewUpdateBuilder()
	ub.Table("account")
	if a.Remark != "" {
		ub.SET("remark", a.Remark)
	}
	if a.Address != "" {
		ub.SET("address", a.Address)
	}
	if a.NickName != "" {
		ub.SET("nick_name", a.NickName)
	}
	if a.Email != "" {
		ub.SET("email", a.Email)
	}
	if a.Phone != "" {
		ub.SET("phone", &a.Phone)
	}
	ub.Where("id = ? ", a.Id)

	if _, err := ub.Exec(dbConn); err != nil {
		return err
	}

	return nil
}

// 更新密码
func UpdateAccountPwd(pwd string, id int) error {
	stmtUpd, err := dbConn.Prepare("UPDATE account SET passwd = ? WHERE id = ?")
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
	stmtOut, err := dbConn.Prepare("SELECT id, user_name FROM account WHERE pid = ?")
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
