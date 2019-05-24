/**
* @Author: yanKoo
* @Date: 2019/3/11 10:48
* @Description: 处理请求的业务逻辑
 */
package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	pb "server/grpc-server/api/talk_cloud"
	cfgWs "server/web-api/configs/web_server"
	tg "server/web-api/dao/group"
	"server/web-api/db"
	"server/web-api/grpc_client_pool"
	"server/web-api/log"
	"server/web-api/model"
	"server/web-api/service"
	"strconv"
)

// @Summary web更新群组中的设备
// @Description logout by account name and pwd, 请求头中Authorization参数设置为登录时返回的sessionId
// @Accept  json
// @Produce  json
// @Param Authorization header string true "登录时返回的sessionId"
// @Param body body model.GroupList true "更新群组中的设备"
// @Success 200 {string} json "{"success":"true","msg": resUpd.ResultMsg.Msg}"
// @Router /group/devices/update [post]
func UpdateGroupDevice(c *gin.Context) {
	gList := &model.GroupList{}
	if err := c.BindJSON(gList); err != nil {
		log.Log.Printf("json parse fail , error : %s", err)
		c.JSON(http.StatusBadRequest, model.ErrorRequestBodyParseFailed)
		return
	}
	// 使用session来校验用户
	if !service.ValidateAccountSession(c.Request, gList.GroupInfo.AccountId) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":      "session is not right.",
			"error_code": "006",
		})
		return
	}
	// TODO 校验更新群组信息的参数合法性
	if gList.GroupInfo.GroupName == "" || gList.GroupInfo.AccountId == 0 || len(gList.DeviceIds) == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "You need at least the group name, the account id, and at least one device id",
			"error_code": "001",
		})
		return
	}

	// 更新群组
	log.Log.Println("update group start rpc")
	conn, err := grpc_client_pool.GetConn(cfgWs.GrpcAddr)
	if err != nil {
		log.Log.Printf("grpc.Dial err : %v", err)
	}
	webCli := pb.NewWebServiceClient(conn)

	deviceIds := make([]int64, 0)
	for _, v := range gList.DeviceIds {
		deviceIds = append(deviceIds, int64(v))
	}

	deviceInfos := make([]*pb.Member, 0)
	for _, v := range gList.DeviceInfo {
		vMap := v.(map[string]interface{})
		log.Log.Println((vMap["id"]).(float64))
		deviceInfos = append(deviceInfos, &pb.Member{
			Id:       int32((vMap["id"]).(float64)),
			IMei:     (vMap["imei"]).(string),
			UserName: (vMap["user_name"]).(string),
			//NickName:  (vMap["nick_name"]).(nil),
			Pwd: (vMap["password"]).(string),
		})
	}
	log.Log.Println("group member update :gList.GroupInfo.Id :", gList.GroupInfo.Id)
	status, _ := strconv.Atoi(gList.GroupInfo.Status)
	resUpd, err := webCli.UpdateGroup(context.Background(), &pb.UpdateGroupReq{
		DeviceIds:   deviceIds,
		DeviceInfos: deviceInfos,
		GroupInfo: &pb.Group{
			Id:        int32(gList.GroupInfo.Id),
			GroupName: gList.GroupInfo.GroupName,
			AccountId: int32(gList.GroupInfo.AccountId),
			Status:    int32(status)},
	})
	if err != nil {
		log.Log.Printf("Update group fail , error: %s", err)
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	}

	log.Log.Println(resUpd)
	c.JSON(http.StatusOK, gin.H{
		"result": "success",
		"msg":    resUpd.ResultMsg.Msg,
	})
}

