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

type AccountDeviceTransReq struct {
	Devices  []*Device      `json:"devices"`
	Sender   Account        `json:"sender"`
	Receiver DeviceReceiver `json:"receiver"`
}

type AccountImportDeviceReq struct {
	DeviceIMei []string `json:"device_imei"`
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
	Id              int             `json:"id"`
	AccountName     string          `json:"account_name"`
	AccountNickName string          `json:"account_nickname"`
	Children        []*AccountClass `json:"children"`
}

// 返回账户及其账户下所有的组和
type AccountGroupsResp struct {
	Message     string       `json:"message"`
	AccountInfo *Account     `json:"account_info"`
	GroupList   []*GroupList `json:"group_list"`
	DeviceList  []*Device    `json:"device_list"`
}

type GroupList struct {
	DeviceIds []int `json:"device_ids"`
	//DeviceInfo []*Device  `json:"device_infos"`
	DeviceInfo []interface{} `json:"device_infos"`
	GroupInfo  *GroupInfo    `json:"group_info"`
}

// Data model
type SessionInfo struct {
	SessionID string `json:"session_id"`
	UserName  string `json:"user_name"`
	UserPwd   string `json:"user_pwd"`
	AccountId int    `json:"account_id"`
	TTL       string `json:"ttl"`
}

type DeviceReceiver struct {
	AccountId   int    `json:"account_id"`
	AccountName string `json:"account_name"`
}

// device
type Device struct {
	Id           int            `json:"id"`
	IMei         string         `json:"imei"`
	UserName     string         `json:"user_name"`
	PassWord     string         `json:"password"`
	AccountId    int            `json:"account_id"`
	CreateTime   string `json:"create_time"`
	LLTime       string `json:"last_login_time"`
	ChangeTime   string `json:"change_time"`
}

// User
type User struct {
	Id         int    `json:"id"`
	IMei       string `json:"imei"`
	UserName   string `json:"user_name"`
	NickName   string `json:"nick_name"`
	PassWord   string `json:"password"`
	UserType   int    `json:"user_type"`  // 用户类型(暂定1是普通用户，2是调度员，3是经销商, 4是超级管理员)
	ParentId   string `json:"parent_id"`  // 如果是普通用户注册的时候，默认是0， 如果是上级用户创建下级账户，就用来表示创建者的id
	AccountId  int    `json:"account_id"` // 只有普通用户才有这个字段，表示这个设备属于哪个账户，如果是非普通用户就是默认为0（因为customer表里面没有0号）
	CreateTime string `json:"create_time"`
	LLTime     string `json:"last_login_time"`
	ChangeTime string `json:"change_time"`
}

type Customer struct {
	Id         int            `json:"id"`
	UId        int            `json:"user_id"`
	PId        int            `json:"parent_id"`
	Email      sql.NullString `json:"email"`
	Phone      sql.NullString `json:"phone"`
	Address    sql.NullString `json:"address"`
	Remark     sql.NullString `json:"remark"`
	Contact    sql.NullString `json:"contact"`
	ChangeTime string         `json:"change_time"`
	CTime      string         `json:"create_time"`
}

type GroupInfo struct {
	Id        int    `json:"id"`
	GroupName string `json:"group_name"`
	AccountId int    `json:"account_id"`
	Status    string `json:"status"`
	CTime     string `json:"c_time"`
}

type Account struct {
	Id          int            `json:"id"`
	Pid         int            `json:"pid"`
	Username    string         `json:"username"`
	NickName    string         `json:"nick_name"`
	Pwd         string         `json:"pwd"`
	Email       sql.NullString `json:"email"`
	PrivilegeId int            `json:"privilege_id"`
	Contact     sql.NullString `json:"contact"`
	RoleId      int            `json:"role_id"`
	State       string         `json:"state"`
	LlTime      string         `json:"ll_time"`
	ChangeTime  string         `json:"change_time"`
	CTime       string         `json:"c_time"`
	Phone       sql.NullString `json:"phone"`
	Remark      sql.NullString `json:"remark"`
	Address     sql.NullString `json:"address"`
}

type CreateAccount struct {
	ConfirmPwd  string `json:"confirm_pwd"`
	Id          int    `json:"id"`
	Pid         int    `json:"pid"`
	Username    string `json:"username"`
	NickName    string `json:"nick_name"`
	Pwd         string `json:"pwd"`
	Email       string `json:"email"`
	PrivilegeId int    `json:"privilege_id"`
	Contact     string `json:"contact"`
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
	LoginId  string `json:"login_id"`
	NickName string `json:"nick_name"`
	Username string `json:"username"`
	TypeId   string `json:"type_id"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Address  string `json:"address"`
	Remark   string `json:"remark"`
	Contact  string `json:"contact"`
}
