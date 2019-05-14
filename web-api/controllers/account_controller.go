/**
* @Author: yanKoo
* @Date: 2019/3/11 10:48
* @Description: 处理请求的业务逻辑
 */
package controllers

import (
	"github.com/gin-gonic/gin"
	"server/web-api/log"
	"net/http"
	"server/web-api/cache"
	tc "server/web-api/dao/customer" // table customer
	td "server/web-api/dao/device"
	tg "server/web-api/dao/group"         // table group
	tgc "server/web-api/dao/group_member" // table group_device
	tlc "server/web-api/dao/location"
	tu "server/web-api/dao/user"
	tuc "server/web-api/dao/user_cache"
	"server/web-api/utils"
	"server/web-api/model"

	"server/web-api/service"
	"strconv"
)

// 创建下级账户
func CreateAccountBySuperior(c *gin.Context) {
	// 1. 取出Post中的表单内容
	uBody := &model.CreateAccount{}
	if err := c.BindJSON(uBody); err != nil {
		log.Log.Println("bind json error : ", err)
		c.JSON(http.StatusUnprocessableEntity, model.ErrorRequestBodyParseFailed)
		return
	}

	// 2. 数据格式合法性校验，首先不能为空，其次每个格式都必须校验
	if uBody.NickName == "" || uBody.Username == "" || uBody.Pwd == "" || uBody.ConfirmPwd == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      " Please fill in the required fields ",
			"error_code": "0001",
		})
		return
	}

	// 校验昵称
	if !utils.CheckNickName(uBody.NickName) {
		log.Log.Println("NickName format error", uBody.NickName)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "昵称只能输入1-20个以字母或者数字开头、可以含中文、下划线的字串。",
			"error_code": "0002",
		})
		return
	}

	if !utils.CheckUserName(uBody.Username) {
		log.Log.Println("Username format error")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "用户名只能输入5-20个包含字母、数字或下划线的字串",
			"error_code": "0003",
		})
		return
	}

	// 名字查重
	aCount, err := tc.GetAccountByName(uBody.Username)
	if err != nil {
		log.Log.Println("db error : ", err)
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	}
	if aCount > 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "User is exist already",
			"error_code": "0005",
		})
		return
	}

	// 校验密码
	if !utils.CheckPwd(uBody.ConfirmPwd) || !utils.CheckPwd(uBody.Pwd) {
		log.Log.Println("Pwd format error")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "密码6位-16位，至少包含一个数字字母",
			"error_code": "0004",
		})
		return
	}
	if uBody.ConfirmPwd != uBody.Pwd {
		log.Log.Println("Confirm Pwd is not match pwd")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "两次输入密码必须一致",
			"error_code": "0005",
		})
		return
	}

	// 只有创建 1是普通用户， 2是调度员， 3是经销商 4是公司，5是超级管理员root
	log.Log.Println("创建等级:", uBody.RoleId)
	if uBody.RoleId < 5 && uBody.RoleId >= 1 {
	} else {
		log.Log.Println("创建权限出错")
		c.JSON(http.StatusUnprocessableEntity, model.ErrorRequestBodyParseFailed)
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
		return
	}

	// 3. 添加账户
	uId, err := tc.AddAccount(uBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	}

	//4. 返回消息内容
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
		log.Log.Printf("Error in GetAccountInfo: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "get accountInfo DB error",
			"error_code": "007",
		})
		return
	}
	// 获取所有用户设备
	deviceAll, err := tu.SelectUserByAccountId(ai.Id)
	if err != nil {
		log.Log.Printf("Error in GetGroups: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "get devices DB error",
			"error_code": "009",
		})
		return
	}

	for _, v := range deviceAll {
		res, _, err := tlc.GetUserLocationInCache(int32(v.Id), cache.GetRedisClient())
		if err != nil {
			log.Log.Printf("GetGpsData error: %+v", err)
		}
		if res != nil {
			v.GPSData = &model.GPS{
				Lng: res.GpsInfo.Longitude,
				Lat: res.GpsInfo.Latitude,
			}
			v.Course = res.GpsInfo.Course
			v.Speed = res.GpsInfo.Speed
			v.LocalTime = res.GpsInfo.LocalTime
		}
	}

	// 获取群组信息
	groups, err := tg.SelectGroupsByAccountId(ai.Id)
	var gList []*model.GroupList
	for i := 0; i < len(groups); i++ {
		us, err := tgc.SelectDevicesByGroupId((*groups[i]).Id)
		if err != nil {
			log.Log.Printf("Error in Get Group devices: %s", err)
		}
		groupMember := make([]interface{}, 0)
		groupOfflineMember := make([]interface{}, 0)
		groupMember = append(groupMember, ai)

		for _, u := range us {
			// 获取该设备在线状态
			online, err := tuc.GetUserStatusFromCache(int32(u.Id), cache.GetRedisClient())
			if err != nil || online != tuc.USER_ONLINE {
				continue
			}
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
				Online:     int(online),
			})
		}
		for _, u := range us {
			// 获取离线的设备
			online, err := tuc.GetUserStatusFromCache(int32(u.Id), cache.GetRedisClient())
			if err != nil {
				online = tuc.USER_OFFLINE
			}
			if online == tuc.USER_OFFLINE {
				groupOfflineMember = append(groupOfflineMember, &model.User{
					Id:         u.Id,
					IMei:       u.IMei,
					UserName:   u.UserName,
					PassWord:   u.PassWord,
					UserType:   u.UserType,
					AccountId:  u.AccountId,
					CreateTime: u.CreateTime,
					LLTime:     u.LLTime,
					ChangeTime: u.ChangeTime,
					Online:     int(online),
				})
			}
		}
		groupMember = append(groupMember, groupOfflineMember...)

		// 群里的用户id
		ids, err := tgc.SelectDeviceIdsByGroupId((*groups[i]).Id)
		if err != nil {
			log.Log.Printf("Error in GetGroups: %s", err)
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
		log.Log.Printf("db error : %s", err)
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	}

	cList := make([]*model.AccountClass, 0)
	for i := 0; i < len(resElem); i++ {
		child, err := tc.GetAccount((*resElem[i]).Id)
		if err != nil {
			log.Log.Printf("db error : %s", err)
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

	ai.Pwd = "" // 不把密码暴露出去
	c.JSON(http.StatusOK, gin.H{
		"message":      "User information obtained successfully",
		"account_info": ai,
		"device_list":  deviceAll,
		"group_list":   gList,
		"tree_data":    resp,
	})
}

// 更新下级账户信息
func UpdateAccountInfo(c *gin.Context) {
	accInf := &model.AccountUpdate{}
	if err := c.BindJSON(accInf); err != nil {
		log.Log.Printf("%s", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "Request body is not correct.",
			"error_code": "001",
		})
		return
	}

	// 校验参数信息 ：校首先必须要有id，其次是每个参数的合法性，首先都不允许为空
	if accInf.LoginId == "" {
		log.Log.Printf("account id is nil")
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

	if err := tc.UpdateAccount(accInf); err != nil {
		log.Log.Println("Update account error :", err)
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
		log.Log.Printf("json parse fail , error : %s", err)
		c.JSON(http.StatusBadRequest, model.ErrorRequestBodyParseFailed)
		return
	}

	// 使用session来校验用户是否登录
	aid, _ := strconv.Atoi(accPwd.Id)
	if !service.ValidateAccountSession(c.Request, aid) {
		c.JSON(http.StatusUnauthorized, model.ErrorNotAuthSession)
		return
	}
	// 校验参数信息 ：校首先必须要有id，都不允许为空
	if accPwd.Id == "" {
		log.Log.Printf("account id is nil")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "account id is null",
			"error_code": "001",
		})
		return
	}

	// 校验密码
	if !utils.CheckPwd(accPwd.ConfirmPwd) || !utils.CheckPwd(accPwd.NewPwd) {
		log.Log.Println("Pwd format error")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "密码6位-16位，至少包含一个数字字母",
			"error_code": "0004",
		})
		return
	}
	// 两次输入的密码必须一致
	if accPwd.ConfirmPwd != accPwd.NewPwd {
		log.Log.Println("Confirm Pwd is not match pwd")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "两次输入密码必须一致",
			"error_code": "0005",
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

	// 判断密码是否正确
	pwd, err := tc.GetAccountPwdByKey(aid)
	if err != nil {
		log.Log.Printf("db error : %s", err)
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	}
	if pwd != accPwd.OldPwd {
		log.Log.Printf("db pwd: %s", pwd)
		log.Log.Printf("input old pwd %s", accPwd.OldPwd)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "Old password is not match",
			"error_code": "002",
		})
		return
	}

	// 更新密码
	id, _ := strconv.Atoi(accPwd.Id)
	if err := tc.UpdateAccountPwd(accPwd.NewPwd, id); err != nil {
		log.Log.Println("Update account errr :", err)
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
	searchId := c.Param("searchId")

	log.Log.Println("searchId:", searchId, "accountId", accountId)
	//使用session来校验用户
	aid, _ := strconv.Atoi(accountId)
	sid, _ := strconv.Atoi(searchId)

	if !service.ValidateAccountSession(c.Request, aid) {
		c.JSON(http.StatusUnauthorized, model.ErrorNotAuthSession)
		return
	}

	// 查询数据返回
	root, err := tc.GetAccount(sid)
	if err != nil {
		log.Log.Printf("GetAccount db error : %s", err)
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	}

	resElem, err := tc.SelectChildByPId(sid)
	if err != nil {
		log.Log.Printf("SelectChildByPId db error : %s", err)
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	}

	cList := make([]*model.AccountClass, 0)
	for i := 0; i < len(resElem); i++ {
		child, err := tc.GetAccount((*resElem[i]).Id)
		if err != nil {
			log.Log.Printf("db error : %s", err)
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
		Id:              sid,
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
	getAdviceId := c.Param("getAdviceId")

	// 使用session来校验用户, 保证上级用户已登录
	aid, _ := strconv.Atoi(accountId)
	getAId, _ := strconv.Atoi(getAdviceId)
	if !service.ValidateAccountSession(c.Request, aid) {
		c.JSON(http.StatusUnauthorized, model.ErrorNotAuthSession)
		return
	}
	// 获取账户信息
	ai, err := tc.GetAccount(getAId)
	if err != nil {
		log.Log.Printf("Error in GetAccountInfo: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "get accountInfo DB error",
			"error_code": "007",
		})
		return
	}
	// 获取所有设备
	deviceAll, err := tu.SelectUserByAccountId(ai.Id)
	if err != nil {
		log.Log.Printf("db error : %s", err)
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
	aidStr := c.Param("accountId")
	accountDevices := &model.AccountDeviceTransReq{}
	if err := c.BindJSON(accountDevices); err != nil {
		log.Log.Printf("json parse fail , error : %s", err)
		c.JSON(http.StatusBadRequest, model.ErrorRequestBodyParseFailed)
		return
	}

	// 使用session来校验用户
	aid, _ := strconv.Atoi(aidStr)
	if !service.ValidateAccountSession(c.Request, aid) {
		c.JSON(http.StatusUnauthorized, model.ErrorNotAuthSession)
		return
	}

	// IMEI号只能是15位数字
	// 结构体为空
	if accountDevices.Devices == nil {
		c.JSON(http.StatusBadRequest, model.ErrorRequestBodyParseFailed)
		return
	}
	for _, v := range accountDevices.Devices {
		if v.IMei == "" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":      "Imei is not correct.",
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
		log.Log.Printf("db error : %s", err)
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": "success",
		"msg":    "trans successful",
	})
}
