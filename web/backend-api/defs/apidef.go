/**
* @Author: yanKoo
* @Date: 2019/3/11 11:13
* @Description:
 */
package defs

// request
type AccountCredential struct {
	Id          int    `json:"id"`
	Pid         int    `json:"pid"`
	Username    string `json:"username"`
	Nick_name   string `json:"nick_name"`
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

// validate
type AccountValidate struct {
	Username string `valid:", between=6|15"`
	Pwd      string `valid:", between=6|15"`
}

//response
type AccountInfo struct {
	Id          int    `json:"id"`
	Username    string `json:"username"`
	PrivilegeId int    `json:"privilege_id"`
	IdCard      string `json:"id_card"`
}

type Message struct {
	Result  bool   `json:"result"`
	Type    string `json:"type"`
	Message string `json:"message"`
	//Sticky string `json:"sticky"`
}

// 返回账户及其账户下所有的组和
type AccountGroupsResp struct {
	Message     string             `json:"message"`
	AccountInfo *AccountCredential `json:"account_info"`
	GroupList   []*GroupList       `json:"group_list"`
	DeviceList  []*Device          `json:"device_list"`
}

type GroupList struct {
	DeviceIds []int      `json:"device_ids"`
	GroupInfo *GroupInfo `json:"group_info"`
}

// Data model
type SessionInfo struct {
	SessionID string `json:"session_id"`
	UserName  string `json:"user_name"`
	AccountId int    `json:"account_id"`
	TTL       string `json:"ttl"`
}
