package group

import (
    "database/sql"
    "fmt"
    "log"
)

func CreateGroup(uid int64, groupName string, db *sql.DB) error {
    if db == nil {
        return fmt.Errorf("db is nil")
    }
    
    res, err := db.Exec("INSERT INTO group(name) VALUES(?)", groupName)
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