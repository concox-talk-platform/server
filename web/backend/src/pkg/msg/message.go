/**
 * Copyright (c) 2019. All rights reserved.
 * Deal with the messages from users
 * Author: tesion
 * Data: April 2nd 2019
 */
package msg

import (
	pb "api/talk_cloud"
	"database/sql"
	"fmt"
	"github.com/smartwalle/dbs"
	"log"
	"strconv"
	"time"
)

type MsgType uint8

const (
	PLAIN_TEXT MsgType = iota // 普通文本
	IMAGE                     // 图片
	AUDIO                     // 音频
	VIDEO                     // 视频
	LOCATION                  // 位置
)

type MsgData struct {
	Id         int64
	Uid        int64
	SenderId   int64
	Type       MsgType
	Content    string
	CreateTime int64
}

type LocationData struct {
	Id        int64
	Uid       int64
	Lat       float64
	Lng       float64
	TimeStamp int64
}

func NewMsg() *MsgData {
	return new(MsgData)
}

func AddMsg(msg *pb.ImMsgReqData, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	timeStr := time.Now().Format("2006-01-02 15:04:05")
	sql := fmt.Sprintf("INSERT INTO message(sender_id, receiver_id, receiver_type, msg_type, content, create_time) "+
		"VALUES(%d,%d,%d,'%d','%s','%s')", msg.Id, msg.ReceiverId, msg.ReceiverType, msg.MsgType, msg.ResourcePath, timeStr)

	_, err := db.Query(sql)
	if err != nil {
		log.Printf("query(%s), error(%s)\n", sql, err)
		return err
	}

	return nil
}

func AddMultiMsg(msg *pb.ImMsgReqData, receiverId []int32, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	ib := dbs.NewInsertBuilder()
	ib.Table("message").Columns("sender_id", "receiver_id", "receiver_type", "gid", "msg_type", "content", "create_time")
	for _, v := range receiverId {
		ib.Values(msg.Id, v, msg.ReceiverType, msg.ReceiverId, msg.MsgType, msg.ResourcePath,
			time.Now().Format("2006-01-02 15:04:05"))
	}
	stmtIns, values, err := ib.ToSQL()
	if err != nil {
		return err
	}
	if _, err := db.Exec(stmtIns, values...); err != nil {
		log.Printf("Add Multi Msg error(%s)\n", err)
		return err
	}

	return nil
}

func GetMsg(uid int32, stat int32, db *sql.DB) ([]*pb.ImMsgReqData, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	sql := fmt.Sprintf(`SELECT sender_id, receiver_id, receiver_type, msg_type, content 
		FROM message WHERE receiver_id=%d AND stat=%d`, uid, stat)
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("query(%s), error(%s)\n", sql, err)
		return nil, err
	}

	defer rows.Close()

	data := make([]*pb.ImMsgReqData, 0)

	for rows.Next() {
		msg := new(pb.ImMsgReqData)
		err = rows.Scan(&msg.Id, &msg.ReceiverId, &msg.ReceiverType, &msg.MsgType, &msg.ResourcePath)
		if err != nil {
			log.Printf("Scan message error: %s\n", err)
			continue
		}

		data = append(data, msg)
	}

	return data, nil
}

func SetMsgStat(msgID int32, stat int, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	sql := fmt.Sprintf("UPDATE message SET stat=%d WHERE receiver_id=%d", stat, msgID)
	_, err := db.Query(sql)
	if err != nil {
		log.Printf("query(%s), error(%s)\n", sql, err)
		return err
	}

	return nil
}

//func SetMultiMsgStat(msgID []int32, stat pb.MsgStat, db *sql.DB) error {
//	if db == nil {
//		return fmt.Errorf("db is nil")
//	}
//
//	sz := len(msgID)
//	if 0 == sz {
//		log.Printf("empty msg id list")
//		return nil
//	}
//
//	sql := fmt.Sprintf("UPDATE message SET stat=%d WHERE id in (", stat)
//	for i := 0; i < sz; i++ {
//		if 0 != i {
//			sql += ","
//		}
//
//		sql += strconv.FormatInt(int64(msgID[i]), 10)
//	}
//
//	//sql += ")"
//
//	_, err := db.Query(sql)
//	if err != nil {
//		log.Printf("query(%s), error(%s)\n", sql, err)
//		return err
//	}
//
//	return nil
//}

func DeleteMsgByID(msgID int64, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	sql := fmt.Sprintf("DELETE FROM message WHERE id=%d", msgID)
	_, err := db.Query(sql)
	if err != nil {
		log.Printf("query(%s), error(%s)\n", sql, err)
		return err
	}

	return nil
}

func DeleteMsg(msgID []int32, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	sz := len(msgID)
	if 0 == sz {
		return nil
	}

	sql := fmt.Sprintf("DELETE FROM message WHERE id in (")
	for i := 0; i < sz; i++ {
		if 0 != i {
			sql += ","
		}

		sql += strconv.FormatInt(int64(msgID[i]), 10)
	}

	//sql += ")"

	_, err := db.Query(sql)
	if err != nil {
		log.Printf("query(%s), error(%s)\n", sql, err)
		return err
	}

	return nil
}

func DeleteMsgByUser(uid int64, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	sql := fmt.Sprintf("DELETE FROM message WHERE uid=%d", uid)
	_, err := db.Query(sql)
	if err != nil {
		log.Printf("query(%s), error(%s)\n", sql, err)
		return err
	}

	return nil
}
