/**
* @Author: yanKoo
* @Date: 2019/3/11 10:48
* @Description: 处理请求的业务逻辑
 */
package controllers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/syssam/go-validator"
	"log"
	"model"
	"net/http"
	tc "pkg/customer"
	"service"
)

// 创建账户
func SignUp(c *gin.Context) {
	// 1. 取出Post中的表单内容
	uBody := &model.Account{}
	if err := c.BindJSON(uBody); err != nil {
		log.Printf("解析错误")
	}
	// 2. TODO 数据格式合法性校验
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

	// 3. 添加账户
	if err := tc.AddAccount(&model.Account{Username: uBody.Username, Pwd: uBody.Pwd}); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	}

	// 4. 生成session，返回消息内容
	id := service.GenerateNewSessionId(uBody.Username)
	c.JSON(http.StatusCreated, gin.H{
		"success":    "true",
		"session_id": id,
	})
}

// 登录
func SignIn(c *gin.Context) {
	// 1. 取出Post请求中的body内容
	signINBody := &model.Account{}
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
	uInfo, err := tc.GetAccount(signINBody.Username)
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
	var vBody = model.AccountValidate{Pwd: signINBody.Pwd, Username: signINBody.Username}
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
	var sInfo = &model.SessionInfo{
		SessionID: service.GetSessionId(c.Request),
		UserName:  signINBody.Username,
		AccountId: uInfo.Id,
	}
	id, err := service.UpdateUserSessionId(sInfo)
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
	signOutBody := &model.Account{}
	if err := c.BindJSON(signOutBody); err != nil {
		log.Printf("%s", err)
		c.JSON(http.StatusBadRequest, model.ErrorRequestBodyParseFailed)
		return
	}

	// 2. TODO 验证session
	if !service.ValidateAccountSession(c.Request, signOutBody.Username) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":      "Request body is not correct.",
			"error_code": "002",
		})
		return
	}
	// 3. 删除session
	if err := service.DeleteExpiredSession(service.GetSessionId(c.Request)); err != nil {
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
