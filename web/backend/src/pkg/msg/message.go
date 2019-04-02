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
	"time"
)

type MsgType uint8
type MsgStat uint8

const (
	MSG_UNREAD MsgStat = iota
	MSG_READ
)

const (
	PLAIN_TEXT MsgType = iota
	AUDIO
	VEDIO
)

type MsgData struct {
	Id int64
	Uid int64
	SenderId int64
	Type MsgType
	Src string
	CreateTime int64
}

func NewMsg() *MsgData {
	return new(MsgData)
}

func AddMsg(msg *MsgData, db *sql.DB) error {
	if db == nil {
		fmt.Errorf("db is nil")
	}

	timeStr := time.Unix(msg.CreateTime, 0).Format("2006-01-02 15:04:05")
	sql := fmt.Sprintf("INSERT INTO message(uid, sender_id, msg_type, src, create_time) " +
		"VALUES(%d,%d,%d,'%s','%s')", msg.Uid, msg.SenderId, msg.Type, msg.Src, timeStr)

	_, err := db.Query(sql)
	if err != nil {
		log.Printf("query(%s), error(%s)\n", sql, err)
		return err
	}

	return nil
}

func GetMsg(uid int64, stat MsgStat, db *sql.DB) ([]*MsgData, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	sql := fmt.Sprintf("SELECT id, uid, sender_id, msg_type, src, UNIX_TIMESTAMP(create_time) " +
		"FROM message WHERE uid=%d AND stat=%d", uid, stat)
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("query(%s), error(%s)\n", sql, err)
		return nil, err
	}

	defer rows.Close()

	data := []*MsgData{}

	for rows.Next() {
		msg := NewMsg()
		err = rows.Scan(&msg.Id,&msg.Uid, &msg.SenderId, &msg.Type, &msg.Src, &msg.CreateTime)
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