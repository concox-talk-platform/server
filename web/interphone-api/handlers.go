/**
* @Author: yanKoo
* @Date: 2019/3/11 10:48
* @Description: 处理请求的业务逻辑
 */
package main

import "C"
import (
	"database/sql"
	"github.com/bittiger/server/web/interphone-api/dbops"
	"github.com/bittiger/server/web/interphone-api/defs"
	"github.com/bittiger/server/web/interphone-api/session"
	"github.com/gin-gonic/gin"
	"github.com/syssam/go-validator"
	"log"
	"net/http"
	"strconv"
)

// 创建账户
func SignUp(c *gin.Context) {
	// 1. 取出Post中的表单内容
	uBody := &defs.AccountCredential{}
	if err := c.BindJSON(uBody); err != nil {
		log.Fatalf("解析错误")
	}
	// 2. TODO 数据格式合法性校验
	aCount, err := dbops.GetAccountByName(uBody.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, defs.ErrorDBError)
		return
	}
	if aCount > 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "User is exist already",
		})
		return
	}

	// 3. 添加账户
	if err := dbops.AddAccountCredential(uBody.Username, uBody.Pwd); err != nil {
		c.JSON(http.StatusInternalServerError, defs.ErrorDBError)
		return
	}

	// 4. 生成session，返回消息内容
	id := session.GenerateNewSessionId(uBody.Username)
	c.JSON(http.StatusCreated, gin.H{
		"success":    "true",
		"session_id": id,
	})
}

// 登录
func SignIn(c *gin.Context) {
	// 1. 取出Post请求中的body内容
	signINBody := &defs.AccountCredential{}
	if err := c.BindJSON(signINBody); err != nil {
		log.Printf("%s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Request body is not correct.",
		})
		return
	}

	// 2. 验证body里面的名字和url的名字
	if c.Param("account_name") != signINBody.Username {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":      "User 和 url 不匹配",
			"error_code": "0021",
		})
		return
	}

	// 3. 数据库查询密码，看是否和发过来的相同
	uInfo, err := dbops.GetAccount(signINBody.Username)
	if err != nil && err != sql.ErrNoRows {
		log.Println("err:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "db ops err",
			"error_code": "003",
		})
		return
	}
	if err == sql.ErrNoRows {
		log.Println("no found this user")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":      "User does not exist.",
			"error_code": "0021",
		})
		return
	}

	log.Printf("Login pwd: %s, loginBody pwd is %s", uInfo.Pwd, signINBody.Pwd)
	if err != nil || uInfo.Pwd != signINBody.Pwd {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":      "User password is wrong.",
			"error_code": "0022",
		})
		return
	}

	// 4. 对登录表单进行验证
	var vBody = defs.AccountValidate{Pwd: signINBody.Pwd, Username: signINBody.Username}
	err = validator.ValidateStruct(vBody)
	if err != nil {
		for _, err := range err.(validator.Errors) {
			log.Println("form validate err: ", err)
		}
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "account or password format is not correct.",
			"error_code": "0023",
		})
		return
	}
	// 5. 更新session
	// 先取session，不用判断这里不用判断session是否过期，因为，你已经login请求了，说明新建立一个session，直接更新session
	var sInfo = &defs.SessionInfo{
		SessionID: GetSessionId(c.Request),
		UserName:  signINBody.Username,
		AccountId: uInfo.Id,
	}
	id, err := session.UpdateUserSessionId(sInfo)
	if err != nil {
		log.Println("login update session err: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Update user SessionId error",
		})
		return
	}

	// 6. 返回登录成功的消息
	c.JSON(http.StatusOK, gin.H{
		"success":    "true",
		"session_id": id,
	})
}

