/**
* @Author: yanKoo
* @Date: 2019/3/11 10:48
* @Description: 处理请求的业务逻辑
 */
package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"model"
	"net/http"
	tc "pkg/customer" // table customer
	td "pkg/device"
	tg "pkg/group"         // table group
	tgc "pkg/group_member" // table group_device
	tu "pkg/user"
	"server/web/backend/src/utils"

	"service"
	"strconv"
)

// 创建下级账户
func CreateAccountBySuperior(c *gin.Context) {
	// 1. 取出Post中的表单内容
	uBody := &model.CreateAccount{}
	if err := c.BindJSON(uBody); err != nil {
		log.Println("bind json error")
		c.JSON(http.StatusUnprocessableEntity, model.ErrorRequestBodyParseFailed)
		return
	}

	// 2. 数据格式合法性校验，首先不能为空，其次每个格式都必须校验
	if uBody.NickName == "" || uBody.Username == "" || uBody.Pwd == "" || uBody.ConfirmPwd == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": " Please fill in the required fields ",
		})
		return
	}

	if !utils.CheckName(uBody.Username) || !utils.CheckName(uBody.NickName) || !utils.CheckPwd(uBody.ConfirmPwd) || !utils.CheckPwd(uBody.Pwd) {
		log.Println("username or nickname or pwd format is error")
		c.JSON(http.StatusUnprocessableEntity, model.ErrorRequestBodyParseFailed)
		return
	}

	// 只有创建 1是普通用户， 2是调度员， 3是经销商 4是超级管理员root
	if uBody.RoleId < 4 && uBody.RoleId >= 1 {
	} else {
		c.JSON(http.StatusUnprocessableEntity, model.ErrorRequestBodyParseFailed)
		return
	}

	// 名字查重
	aCount, err := tc.GetAccountByName(uBody.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	}
	if aCount > 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "User is exist already",
		})
		return
	}

	// 判断用户类型是否符合上级创建下级
	parentAccount, err := tc.GetAccount(uBody.Pid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	}

	// 只能给下级创建
	if parentAccount.RoleId <= uBody.RoleId {
		c.JSON(http.StatusInternalServerError, model.ErrorCreateAccountError)
	}

	// 3. 添加账户
	uId, err := tc.AddAccount(uBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	}

	//4. 生成session，返回消息内容
	//id := service.GenerateNewSessionId(uBody.Username)
	c.JSON(http.StatusCreated, gin.H{
		"success":    "true",
		"account_id": uId,
	})
}

// 获取账户信息
func GetAccountInfo(c *gin.Context) {
	aName := c.Param("account_name")
	if aName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "request url is not correct.",
			"error_code": "001",
		})
		return
	}

	// 使用session来校验用户
	if !service.ValidateAccountSession(c.Request, aName) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":      "session is not right.",
			"error_code": "006",
		})
		return
	}

	// 获取账户信息
	ai, err := tc.GetAccount(aName)
	if err != nil {
		log.Printf("Error in GetAccountInfo: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "get accountInfo DB error",
			"error_code": "007",
		})
		return
	}
	// 获取所有用户设备
	deviceAll, err := tu.SelectUserByAccountId(ai.Id)
	if err != nil {
		log.Printf("Error in GetGroups: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "get devices DB error",
			"error_code": "009",
		})
		return
	}

	// 获取群组信息
	groups, err := tg.SelectGroupsByAccountId(ai.Id)
	var gList []*model.GroupList
	for i := 0; i < len(groups); i++ {
		us, err := tgc.SelectDevicesByGroupId((*groups[i]).Id)
		if err != nil {
			log.Printf("Error in Get Group devices: %s", err)
		}
		log.Println(us)
		groupMember := make([]interface{}, 0)
		groupMember = append(groupMember, ai)
		for _, u := range us {
			groupMember = append(groupMember, &model.User{
				Id:         u.Id,
				IMei:       u.IMei,
				UserName:   u.UserName,
				PassWord:   u.PassWord,
				UserType:   u.UserType,
				AccountId:  u.AccountId,
				CreateTime: u.CreateTime,
				LLTime:     u.LLTime,
				ChangeTime: u.ChangeTime,
			})
		}
		// 群里的用户id
		ids, err := tgc.SelectDeviceIdsByGroupId((*groups[i]).Id)
		if err != nil {
			log.Printf("Error in GetGroups: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":      "get GetGroups DB error",
				"error_code": "008",
			})
			return
		}
		gListEle := &model.GroupList{DeviceInfo: groupMember, DeviceIds: ids, GroupInfo: groups[i]}
		gList = append(gList, gListEle)
	}

	resElem, err := tc.SelectChildByPId(ai.Id)
	if err != nil {
		log.Printf("db error : %s", err)
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	}

	cList := make([]*model.AccountClass, 0)
	for i := 0; i < len(resElem); i++ {
		child, err := tc.GetAccount((*resElem[i]).Id)
		if err != nil {
			log.Printf("db error : %s", err)
			c.JSON(http.StatusInternalServerError, model.ErrorDBError)
			return
		}

		cList = append(cList, &model.AccountClass{
			Id:              child.Id,
			AccountName:     child.Username,
			AccountNickName: child.NickName,
		})
	}

	resp := &model.AccountClass{
		Id:              ai.Id,
		AccountName:     ai.Username,
		AccountNickName: ai.NickName,
		Children:        cList,
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "User information obtained successfully",
		"account_info": ai,
		"device_list":  deviceAll,
		"group_list":   gList,
		"tree_data":    resp,
	})
}

