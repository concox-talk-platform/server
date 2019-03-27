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
	ta "pkg/account" // table account
	td "pkg/device"
	tg "pkg/group"         // table group
	tgd "pkg/group_device" // table group_device

	"service"
	"strconv"
)

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
	ai, err := ta.GetAccount(aName)
	if err != nil {
		log.Printf("Error in GetAccountInfo: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "get accountInfo DB error",
			"error_code": "007",
		})
		return
	}
	// 获取所有设备,并且在组里面
	deviceAll, err := td.SelectDeviceByAccountId(ai.Id)
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
		ds, err := tgd.SelectDevicesByGroupId((*groups[i]).Id)
		if err != nil {
			log.Printf("Error in Get Group devices: %s", err)
		}
		log.Println(ds)
		devices := make([]*model.Device, 0)
		for _, d := range ds {
			devices = append(devices, &model.Device{
				Id:           d.Id,
				IMei:         d.IMei,
				UserName:     d.UserName,
				PassWord:     d.PassWord,
				AccountId:    d.AccountId,
				Status:       d.Status,
				ActiveStatus: d.ActiveStatus,
				BindStatus:   d.BindStatus,
				CrateTime:    d.CrateTime,
				LLTime:       d.LLTime,
				ChangeTime:   d.ChangeTime,
			})
		}
		ids, err := tgd.SelectDeviceIdsByGroupId((*groups[i]).Id)
		if err != nil {
			log.Printf("Error in GetGroups: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":      "get GetGroups DB error",
				"error_code": "008",
			})
			return
		}
		gListEle := &model.GroupList{DeviceInfo: devices, DeviceIds: ids, GroupInfo: groups[i]}
		gList = append(gList, gListEle)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "User information obtained successfully",
		"account_info": ai,
		"device_list":  deviceAll,
		"group_list":   gList,
	})
}

// 更新账户信息
func UpdateAccountInfo(c *gin.Context) {
	accInf := &model.AccountUpdate{}
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
	if !service.ValidateAccountSession(c.Request, accInf.Id) {
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

	if err := ta.UpdateAccount(accInf); err != nil {
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
	pwd, err := ta.GetAccountPwdByKey(aid)
	if err != nil {
		log.Fatalf("db error : %s", err)
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
	if err := ta.UpdateAccountPwd(accPwd.NewPwd, id); err != nil {
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
	root, err := ta.GetAccount(aid)
	if err != nil {
		log.Fatalf("db error : %s", err)
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	}

	resElem, err := ta.SelectChildByPId(aid)
	if err != nil {
		log.Fatalf("db error : %s", err)
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	}

	cList := make([]*model.AccountClass, 0)
	for i := 0; i < len(resElem); i++ {
		child, err := ta.GetAccount((*resElem[i]).Id)
		if err != nil {
			log.Fatalf("db error : %s", err)
			c.JSON(http.StatusInternalServerError, model.ErrorDBError)
			return
		}

		cList = append(cList, &model.AccountClass{
			Id:          child.Id,
			AccountName: child.Username,
		})
	}

	resp := &model.AccountClass{
		Id:          aid,
		AccountName: root.Username,
		Children:    cList,
	}

	c.JSON(http.StatusOK, gin.H{
		"result":    "success",
		"tree_data": resp,
	})
}

func GetAccountDevice(c *gin.Context)  {
	accountId := c.Param("accountId")

	// 使用session来校验用户 TODO 考虑加一个
	aid, _ := strconv.Atoi(accountId)
	//if !service.ValidateAccountSession(c.Request, aid) {
	//	c.JSON(http.StatusUnauthorized, model.ErrorNotAuthSession)
	//	return
	//}
	// 获取账户信息
	ai, err := ta.GetAccount(aid)
	if err != nil {
		log.Printf("Error in GetAccountInfo: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "get accountInfo DB error",
			"error_code": "007",
		})
		return
	}
	// 获取所有设备
	deviceAll, err := td.SelectDeviceByAccountId(ai.Id)
	if err != nil {
		log.Fatalf("db error : %s", err)
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"account_info" : ai,
		"devices" : deviceAll,
	})
}

// 转移设备
func TransAccountDevice(c *gin.Context)  {
	aidStr := c.Param("accountId")
	accountDevices := &model.AccountDeviceTransReq{}
	if err := c.BindJSON(accountDevices); err != nil {
		log.Printf("json parse fail , error : %s", err)
		c.JSON(http.StatusBadRequest, model.ErrorRequestBodyParseFailed)
		return
	}

	// 使用session来校验用户
	aid, _ := strconv.Atoi(aidStr)
	if !service.ValidateAccountSession(c.Request, aid) {
		c.JSON(http.StatusUnauthorized, model.ErrorNotAuthSession)
		return
	}

	// 更新设备
	if err := td.MultiUpdateDevice(accountDevices); err != nil {
		log.Fatalf("db error : %s", err)
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result" : "success",
		"msg" : "trans successful",
	})
}