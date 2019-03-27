/**
* @Author: yanKoo
* @Date: 2019/3/16 15:32
* @Description:
 */
package model

import "database/sql"

// request
// update pwd modal
type AccountPwd struct {
	Id         string `json:"id"`
	OldPwd     string `json:"old_pwd"`
	NewPwd     string `json:"new_pwd"`
	ConfirmPwd string `json:"confirm_pwd"`
}

// validate
type AccountValidate struct {
	Username string `valid:", between=6|15"`
	Pwd      string `valid:", between=6|15"`
}

//response
type Message struct {
	Result  bool   `json:"result"`
	Type    string `json:"type"`
	Message string `json:"message"`
	//Sticky string `json:"sticky"`
}

// 账户层级关系
type AccountClass struct {
	Id int `json:"id"`
	AccountName string `json:"account_name"`
	Children []*AccountClass `json:"children"`
}

// 返回账户及其账户下所有的组和
type AccountGroupsResp struct {
	Message     string             `json:"message"`
	AccountInfo *Account `json:"account_info"`
	GroupList   []*GroupList       `json:"group_list"`
	DeviceList  []*Device          `json:"device_list"`
}

type GroupList struct {
	DeviceIds  []int      `json:"device_ids"`
	DeviceInfo []*Device  `json:"device_infos"`
	GroupInfo  *GroupInfo `json:"group_info"`
}

// Data model
type SessionInfo struct {
	SessionID string `json:"session_id"`
	UserName  string `json:"user_name"`
	UserPwd   string `json:"user_pwd"`
	AccountId int    `json:"account_id"`
	TTL       string `json:"ttl"`
}

// device
type Device struct {
	Id           int            `json:"id"`
	IMei         string         `json:"imei"`
	UserName     string         `json:"user_name"`
	PassWord     string         `json:"password"`
	AccountId    int            `json:"account_id"`
	Status       sql.NullString `json:"status"`
	ActiveStatus sql.NullString `json:"active_status"`
	BindStatus   sql.NullString `json:"bind_status"`
	CrateTime    sql.NullString `json:"crate_time"`
	LLTime       sql.NullString `json:"last_login_time"`
	ChangeTime   sql.NullString `json:"change_time"`
}

type GroupInfo struct {
	Id        int    `json:"id"`
	GroupName string `json:"group_name"`
	AccountId int    `json:"account_id"`
	Status    string `json:"status"`
	CTime     string `json:"c_time"`
}

type Account struct {
	Id          int    `json:"id"`
	Pid         int    `json:"pid"`
	Username    string `json:"username"`
	NickName   string `json:"nick_name"`
	Pwd         string `json:"pwd"`
	Email       string `json:"email"`
	PrivilegeId int    `json:"privilege_id"`
	RoleId      int    `json:"role_id"`
	State       string `json:"state"`
	LlTime      string `json:"ll_time"`
	ChangeTime  string `json:"change_time"`
	CTime       string `json:"c_time"`
	Phone       string `json:"phone"`
	Remark      string `json:"remark"`
	Address     string `json:"address"`
}

type AccountUpdate struct {
	Id       string `json:"id"`
	NickName string `json:"nick_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Address  string `json:"address"`
	Remark   string `json:"remark"`
}