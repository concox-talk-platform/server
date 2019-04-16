package user

import (
	pb "api/talk_cloud"
	"database/sql"
	"db"
	"errors"
	"github.com/smartwalle/dbs"
	"log"
	"model"
	"time"
)

var dbConn = db.DBHandler

// 增加设备
func AddUser(u *model.User, db ...interface{}) error {
	var stmtInsB = dbs.NewInsertBuilder()

	stmtInsB.Table("user")
	//stmtInsB.SET("id",u.Id)
	stmtInsB.SET("imei", u.IMei)
	stmtInsB.SET("name", u.UserName)
	stmtInsB.SET("passwd", u.PassWord)
	stmtInsB.SET("cid", u.AccountId)
	//stmtInsB.SET("pid", u.ParentId)
	stmtInsB.SET("nick_name", u.NickName) // 注册的时候默认把username当做昵称
	stmtInsB.SET("user_type", 1)
	t := time.Now()
	ctime := t.Format("2006-1-2 15:04:05")
	stmtInsB.SET("last_login_time", ctime)
	stmtInsB.SET("create_time", ctime)

	if _, err := stmtInsB.Exec(dbConn); err != nil {
		return err
	}

	return nil
}

// 用过用户名查重，用在app GRpc注册
func GetUserByName(userName string) (int, error) {
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

// 通过关键词查找用户名
func SelectUserByKey(key interface{}) (*model.User, error) {
	var stmtOut *sql.Stmt
	var err error
	switch t := key.(type) {
	case int:
		stmtOut, err = dbConn.Prepare("SELECT id, name, nick_name, passwd, imei, user_type, pid, cid, lock_gid, create_time, last_login_time, change_time FROM `user` WHERE id = ?")
	case string:
		stmtOut, err = dbConn.Prepare("SELECT id, name, nick_name, passwd, imei, user_type, pid, cid, lock_gid, create_time, last_login_time, change_time  FROM `user` WHERE name = ?")
	default:
		_ = t
		return nil, err
	}
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	var (
		id, userType, cId, lockGId                                    int
		pId, userName, nickName, pwd, iMei, cTime, llTime, changeTime string
	)
	err = stmtOut.QueryRow(key).Scan(&id, &userName, &nickName, &pwd, &iMei, &userType, &pId, &cId, &lockGId, &cTime, &llTime, &changeTime)
	if err != nil {
		return nil, err
	}

	res := &model.User{
		Id:          id,
		IMei:        iMei,
		UserName:    userName,
		PassWord:    pwd,
		NickName:    nickName,
		UserType:    userType,
		ParentId:    pId,
		AccountId:   cId,
		LockGroupId: lockGId,
		CreateTime:  cTime,
		LLTime:      llTime,
		ChangeTime:  changeTime,
	}

	defer func() {
		if err := stmtOut.Close(); err != nil {
			log.Println("Statement close fail")
		}
	}()
	return res, nil
}

// 查找设备
func SelectUserByAccountId(aid int) (interface{}, error) {
	var stmtOut *sql.Stmt
	var err error
	stmtOut, err = dbConn.Prepare("SELECT id, imei, name, passwd, cid, create_time, last_login_time, change_time FROM user WHERE cid = ?")

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
			id, accountId             int
			userName, pwd, iMei       string
			cTime, llTime, changeTime string
		)
		if err := rows.Scan(&id, &iMei, &userName, &pwd, &accountId, &cTime, &llTime, &changeTime); err != nil {
			return res, err
		}

		d := &model.Device{
			Id: id,
			IMei: iMei,
			UserName:  userName, //PassWord: pwd,
			AccountId: accountId,
			CreateTime: cTime,
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

// 设置用户锁定默认组
func SetLockGroupId(req *pb.SetLockGroupIdReq, db *sql.DB) error {
	if db == nil {
		return errors.New("set Lock group Id error, db is nil")
	}

	updStmt, err := db.Prepare("UPDATE`user` SET lock_gid = ? WHERE id = ?")
	if err != nil {
		return errors.New("set Lock group Id error, updStmt error " + err.Error())
	}

	_, err = updStmt.Exec(req.GId, req.UId)
	if err != nil {
		return errors.New("set Lock group Id error, updStmt.Exec error " + err.Error())
	}

	return nil
}
