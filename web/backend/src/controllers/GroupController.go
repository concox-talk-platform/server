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
	tg "pkg/group"
)

// 创建群组
func CreateGroup(c *gin.Context) {
	gInfo := &model.GroupList{}

	if err := c.BindJSON(gInfo); err != nil {
		log.Printf("json parse fail , error : %s", err)
		c.JSON(http.StatusBadRequest, model.ErrorRequestBodyParseFailed)
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
	if err := tg.CreateGroupByWeb(gInfo); err != nil {
		log.Printf("create group fail , error: %s", err)
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
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
