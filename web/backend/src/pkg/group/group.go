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

func CreateGroup(uid, user_type int64, groupName string, db *sql.DB) error {
    if db == nil {
        return fmt.Errorf("db is nil")
    }
    
    res, err := db.Exec("INSERT INTO user_group(account_id,user_type,group_name) VALUES(?,?,?)", uid, user_type, groupName)
    if err != nil {
        log.Printf("query error(%s)\n", err)
        return err
    }
    
    group_id, err := res.LastInsertId()
    if err != nil {
        log.Printf("get last insert id error: %s", err)
        return err
    }
    
    return AddGroupMember(uid, group_id, GROUP_MANAGER, db)
}

func AddGroupMember(uid, gid int64, userType RoleType,  db *sql.DB) error {
    if db == nil {
        return fmt.Errorf("db is nil")
    }
    
    sql := fmt.Sprintf("INSERT INTO group_member(gid, uid, role_type) VALUES(%d, %d, %d)", gid, uid, userType)
    rows, err := db.Query(sql)
    if err != nil {
        log.Printf("query(%s), error(%s)", sql, err)
        return err
    }
    
    defer rows.Close()
    
    return nil
}

func RemoveGroupMember(uid, gid uint64, db *sql.DB) error {
    if db == nil {
        return fmt.Errorf("db is nil")
    }

    sql := fmt.Sprintf("DELETE FROM group_member WHERE uid=%d AND gid=%d", uid, gid)
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

    sql := fmt.Sprintf("DELETE FROM user_group WHERE id=%d", gid)
    _, err := db.Query(sql)
    if err != nil {
        log.Printf("remove group(%d) error: %s\n", gid, err)
        return err
    }

    return nil
}

func ClearGroupMember(gid uint64, db *sql.DB) error {
    if db == nil {
        return fmt.Errorf("db is nil")
    }

    sql := fmt.Sprintf("DELETE FROM group_member WHERE gid=%d", gid)
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

    sql := fmt.Sprintf("SELECT g.id, g.group_name " +
        "FROM user_group AS g RIGHT LEFT JOIN group_member AS gm " +
        "ON g.id=gm.group_id WHERE gm.uid=%d", uid)

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

    sql := fmt.Sprintf("SELECT id, group_name FROM user_group WHERE group_name LIKE '%%s%'", target)
    rows, err := db.Query(sql)
    if err != nil {
        log.Printf("query(%s), error(%s)\n", sql, err)
        return nil, err
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

func GetGruopMembers(gid uint64, db *sql.DB) (*pb.GrpMemberList, error){
    if db == nil {
        return nil, fmt.Errorf("db is nil")
    }

    sql := fmt.Sprintf("SELECT u.id, u.name, u.user_type, gm.role_type " +
        "FROM user AS u RIGHT JOIN group_member AS gm ON gm.uid=u.id WHERE gm.gid=%d AND gm.stat=1", gid)

    rows, err := db.Query(sql)
    if err != nil {
        log.Printf("query(%s), error(%s)\n", sql, err)
        return nil, err
    }

    defer rows.Close()

    grpMems := new(pb.GrpMemberList)
    grpMems.Gid = gid

    for rows.Next() {
        gm := new(pb.UserRecord)
        err = rows.Scan(&gm.Uid, &gm.Name, &gm.UserType, &gm.GrpRole)
        if err != nil {
            return nil, err
        }

        grpMems.UsrList = append(grpMems.UsrList, gm)
    }

    return grpMems, nil
}