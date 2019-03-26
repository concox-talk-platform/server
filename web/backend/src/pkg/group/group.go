/**
 * Copyrights (c) 2019. All rights reserved.
 * Group handlers
 * Author: tesion
 * Date: March 26 2019
 */
package group

import (
    pb "api/talk_cloud"
    "database/sql"
    "fmt"
    "log"
)


func CreateGroup(uid int64, groupName string, db *sql.DB) error {
    if db == nil {
        return fmt.Errorf("db is nil")
    }
    
    res, err := db.Exec("INSERT INTO group(group_name) VALUES(?)", groupName)
    if err != nil {
        log.Printf("query error(%s)\n", err)
        return err
    }
    
    group_id, err := res.LastInsertId()
    if err != nil {
        log.Printf("get last insert id error: %s", err)
        return err
    }
    
    return AddGroupUser(uid, group_id, GROUP_MANAGER, db)
}

func AddGroupUser(uid, gid int64, userType RoleType,  db *sql.DB) error {
    if db == nil {
        return fmt.Errorf("db is nil")
    }
    
    sql := fmt.Sprintf("INSERT INTO group_device(group_id, device_id, role_type) VALUES(%d, %d, %d)", gid, uid, userType)
    rows, err := db.Query(sql)
    if err != nil {
        log.Printf("query(%s), error(%s)", sql, err)
        return err
    }
    
    defer rows.Close()
    
    return nil
}

func RemoveGroupUser(uid, gid uint64, db *sql.DB) error {
    if db == nil {
        return fmt.Errorf("db is nil")
    }

    sql := fmt.Sprintf("DELETE FROM group_device WHERE device_id=%d AND group_id=%d", uid, gid)
    _, err := db.Query(sql)
    if err != nil {
        log.Printf("query(%s), error(%s)", sql, err)
        return err
    }

    return nil
}

func RemoveGroup(gid uint64, db *sql.DB) error {
    if db == nil {
        return fmt.Errorf("db is nil")
    }

    sql := fmt.Sprintf("DELETE FROM group WHERE id=%d", gid)
    _, err := db.Query(sql)
    if err != nil {
        log.Printf("remove group(%d) error: %s\n", gid, err)
        return err
    }

    return nil
}

func ClearGroupUser(gid uint64, db *sql.DB) error {
    if db == nil {
        return fmt.Errorf("db is nil")
    }

    sql := fmt.Sprintf("DELETE FROM group_device WHERE group_id=%d", gid)
    _, err := db.Query(sql)
    if err != nil {
        log.Printf("clear gruop(%d) user error: %s\n", gid, err)
        return err
    }

    return nil
}

func GetGroupList(uid uint64, db *sql.DB) (*pb.GroupListRsp, error) {
    if db == nil {
        return nil, fmt.Errorf("db is nil")
    }

    sql := fmt.Sprintf("SELECT g.id, g.name " +
        "FROM group AS g RIGHT LEFT JOIN group_device AS gd " +
        "ON g.id=gd.group_id WHERE gd.device_id=%d", uid)

    rows, err := db.Query(sql)
    if err != nil {
        log.Printf("query(%s), error(%s)", sql, err)
        return nil, err
    }

    defer rows.Close()

    groups := &pb.GroupListRsp{Uid:uid, GroupList:nil}

    for rows.Next() {
        group := new(pb.GroupRecord)
        err = rows.Scan(&group.Gid, &group.GroupName)
        if err != nil {
            return nil, err
        }

        groups.GroupList = append(groups.GroupList, group)
    }

    return groups, nil
}

func SearchGroup(target string, db *sql.DB) (*pb.GroupListRsp, error) {
    if db == nil {
        return nil, fmt.Errorf("db is nil")
    }

    sql := fmt.Sprintf("SELECT id, group_name FROM group WHERE group_name LIKE '%%s%'", target)
    rows, err := db.Query(sql)
    if err != nil {
        log.Printf("query(%s), error(%s)\n", sql, err)
    }

    defer rows.Close()

    groups := &pb.GroupListRsp{GroupList:nil}

    for rows.Next() {
        group := new(pb.GroupRecord)
        err = rows.Scan(&group.Gid, &group.GroupName)
        if err != nil {
            return nil, err
        }

        groups.GroupList = append(groups.GroupList, group)
    }

    return groups, nil
}