/*
@Time : 2019/3/29 15:33 
@Author : yanKoo
@File : DeviceController
@Software: GoLand
@Description: 超级管理员导入设备，调用mysql的GRPC的server端的方法
*/
package controllers

import (
	pb "api/talk_cloud"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"model"
	"net/http"
	"service"
)

// 导入设备 TODO
func ImportDeviceByRoot(c * gin.Context) {
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

	aiDReq := &model.AccountImportDeviceReq{}
	if err := c.BindJSON(aiDReq); err != nil {
		log.Printf("%s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "Request body is not correct.",
			"error_code": "001",
		})
		return
	}

	res, err := service.WebCli.ImportDeviceByRoot(context.Background(), &pb.ImportDeviceReq{
		DeviceImei:aiDReq.DeviceIMei,
		AccountId: 1,
	})
	if err != nil {
		log.Println("import device error : ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":"import device error, Please try again later.",
			"error_code":"500",
		})
		return
	}
	log.Println(res.Result.Msg, res.Result.StateCode)
	c.JSON(http.StatusOK, gin.H{
		"msg":"import device successful.",
	})
}