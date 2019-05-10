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
	cfgWs "configs/web_server"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"model"
	"net/http"
	"service"
	"service/grpc_client_pool"
)

// 导入设备
func ImportDeviceByRoot(c *gin.Context) {
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

	log.Println("ImportDeviceByRoot start rpc")
	conn, err := grpc_client_pool.GetConn(cfgWs.GrpcAddr)
	if err != nil {
		log.Printf("grpc.Dial err : %v", err)
	}
	log.Printf("%+v", aiDReq)
	webCli := pb.NewWebServiceClient(conn)
	res, err := webCli.ImportDeviceByRoot(context.Background(), &pb.ImportDeviceReq{
		AccountId: 1,
		Devices:   aiDReq.Devices,
	})
	if err != nil {
		log.Println("Import device error : ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "Import device error, Please try again later.",
			"error_code": "500",
		})
		return
	}

	c.JSON(int(res.Result.Code), gin.H{
		"msg": res.Result.Msg,
	})
}

func UpdateDeviceInfo(c *gin.Context) {
	d := &pb.DeviceUpdate{}
	if err := c.BindJSON(d); err != nil {
		log.Printf("%s", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "Request body is not correct.",
			"error_code": "001",
		})
		return
	}

	// 校验参数信息 ：校首先必须要有id，其次是每个参数的合法性，首先都不允许为空
	if d.LoginId == 0 {
		log.Printf("account id is nil")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "The account id cannot be empty",
			"error_code": "003",
		})
		return
	}

	// 使用session来校验用户
	if !service.ValidateAccountSession(c.Request, int(d.LoginId)) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":      "session is not right.",
			"error_code": "006",
		})
		return
	}

	log.Println("ImportDeviceByRoot start rpc")
	conn, err := grpc_client_pool.GetConn(cfgWs.GrpcAddr)
	if err != nil {
		log.Printf("grpc.Dial err : %v", err)
	}
	log.Printf("%+v", d)
	webCli := pb.NewWebServiceClient(conn)
	res, err := webCli.UpdateDeviceInfo(context.Background(), &pb.UpdDInfoReq{
		DeviceInfo: d,
	})
	if err != nil {
		log.Println("Import device error : ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "Import device error, Please try again later.",
			"error_code": "500",
		})
		return
	}

	c.JSON(int(res.Res.Code), gin.H{
		"msg": res.Res.Msg,
	})
}
