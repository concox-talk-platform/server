/**
* @Author: yanKoo
* @Date: 2019/3/16 15:32
* @Description:
 */
package defs

// device
type Device struct {
	Id           int    `json:"id"`
	IMei         string `json:"imei"`
	AliasName    string `json:"alias_name"`
	AccountId    int    `json:"account_id"`
	Status       string `json:"status"`
	ActiveStatus string `json:"active_status"`
	BindStatus   string `json:"bind_status"`
	CrateTime    string `json:"crate_time"`
}

type GroupInfo struct {
	Id        int    `json:"id"`
	GroupName string `json:"group_name"`
	AccountId int    `json:"account_id"`
	Status    string `json:"status"`
	CTime     string `json:"c_time"`
}

// updateAccount modal
type Account struct {
	Id       string `json:"id"`
	NickName string `json:"nick_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Address  string `json:"address"`
	Remark   string `json:"remark"`
}

// update pwd modal
type AccountPwd struct {
	Id         string `json:"id"`
	OldPwd     string `json:"old_pwd"`
	NewPwd     string `json:"new_pwd"`
	ConfirmPwd string `json:"confirm_pwd"`
}
