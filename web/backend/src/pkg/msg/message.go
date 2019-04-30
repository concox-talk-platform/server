/**
 * Copyright (c) 2019. All rights reserved.
 * Deal with the messages from users
 * Author: tesion
 * Data: April 2nd 2019
 */
package msg

import (
	pb "api/talk_cloud"
	cfgComm "configs/common"
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

	IM_MSG_FROM_UPLOAD_RECEIVER_IS_USER  = 1 // APP和web通过httpClient上传的IM信息是发给个人
	IM_MSG_FROM_UPLOAD_RECEIVER_IS_GROUP = 2 // APP和web通过httpClient上传的IM信息是发给群组
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

	timeStr := time.Now().Format(cfgComm.TimeLayout)
	sql := fmt.Sprintf("INSERT INTO message(sender_id, s_name, send_time, recv_name, receiver_id, receiver_type, msg_type, content, create_time) "+
		"VALUES(%d,'%s','%s','%s',%d,%d,'%d','%s','%s')", msg.Id, msg.SenderName, msg.SendTime, msg.ReceiverName, msg.ReceiverId, msg.ReceiverType, msg.MsgType, msg.ResourcePath, timeStr)

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
	ib.Table("message").Columns("sender_id", "s_name","send_time", "recv_name", "receiver_id", "receiver_type", "gid", "msg_type", "content", "create_time")
	for _, v := range receiverId {
		ib.Values(msg.Id, msg.SenderName, msg.SendTime, msg.ReceiverName, v, msg.ReceiverType, msg.ReceiverId, msg.MsgType, msg.ResourcePath,
			time.Now().Format(cfgComm.TimeLayout))
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

	sql := fmt.Sprintf(`SELECT sender_id, s_name, send_time, recv_name, receiver_id, receiver_type, gid, msg_type, content
		FROM message WHERE receiver_id=%d AND stat=%d`, uid, stat)
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("query(%s), error(%s)\n", sql, err)
		return nil, err
	}

	defer rows.Close()

	data := make([]*pb.ImMsgReqData, 0)

	var (
		receiverId int32
		gid        int32
	)
	for rows.Next() {
		msg := new(pb.ImMsgReqData)
		err = rows.Scan(&msg.Id, &msg.SenderName, &msg.SendTime, &msg.ReceiverName, &receiverId, &msg.ReceiverType, &gid, &msg.MsgType, &msg.ResourcePath)
		if msg.ReceiverType == IM_MSG_FROM_UPLOAD_RECEIVER_IS_GROUP {
			msg.ReceiverId = gid
		}
		if msg.ReceiverType == IM_MSG_FROM_UPLOAD_RECEIVER_IS_USER {
			msg.ReceiverId = receiverId
		}

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