// 退出登录
func SignOut(c *gin.Context) {
	// 1. 取出body中的内容
	signOutBody := &defs.AccountCredential{}
	if err := c.BindJSON(signOutBody); err != nil {
		log.Printf("%s", err)
		c.JSON(http.StatusBadRequest, defs.ErrorRequestBodyParseFailed)
		return
	}

	// 2. TODO 验证session
	if !validateAccountSession(c.Request, signOutBody.Username) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":      "Request body is not correct.",
			"error_code": "002",
		})
		return
	}
	// 3. 删除session
	if err := session.DeleteExpiredSession(GetSessionId(c.Request)); err != nil {
		log.Printf("SignOut delete session error : %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "SignOut delete session DB error",
			"error_code": "003",
		})
		return
	}

	// 4. 返回消息
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg":     "delete session id is successful",
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
	if !validateAccountSession(c.Request, aName) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":      "session is not right.",
			"error_code": "006",
		})
		return
	}

	// 获取账户信息
	ai, err := dbops.GetAccount(aName)
	if err != nil {
		log.Printf("Error in GetAccountInfo: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "get accountInfo DB error",
			"error_code": "007",
		})
		return
	}

	// 获取群组信息
	groups, err := dbops.SelectGroupsByAccountId(ai.Id)
	var gList []*defs.GroupList
	for i := 0; i < len(groups); i++ {
		ids, err := dbops.SelectDeviceIdsByGroupId((*groups[i]).Id)
		if err != nil {
			log.Printf("Error in GetGroups: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":      "get GetGroups DB error",
				"error_code": "008",
			})
			return
		}
		gListEle := &defs.GroupList{DeviceIds: ids, GroupInfo: groups[i]}
		gList = append(gList, gListEle)
	}

	// 获取所有设备,并且在组里面
	devices, err := dbops.SelectDeviceByAccountId(ai.Id)
	if err != nil {
		log.Printf("Error in GetGroups: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "get devices DB error",
			"error_code": "009",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "User information obtained successfully",
		"account_info": ai,
		"group_list":   gList,
		"device_list":  devices,
	})
}

// 更新账户信息
func UpdateAccountInfo(c *gin.Context) {
	accInf := &defs.Account{}
	if err := c.BindJSON(accInf); err != nil {
		log.Printf("%s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "Request body is not correct.",
			"error_code": "001",
		})
		return
	}

	// 校验参数信息 ：校首先必须要有id，其次是每个参数的合法性，首先都不允许为空，TODO 其次每一个参数的合法性，比如只能有汉字
	if accInf.Id == "" {
		log.Printf("account id is nil")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "The account id cannot be empty",
			"error_code": "003",
		})
		return
	}

	// 使用session来校验用户
	if !validateAccountSession(c.Request, accInf.Id) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":      "session is not right.",
			"error_code": "006",
		})
		return
	}

	if accInf.Phone == "" && accInf.Email == "" && accInf.NickName == "" && accInf.Address == "" && accInf.Remark == "" {
		log.Printf("All changed parameters are null")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "Change at least one parameter",
			"error_code": "004",
		})
		return
	}

	if err := dbops.UpdateAccount(accInf); err != nil {
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
	accPwd := &defs.AccountPwd{}
	if err := c.BindJSON(accPwd); err != nil {
		log.Printf("json parse fail , error : %s", err)
		c.JSON(http.StatusBadRequest, defs.ErrorRequestBodyParseFailed)
		return
	}

	// 使用session来校验用户
	aid, _ := strconv.Atoi(accPwd.Id)
	if !validateAccountSession(c.Request, aid) {
		c.JSON(http.StatusUnauthorized, defs.ErrorNotAuthSession)
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
	pwd, err := dbops.GetAccountCredential(aid)
	if err != nil {
		log.Fatalf("db error : %s", err)
		c.JSON(http.StatusInternalServerError, defs.ErrorDBError)
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
	if err := dbops.UpdateAccountPwd(accPwd.NewPwd, id); err != nil {
		log.Println("Update account errr :", err)
		c.JSON(http.StatusInternalServerError, defs.ErrorDBError)
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"result": "success",
			"msg":    "Password changed successfully",
		})
	}
}

// 创建群组
func CreateGroup(c *gin.Context) {
	gInfo := &defs.GroupList{}

	if err := c.BindJSON(gInfo); err != nil {
		log.Printf("json parse fail , error : %s", err)
		c.JSON(http.StatusBadRequest, defs.ErrorRequestBodyParseFailed)
		return
	}

	// TODO 校验创建群组信息的参数合法性
	if gInfo.GroupInfo.GroupName == "" || gInfo.GroupInfo.AccountId == 0 || len(gInfo.DeviceIds) == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "You need at least the group name, the account id, and at least one device id",
			"error_code": "001",
		})
		return
	}

	// 创建群组
	if err := dbops.CreateGroup(gInfo); err != nil {
		log.Printf("create group fail , error: %s", err)
		c.JSON(http.StatusInternalServerError, defs.ErrorDBError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "success",
		"msg":    "Create group successfully",
	})
}

// TODO 群组的成员更新
func UpdateGroup(c *gin.Context) {

}

func DeleteGroup() {

}
