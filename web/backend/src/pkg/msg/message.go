/**
 * Copyright (c) 2019. All rights reserved.
 * Deal with the messages from users
 * Author: tesion
 * Data: April 2nd 2019
 */
package msg

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
	pb "api/talk_cloud"
)

type MsgType uint8
type MsgStat uint8

const (
	MSG_UNREAD MsgStat = iota
	MSG_READ
)

const (
	PLAIN_TEXT MsgType = iota	// 普通文本
	IMAGE						// 图片
	AUDIO						// 音频
	VIDEO						// 视频
	LOCATION					// 位置
)

type MsgData struct {
	Id int64
	Uid int64
	SenderId int64
	Type MsgType
	Content string
	CreateTime int64
}

type LocationData struct {
	Id int64
	Uid int64
	Lat float64
	Lng float64
	TimeStamp int64
}

func NewMsg() *MsgData {
	return new(MsgData)
}

func AddMsg(msg *MsgData, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	timeStr := time.Unix(msg.CreateTime, 0).Format("2006-01-02 15:04:05")
	sql := fmt.Sprintf("INSERT INTO message(uid, sender_id, msg_type, content, create_time) " +
		"VALUES(%d,%d,%d,'%s','%s')", msg.Uid, msg.SenderId, msg.Type, msg.Content, timeStr)

	_, err := db.Query(sql)
	if err != nil {
		log.Printf("query(%s), error(%s)\n", sql, err)
		return err
	}

	return nil
}

func AddMultiMsg(req *pb.MsgNewReq, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	sz := len(req.Uids)
	if 0 == sz {
		log.Printf("empty uids\n")
		return nil
	}

	timeStr := time.Unix(int64(req.CreateTime), 0).Format("2006-01-02 15:04:05")
	sql := fmt.Sprintf("INSERT INTO message(uid, sender_id, msg_type, content, create_time) VALUES")

	var elem string
	for i := 0; i < sz; i++ {
		if 0 != i {
			sql += ","
		}
		elem = fmt.Sprintf("(%d,%d,%d,'%s','%s')", req.Uids[i], req.SenderId, req.MsgType, req.Content, timeStr)
		sql += elem
	}

	_, err := db.Query(sql)
	if err != nil {
		log.Printf("query(%s), error(%s)\n", sql, err)
		return err
	}

	return nil
}

func GetMsg(uid int32, stat pb.MsgStat, db *sql.DB) ([]*pb.MsgData, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	sql := fmt.Sprintf("SELECT id, sender_id, msg_type, content, UNIX_TIMESTAMP(create_time) " +
		"FROM message WHERE uid=%d AND stat=%d", uid, stat)
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("query(%s), error(%s)\n", sql, err)
		return nil, err
	}

	defer rows.Close()

	data := make([]*pb.MsgData, 0)

	for rows.Next() {
		msg := new(pb.MsgData)
		err = rows.Scan(&msg.MsgId, &msg.SenderId, &msg.MsgType, &msg.Content, &msg.CreateTime)
		if err != nil {
			log.Printf("Scan message error: %s\n", err)
			continue
		}

		data = append(data, msg)
	}

	return data, nil
}

func SetMsgStat(msgID int64, stat MsgStat, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	sql := fmt.Sprintf("UPDATE message SET stat=%d WHERE id=%d", stat, msgID)
	_, err := db.Query(sql)
	if err != nil {
		log.Printf("query(%s), error(%s)\n", sql, err)
		return err
	}

	return nil
}

func SetMultiMsgStat(msgID []int32, stat pb.MsgStat, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	sz := len(msgID)
	if 0 == sz {
		log.Printf("empty msg id list")
		return nil
	}

	sql := fmt.Sprintf("UPDATE message SET stat=%d WHERE id in (", stat)
	for i := 0; i < sz; i++ {
		if 0 != i {
			sql += ","
		}

		sql += strconv.FormatInt(int64(msgID[i]), 10)
	}

	sql += ")"

	_, err := db.Query(sql)
	if err != nil {
		log.Printf("query(%s), error(%s)\n", sql, err)
		return err
	}

	return nil
}

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

	sql += ")"

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