// 更新账户信息 如果父id == 0 就是修改自己。
func UpdateAccountInfo(c *gin.Context) {
	accInf := &model.AccountUpdate{}
	if err := c.BindJSON(accInf); err != nil {
		log.Printf("%s", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "Request body is not correct.",
			"error_code": "001",
		})
		return
	}

	// 校验参数信息 ：校首先必须要有id，其次是每个参数的合法性，首先都不允许为空
	if accInf.LoginId == "" {
		log.Printf("account id is nil")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "The account id cannot be empty",
			"error_code": "003",
		})
		return
	}

	loginId, _ := strconv.Atoi(accInf.LoginId)

	// 使用session来校验用户
	if !service.ValidateAccountSession(c.Request, loginId) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":      "session is not right.",
			"error_code": "006",
		})
		return
	}

	/*if accInf.Phone == "" && accInf.Email == "" &&
		accInf.NickName == "" && accInf.Address == "" &&
		accInf.Remark == "" && accInf.Contact == "" {
		log.Printf("All changed parameters are null")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "Change at least one parameter",
			"error_code": "004",
		})
		return
	}*/

	if err := tc.UpdateAccount(accInf); err != nil {
		log.Println("Update account error :", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "get devices DB error",
			"error_code": "009",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": "true",
			"msg":     "update account success",
		})
	}
}

// 更新用户密码
func UpdateAccountPwd(c *gin.Context) {
	accPwd := &model.AccountPwd{}
	if err := c.BindJSON(accPwd); err != nil {
		log.Printf("json parse fail , error : %s", err)
		c.JSON(http.StatusBadRequest, model.ErrorRequestBodyParseFailed)
		return
	}

	// 使用session来校验用户
	aid, _ := strconv.Atoi(accPwd.Id)
	if !service.ValidateAccountSession(c.Request, aid) {
		c.JSON(http.StatusUnauthorized, model.ErrorNotAuthSession)
		return
	}
	// 校验参数信息 ：校首先必须要有id，其次是每个参数的合法性，首先都不允许为空，TODO 其次每一个参数的合法性，比如只能有汉字
	if accPwd.Id == "" {
		log.Printf("account id is nil")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "account id is null",
			"error_code": "001",
		})
		return
	}

	// 判断密码是否正确
	pwd, err := tc.GetAccountPwdByKey(aid)
	if err != nil {
		log.Printf("db error : %s", err)
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	}
	if pwd != accPwd.OldPwd {
		log.Printf("db pwd: %s", pwd)
		log.Printf("input old pwd %s", accPwd.OldPwd)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "Old password is not match",
			"error_code": "002",
		})
		return
	}

	// 新密码不能和旧密码不同
	if accPwd.NewPwd == accPwd.OldPwd {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "The new password cant't be the same as the old password.",
			"error_code": "003",
		})
		return
	}

	// 两次输入的密码必须一致
	if accPwd.NewPwd != accPwd.ConfirmPwd {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "The two passwords don't match",
			"error_code": "004",
		})
		return
	}

	// 更新密码
	id, _ := strconv.Atoi(accPwd.Id)
	if err := tc.UpdateAccountPwd(accPwd.NewPwd, id); err != nil {
		log.Println("Update account errr :", err)
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"result": "success",
			"msg":    "Password changed successfully",
		})
	}
}