// @Summary web创建群组
// @Description web创建群组, 请求头中Authorization参数设置为登录时返回的sessionId
// @Accept  json
// @Produce  json
// @Param Authorization header string true "登录时返回的sessionId"
// @Param body body model.GroupList true "创建群组"
// @Success 200 {string} json "{"success":"true","msg": resUpd.ResultMsg.Msg}"
// @Router /group [post]
func CreateGroup(c *gin.Context) {
	gList := &model.GroupList{}
	if err := c.BindJSON(gList); err != nil {
		log.Log.Printf("json parse fail , error : %s", err)
		c.JSON(http.StatusBadRequest, model.ErrorRequestBodyParseFailed)
		return
	}
	// 使用session来校验用户
	if !service.ValidateAccountSession(c.Request, gList.GroupInfo.AccountId) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":      "session is not right.",
			"error_code": "006",
		})
		return
	}
	// TODO 校验创建群组信息的参数合法性
	if gList.GroupInfo.GroupName == "" || gList.GroupInfo.AccountId == 0 || len(gList.DeviceIds) == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":      "You need at least the group name, the account id, and at least one device id",
			"error_code": "001",
		})
		return
	}

	// 组名查重
	res, err := tg.CheckDuplicateGName(gList.GroupInfo)
	if err != nil {
		log.Log.Printf("CheckDuplicateGName fail , error: %s", err)
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	}
	if res > 0 {
		log.Log.Printf("CheckDuplicateGName error: %s", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg":  "group name duplicate",
			"code": "422",
		})
		return
	}

	// 创建群组
	log.Log.Println("start rpc")
	conn, err := grpc_client_pool.GetConn(cfgWs.GrpcAddr)
	if err != nil {
		log.Log.Printf("grpc.Dial err : %v", err)
	}
	webCli := pb.NewTalkCloudClient(conn)
	log.Log.Printf("++++++++++++++webCli: %+v", webCli)
	var deviceIds string
	for _, v := range gList.DeviceIds {
		if v == -1 {
			deviceIds = "-1"
		} else {

		}
	}

	deviceInfos := make([]*pb.Member, 0)
	for _, v := range gList.DeviceInfo {
		vMap := v.(map[string]interface{})
		log.Log.Println((vMap["id"]).(float64))
		deviceInfos = append(deviceInfos, &pb.Member{
			Id:       int32((vMap["id"]).(float64)),
			IMei:     (vMap["imei"]).(string),
			UserName: (vMap["user_name"]).(string),
			Pwd:      (vMap["password"]).(string),
		})
	}
	status, _ := strconv.Atoi(gList.GroupInfo.Status)

	log.Log.Println("gList.GroupInfo.GroupName:", gList.GroupInfo.GroupName)

	resCreate, err := webCli.CreateGroup(context.Background(), &pb.CreateGroupReq{
		DeviceIds:   deviceIds,
		DeviceInfos: deviceInfos,
		GroupInfo: &pb.Group{
			Id:        int32(gList.GroupInfo.Id),
			GroupName: gList.GroupInfo.GroupName,
			AccountId: int32(gList.GroupInfo.AccountId),
			Status:    int32(status)},
	})
	log.Log.Printf("group: %+v", resCreate.GroupInfo.GroupName)
	if err != nil {
		log.Log.Printf("create group fail , error: %s", err)
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	}

	log.Log.Println(resCreate)
	c.JSON(http.StatusOK, gin.H{
		"result":     "success",
		"group_info": resCreate.GroupInfo,
		"msg":        resCreate.Res.Msg,
	})
}


// @Summary web更新群组信息，目前只更新群组名字
// @Description web创建群组, 请求头中Authorization参数设置为登录时返回的sessionId
// @Accept  json
// @Produce  json
// @Param Authorization header string true "登录时返回的sessionId"
// @Param body body model.GroupInfo true "web更新群组"
// @Success 200 {string} json "{"success":"true","msg": resUpd.ResultMsg.Msg}"
// @Router /group/update [post]
func UpdateGroup(c *gin.Context) {
	gI := &model.GroupInfo{}
	if err := c.BindJSON(gI); err != nil {
		log.Log.Printf("json parse fail , error : %s", err)
		c.JSON(http.StatusBadRequest, model.ErrorRequestBodyParseFailed)
		return
	}

	// 使用session来校验用户
	if !service.ValidateAccountSession(c.Request, gI.AccountId) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":      "session is not right.",
			"error_code": "006",
		})
		return
	}
	// 组名查重
	res, err := tg.CheckDuplicateGName(gI)
	if err != nil {
		log.Log.Printf("CheckDuplicateGName fail , error: %s", err)
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	}
	if res > 0 {
		log.Log.Printf("CheckDuplicateGName error: %s", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg":  "group name duplicate",
			"code": "422",
		})
		return
	}

	if err := tg.UpdateGroup(gI, db.DBHandler); err != nil {
		log.Log.Printf("update group fail , error: %s", err)
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "success",
		"msg":    "Update group successfully",
	})
}

// @Summary web群组删除，目前只更新群组名字
// @Description web创建群组, 请求头中Authorization参数设置为登录时返回的sessionId
// @Accept  json
// @Produce  json
// @Param Authorization header string true "登录时返回的sessionId"
// @Param body body model.GroupInfo true "web更新群组"
// @Success 200 {string} json "{"success":"true","msg": resUpd.ResultMsg.Msg}"
// @Router /group/delete [post]
func DeleteGroup(c *gin.Context) {
	gI := &model.GroupInfo{}
	if err := c.BindJSON(gI); err != nil {
		log.Log.Printf("json parse fail , error : %s", err)
		c.JSON(http.StatusBadRequest, model.ErrorRequestBodyParseFailed)
		return
	}
	// 使用session来校验用户
	if !service.ValidateAccountSession(c.Request, gI.AccountId) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":      "session is not right.",
			"error_code": "006",
		})
		return
	}
	log.Log.Println("start rpc")
	conn, err := grpc_client_pool.GetConn(cfgWs.GrpcAddr)
	if err != nil {
		log.Log.Printf("grpc.Dial err : %v", err)
	}
	webCli := pb.NewWebServiceClient(conn)

	if _, err := webCli.DeleteGroup(context.Background(), &pb.Group{Id: int32(gI.Id)}); err != nil {
		log.Log.Printf("update group fail , error: %s", err)
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "success",
		"msg":    "Delete group successfully",
	})
}