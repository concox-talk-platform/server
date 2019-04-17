/*
@Time : 2019/4/11 10:30 
@Author : yanKoo
@File : file_controller
@Software: GoLand
@Description:
*/
package controllers

import (
	pb "api/talk_cloud"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"model"
	"net/http"
	"os"
	"service"
	"service/client_pool"
	"strconv"
)

const (
	FILE_BASE_DIR   = "/opt/web/jimi_talk.com/data/"
	FILE_BASE_URL   = "http://ptt.jimilab.com:81/data/"
	MAX_UPLOAD_SIZE = 1024 * 1024 * 500 // 1024 byte * 1024 * 500 = 500mb
)

// 进行文件大小,存在等判断，body里面等参数的判断
func uploadFilePre(c *gin.Context) (int32, error) {
	r := c.Request
	// 判断文件大小
	r.Body = http.MaxBytesReader(c.Writer, r.Body, MAX_UPLOAD_SIZE)
	if err := c.Request.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "File is too big" + " :" + err.Error(),
		})
		return -1, nil
	}

	// 判断文件类型  TODO
	fileType := int32(1)
	return fileType, nil
}

// 文件存储
func fileStore(c *gin.Context) (string, *model.UpLoadFileParmas, error) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Println("err: ", err)
	} else {
		log.Println("file:", file)
	}

	// 获取上传文件所带参数
	fParams := &model.UpLoadFileParmas{}
	if err := json.Unmarshal([]byte(c.Request.FormValue("params")), fParams); err != nil {
		log.Println("upload file unmarshal params error: ", err)
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
	}

	// 创建多级目录
	uIMDir := FILE_BASE_DIR + strconv.Itoa(fParams.Id) + "/"
	if err := os.MkdirAll(uIMDir, os.ModePerm); err != nil {
		log.Println("upload file mkdir error: ", err)
		c.JSON(http.StatusInternalServerError, model.ErrorDBError)
	}
	//写入文件
	fName := header.Filename
	out, err := os.Create(uIMDir + fName)
	if err != nil {

	}
	defer func() {
		_ = out.Close()
	}()
	_, err = io.Copy(out, file)
	if err != nil {

	}

	return FILE_BASE_URL + strconv.Itoa(fParams.Id), fParams, nil
}

func UploadFile(c *gin.Context) {
	log.Println("enter upload file.")
	fileType, err := uploadFilePre(c)
	if err != nil {
		return
	}

	// 保存文件
	fUrl, fParams, err := fileStore(c)
	if err != nil {
		return
	}

	log.Println("url: ", fUrl, "fParams: ", fParams)

	// TODO 调用调用GRPC接口，转发数据

	conn, err := client_pool.GetConn(service.Addr)
	if err != nil {
		log.Printf("grpc.Dial err : %v", err)
	}
	webCli := pb.NewTalkCloudClient(conn)

	res, err := webCli.ImMessagePublish(context.Background(), &pb.ImMsgReqData{
		Id:           int32(fParams.Id),
		ReceiverType: int32(fParams.ReceiverType),
		ReceiverId:   int32(fParams.ReceiverId),
		ResourcePath: fUrl,
		MsgType:      fileType,
	})
	if err != nil {
		log.Println("get data push stream error: ", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"msg":  "Uploaded successfully",
		"code": res.Result.Code,
	})
}
