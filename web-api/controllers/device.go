/*
@Time : 2019/3/29 15:33
@Author : yanKoo
@File : DeviceController
@Software: GoLand
@Description: 超级管理员导入设备，调用mysql的GRPC的server端的方法
*/
package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	pb "server/grpc-server/api/talk_cloud"
	cfgWs "server/web-api/configs/web_server"
	"server/web-api/grpc_client_pool"
	"server/web-api/log"
	"server/web-api/model"
	"server/web-api/service"
	"server/web-api/utils"
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
		log.Log.Printf("%s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "Request body is not correct.",
			"error_code": "001",
		})
		return
	}

	conn, err := grpc_client_pool.GetConn(cfgWs.GrpcAddr)
	if err != nil {
		log.Log.Printf("grpc.Dial err : %v", err)
	}
	log.Log.Printf("%+v", aiDReq)
	webCli := pb.NewWebServiceClient(conn)

	var errDevices []*model.Device
	var duliDevices []*model.Device
	var errIdx int
	var duliIdx int
	dinfo := make([]*pb.DeviceInfo, 0)
	for _, v := range aiDReq.Devices {
		// 校验imei
		if utils.CheckIMei(v.IMei) {
			// imei查重
			if r, err := webCli.SelectDeviceByImei(context.Background(), &pb.ImeiReq{Imei: v.IMei}); err != nil {
				log.Log.Info("Select id by imei with error in web: ", err)
			} else {
				if r.Id > 0 {
					v.Id = duliIdx
					duliDevices = append(duliDevices, v)
					duliIdx++
					continue
				}
				dinfo = append(dinfo, &pb.DeviceInfo{
					IMei:       v.IMei,
					DeviceType: v.DeviceType,
					ActiveTime: v.ActiveTime,
					SaleTime:   v.SaleTime,
				})
			}
		} else {
			v.Id = errIdx
			errDevices = append(errDevices, v)
			errIdx++
		}

	}
	if len(dinfo) == 0 {
		// 返回格式不正确的数据
		c.JSON(http.StatusOK, gin.H{
			"error":        "Import some device error, Please try again later.",
			"err_devices":  errDevices,
			"deli_devices": duliDevices,
			"error_code":   "422",
		})
		return
	}

	log.Log.Println("ImportDeviceByRoot start rpc")
	res, err := webCli.ImportDeviceByRoot(context.Background(), &pb.ImportDeviceReq{
		AccountId: 1,
		Devices:   dinfo,
	})
	if err != nil {
		log.Log.Println("Import device error : ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "Import device error, Please try again later.",
			"error_code": "500",
			"msg":        err,
		})
		return
	}

	if len(errDevices) != 0 {
		// 返回格式不正确的数据
		c.JSON(http.StatusOK, gin.H{
			"error":      "Import some device error, Please try again later.",
			"devices":    errDevices,
			"error_code": "422",
		})

		if len(dinfo) == 0 {
			return
		}
	} else {
		c.JSON(int(res.Result.Code), gin.H{
			"err_devices":  errDevices,
			"deli_devices": duliDevices,
			"msg":          res.Result.Msg,
		})
	}
}

func UpdateDeviceInfo(c *gin.Context) {
	d := &pb.DeviceUpdate{}
	if err := c.BindJSON(d); err != nil {
		log.Log.Printf("%s", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "Request body is not correct.",
			"error_code": "001",
		})
		return
	}

	// 校验参数信息 ：校首先必须要有id，其次是每个参数的合法性，首先都不允许为空
	if d.LoginId == 0 {
		log.Log.Printf("account id is nil")
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

	log.Log.Println("ImportDeviceByRoot start rpc")
	conn, err := grpc_client_pool.GetConn(cfgWs.GrpcAddr)
	if err != nil {
		log.Log.Printf("grpc.Dial err : %v", err)
	}
	log.Log.Printf("%+v", d)
	webCli := pb.NewWebServiceClient(conn)
	res, err := webCli.UpdateDeviceInfo(context.Background(), &pb.UpdDInfoReq{
		DeviceInfo: d,
	})
	if err != nil {
		log.Log.Println("Import device error : ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "Update device error, Please try again later.",
			"error_code": "500",
		})
		return
	}

	c.JSON(int(res.Res.Code), gin.H{
		"msg": res.Res.Msg,
	})
}