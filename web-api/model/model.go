/**
* @Author: yanKoo
* @Date: 2019/3/16 15:32
* @Description:
 */
package model

import (
	"database/sql"
)

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
	Receiver DeviceReceiver `json:"receiver"`
	Sender   Account        `json:"sender"`
}

type AccountImportDeviceReq struct {
	Devices []*Device        `json:"devices"  binding:"required"`
	DType   string           `json:"d_type"`
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
	DeviceIds  []int         `json:"device_ids"`
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
	Id         int     `json:"id"`
	IMei       string  `json:"imei"`
	UserName   string  `json:"user_name"`
	NickName   string  `json:"nick_name"`
	PassWord   string  `json:"password"`
	AccountId  int     `json:"account_id"`
	CreateTime string  `json:"create_time"`
	LLTime     string  `json:"last_login_time"`
	ChangeTime string  `json:"change_time"`
	LocalTime  uint64  `json:"local_time"`
	GPSData    *GPS    `json:"gps_data"`
	Speed      float32 `json:"speed"`
	Course     float32 `json:"course"`
	DeviceType string  `json:"device_type"`
	ActiveTime string  `json:"active_time"`
	SaleTime   string  `json:"sale_time"`
}

type GPS struct {
	Lng float64 `json:"lng"`
	Lat float64 `json:"lat"`
}

// User
type User struct {
	Id          int    `json:"id"`
	IMei        string `json:"imei"`
	UserName    string `json:"user_name"`
	NickName    string `json:"nick_name"`
	PassWord    string `json:"password"`
	UserType    int    `json:"user_type"`  // 用户类型(暂定1是普通用户，2是调度员，3是经销商, 4是超级管理员)
	ParentId    string `json:"parent_id"`  // 如果是普通APP用户和设备注册的时候，默认是0， 如果是上级用户创建下级账户，就用来表示创建者的id
	AccountId   int    `json:"account_id"` // 只有普通用户才有这个字段，表示这个设备属于哪个账户，如果是非普通用户就是默认为0（因为customer表里面没有0号）
	LockGroupId int    `json:"lock_group_id"`
	CreateTime  string `json:"create_time"`
	LLTime      string `json:"last_login_time"`
	ChangeTime  string `json:"change_time"`
	Online      int    `json:"online"`
	DeviceType  string `json:"device_type"`
	ActiveTime  string `json:"active_time"`
	SaleTime    string `json:"sale_time"`
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
	OnlineNum int    `json:"online_num"`
}


type AccountForSwag struct {
	Username    string         `json:"username" binding:"required" example:"elephant"`
	Pwd         string         `json:"pwd" example:"123456"`
}
type Account struct {
	Id          int            `json:"id"`
	Pid         int            `json:"pid"`
	Username    string         `json:"username" binding:"required" example:"elephant"`
	NickName    string         `json:"nick_name"`
	Pwd         string         `json:"pwd" example:"123456"`
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

type CreateAccountForSwag struct {
	ConfirmPwd  string `json:"confirm_pwd" example:"111111"`
	Pid         int    `json:"pid" example:"13"`
	Username    string `json:"username" example:"dis007"`
	NickName    string `json:"nick_name" example:"nickDis001"`
	Pwd         string `json:"pwd" example:"123456"`
	RoleId      int    `json:"role_id" example:"3"`
}

type CreateAccount struct {
	ConfirmPwd  string `json:"confirm_pwd" example:"111111"`
	Id          int    `json:"id" example:"5"`
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

type ImMsgData struct {
	Id           int    `json:"id"`
	SenderName   string `json:"SenderName"`
	ReceiverType int    `json:"ReceiverType"`
	ReceiverId   int    `json:"ReceiverId"`
	ReceiverName string `json:"ReceiverName"`
	ResourcePath string `json:"ResourcePath"`
	MsgType      int    `json:"MsgType"`
	SendTime     string `json:"SendTime"`
}

type FileContext struct {
	UserId         int
	FilePath       string
	FileType       int32
	FileParams     *ImMsgData
	FileName       string
	FileSize       int
	FileMD5        string
	FileFastId     string
	FileUploadTime string
}

// just for swagger use
type DeviceUpdate struct {
	Id                   int32    `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	IMei                 string   `protobuf:"bytes,2,opt,name=IMei,proto3" json:"IMei,omitempty"`
	NickName             string   `protobuf:"bytes,3,opt,name=NickName,proto3" json:"NickName,omitempty"`
	LoginId              int32    `protobuf:"varint,4,opt,name=LoginId,proto3" json:"LoginId,omitempty"`
	CreateTime           string   `protobuf:"bytes,6,opt,name=CreateTime,proto3" json:"CreateTime,omitempty"`
	LLTime               string   `protobuf:"bytes,7,opt,name=LLTime,proto3" json:"LLTime,omitempty"`
	ChangeTime           string   `protobuf:"bytes,8,opt,name=ChangeTime,proto3" json:"ChangeTime,omitempty"`
	LocalTime            uint64   `protobuf:"varint,9,opt,name=LocalTime,proto3" json:"LocalTime,omitempty"`
	DeviceType           string   `protobuf:"bytes,12,opt,name=DeviceType,proto3" json:"DeviceType,omitempty"`
	ActiveTime           string   `protobuf:"bytes,13,opt,name=ActiveTime,proto3" json:"ActiveTime,omitempty"`
	SaleTime             string   `protobuf:"bytes,14,opt,name=SaleTime,proto3" json:"SaleTime,omitempty"`
}