// 获取账户下级目录
func GetAccountClass(c *gin.Context) {
	accountId := c.Param("accountId")

	// 使用session来校验用户
	aid, _ := strconv.Atoi(accountId)
	if !service.ValidateAccountSession(c.Request, aid) {
		c.JSON(http.StatusUnauthorized, model.ErrorNotAuthSession)
		return
	}

	// 查询数据返回
	root, err := tc.GetAccount(aid)
	if err != nil {
		log.Printf("db error : %s", err)
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	}

	resElem, err := tc.SelectChildByPId(aid)
	if err != nil {
		log.Printf("db error : %s", err)
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	}

	cList := make([]*model.AccountClass, 0)
	for i := 0; i < len(resElem); i++ {
		child, err := tc.GetAccount((*resElem[i]).Id)
		if err != nil {
			log.Printf("db error : %s", err)
			c.JSON(http.StatusInternalServerError, model.ErrorDBError)
			return
		}

		cList = append(cList, &model.AccountClass{
			Id:              child.Id,
			AccountName:     child.Username,
			AccountNickName: child.NickName,
		})
	}

	resp := &model.AccountClass{
		Id:              aid,
		AccountName:     root.Username,
		AccountNickName: root.NickName,
		Children:        cList,
	}

	c.JSON(http.StatusOK, gin.H{
		"result":    "success",
		"tree_data": resp,
	})
}

// 获取账户的设备信息
func GetAccountDevice(c *gin.Context) {
	accountId := c.Param("accountId")

	// 使用session来校验用户 TODO 考虑加一个, 保证只有上级用户和自己能访问这个接口
	aid, _ := strconv.Atoi(accountId)
	//if !service.ValidateAccountSession(c.Request, aid) {
	//	c.JSON(http.StatusUnauthorized, model.ErrorNotAuthSession)
	//	return
	//}
	// 获取账户信息
	ai, err := tc.GetAccount(aid)
	if err != nil {
		log.Printf("Error in GetAccountInfo: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "get accountInfo DB error",
			"error_code": "007",
		})
		return
	}
	// 获取所有设备
	deviceAll, err := tu.SelectUserByAccountId(ai.Id)
	if err != nil {
		log.Printf("db error : %s", err)
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"account_info": ai,
		"devices":      deviceAll,
	})
}

// 转移设备
func TransAccountDevice(c *gin.Context) {
	//aidStr := c.Param("accountId")
	accountDevices := &model.AccountDeviceTransReq{}
	if err := c.BindJSON(accountDevices); err != nil {
		log.Printf("json parse fail , error : %s", err)
		c.JSON(http.StatusBadRequest, model.ErrorRequestBodyParseFailed)
		return
	}

	// 使用session来校验用户
	//aid, _ := strconv.Atoi(aidStr)
	//if !service.ValidateAccountSession(c.Request, aid) {
	//	c.JSON(http.StatusUnauthorized, model.ErrorNotAuthSession)
	//	return
	//}

	for _, v := range accountDevices.Devices {
		if v.IMei == "" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":      "Imei can't be empty.",
				"error_code": "001",
			})
			return
		}
	}
	if accountDevices.Receiver.AccountId == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "receiver id can't be empty.",
			"error_code": "001",
		})
		return
	}

	// 更新设备
	if err := td.MultiUpdateDevice(accountDevices); err != nil {
		log.Printf("db error : %s", err)
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": "success",
		"msg":    "trans successful",
	})
}
