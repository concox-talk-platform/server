/**
* @Author: yanKoo
* @Date: 2019/3/11 10:48
* @Description: 处理请求的业务逻辑
 */
package controllers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	tc "server/common/dao/customer"
	"server/common/utils"
	"server/web-api/model"
	ss "server/web-api/service"
)

// 登录
func SignIn(c *gin.Context) {
	// 1. 取出Post请求中的body内容
	logrus.Info("logrous info")
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

	// 3. 对登录表单进行验证
	if !utils.CheckUserName(signINBody.Username) {
		log.Println("Username format error")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "用户名只能输入5-20个包含字母、数字或下划线的字串",
			"error_code": "0022",
		})
		return
	}

	// 校验密码
	if !utils.CheckPwd(signINBody.Pwd) {
		log.Println("Pwd format error")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "密码6位-16位，至少包含一个数字字母",
			"error_code": "0023",
		})
		return
	}

	// 4. 数据库查询密码，看是否和发过来的相同
	uInfo, err := tc.GetAccount(signINBody.Username)
	if err != nil && err != sql.ErrNoRows {
		log.Println("login db err:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "process error, please try again.",
			"error_code": "003",
		})
		return
	}
	if err == sql.ErrNoRows {
		log.Println("no found this user")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":      "User does not exist.",
			"error_code": "0024",
		})
		return
	}

	log.Printf("Login pwd: %s, loginBody pwd is %s", uInfo.Pwd, signINBody.Pwd)
	if err != nil || uInfo.Pwd != signINBody.Pwd {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":      "User password is wrong.",
			"error_code": "0025",
		})
		return
	}

	// 5. 插入session，由于有用name有用id来验证的，所以设置两个key
	// 先取session，不用判断这里不用判断session是否过期，因为，你已经login请求了，说明新建立一个session，直接更新session
	sId, err := ss.InsertSessionInfo(uInfo)
	if err != nil {
		log.Println("login update session err: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Process error, please try again later.",
		})
		return
	}

	// 6. 返回登录成功的消息
	c.JSON(http.StatusOK, gin.H{
		"success":    "true",
		"session_id": sId,
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

	// 2. 验证session
	if !ss.ValidateAccountSession(c.Request, signOutBody.Username) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":      "Session is not exist.",
			"error_code": "401",
		})
		return
	}
	// 3. 删除session
	if err := ss.DeleteSessionInfo(ss.GetSessionId(c.Request), signOutBody); err != nil {
		log.Printf("SignOut delete session db error : %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "Process error, please try again later.",
			"error_code": "003",
		})
		return
	}

	// 4. 返回消息
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg":     "SignOut is successful",
	})
}
