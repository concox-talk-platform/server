/**
* @Author: yanKoo
* @Date: 2019/3/11 11:03
* @Description: 用来定义返回错误的结构体和消息格式
 */
package model

import "github.com/gin-gonic/gin"

// 错误结构体
var (
	ErrorRequestBodyParseFailed = gin.H{
		"error":      "Request body is not correct.",
		"error_code": "001",
	} // 不能解析消息体

	ErrorDBError = gin.H{
		"error":      "DB ops failed",
		"error_code": "003",
	} // 数据库操作错误

	ErrorNotAuthSession = gin.H{
		"error":      "session is not right.",
		"error_code": "006",
	} // 账户不合法，不存在
)
