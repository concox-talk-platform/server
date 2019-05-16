/*
@Time : 2019/4/11 10:30 
@Author : yanKoo
@File : file_controller
@Software: GoLand
@Description:
*/
package controllers

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"
	pb "server/grpc-server/api/talk_cloud"
	cfgComm "server/web-api/configs/common"
	cfgWs "server/web-api/configs/web_server"
	tfi "server/web-api/dao/file_info"
	"server/web-api/grpc_client_pool"
	"server/web-api/log"
	"server/web-api/model"
	"server/web-api/service/fdfs_client"
	"server/web-api/utils"
	"strconv"
	"time"
)

var upGrader = websocket.Upgrader{
	HandshakeTimeout: time.Duration(600),
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	MAX_UPLOAD_SIZE = 1024 * 1024 * 500 // 1024 byte * 1024 * 500 = 500mb

	FIRST_LOGIN_DATA                = 1 // 初次登录返回的数据。比如用户列表，群组列表，该用户的个人信息
	OFFLINE_IM_MSG                  = 2 // 用户离线时的IM数据
	IM_MSG_FROM_UPLOAD_OR_WS_OR_APP = 3 // APP和web通过httpClient上传的IM信息
	KEEP_ALIVE_MSG                  = 4 // 用户登录后，每隔interval秒向stream发送一个消息，测试能不能连通
	LOGOUT_NOTIFY_MSG               = 5 // 用户掉线之后，通知和他在一个组的其他成员
	LOGIN_NOTIFY_MSG                = 6 // 用户上线之后，通知和他在一个组的其他成员
	IM_SOS_MSG                      = 7 // 版本不同，暂定这里

	IM_MSG_WORKDONE  = 1
	IM_MSG_WORKWRONG = -1

	IM_TEXT_MSG         = 1 // 普通文本
	IM_IMAGE_MSG        = 2 // 图片
	IM_VOICE_MSG        = 3 // 音频文件
	IM_VIDEO_MSG        = 4 // 视频文件
	IM_PDF_MSG          = 5 // PDF文件
	IM_UNKNOWN_TYPE_MSG = 10000
)

type worker struct {
	uId        int32
	cliStream  *pb.TalkCloud_DataPublishClient
	ws         *websocket.Conn
	Data       chan interface{}
	mt         int
	WorkerDone chan int
}

var conn *grpc.ClientConn
var client *fdfs_client.Client

func init() {
	// 调用调用GRPC接口，转发数据
	var err error
	conn, err = grpc_client_pool.GetConn(cfgWs.GrpcAddr)
	if err != nil {
		log.Log.Printf("grpc.Dial err : %v", err)
	}

	client, err = fdfs_client.NewClientWithConfig()
	if err != nil {
		log.Log.Printf("Client: %+v NewClientWithConfig fastdfs error: %+v", client, err)
	}
	// defer client.Destory() TODO destory?
}

func UploadFile(c *gin.Context) {
	log.Log.Println("start upload file.")
	err := uploadFilePre(c)
	if err != nil {
		log.Log.Println("uploadFilePre error: ", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "Uploaded File params error, please try again later.", "code": 001})
		return
	}

	// 保存文件
	fContext, err := fileStore(c)
	if err != nil {
		log.Log.Println("fileStore", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Upload File fail, please try again later.", "code": 002})
		return
	}

	log.Log.Println("url: ", fContext.FilePath, "fParams: ", fContext.FileParams)

	webCli := pb.NewTalkCloudClient(conn)

	res, err := webCli.ImMessagePublish(context.Background(), &pb.ImMsgReqData{
		Id:           int32(fContext.FileParams.Id),
		SenderName:   fContext.FileParams.SenderName,
		ReceiverType: int32(fContext.FileParams.ReceiverType),
		ReceiverId:   int32(fContext.FileParams.ReceiverId),
		ResourcePath: fContext.FilePath,
		MsgType:      fContext.FileType,
		ReceiverName: fContext.FileParams.ReceiverName,
		SendTime:     fContext.FileParams.SendTime,
		MsgCode:      strconv.FormatInt(time.Now().Unix(), 10),
	})
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{"msg": "Uploaded File, please try again later.", "code": 001})
		return
	}
	log.Log.Printf("upload file success by grpc: %+v", res)

	c.JSON(http.StatusCreated, gin.H{
		"msg":          "Uploaded successfully",
		"code":         res.Result.Code,
		"MsgCode":      res.MsgCode,
		"resourcePath": fContext.FilePath,
		"resourceName": fContext.FileName,
	})
}

// 进行文件大小,存在等判断，body里面等参数的判断
func uploadFilePre(c *gin.Context) error {
	r := c.Request
	// 判断文件大小
	r.Body = http.MaxBytesReader(c.Writer, r.Body, MAX_UPLOAD_SIZE)
	if err := c.Request.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		log.Log.Println("File is too big.")
		return err
	}

	// TODO 允许重复上传文件?
	return nil
}

// 文件存储
func fileStore(c *gin.Context) (*model.FileContext, error) {
	file, header, err := c.Request.FormFile("file") // TODO 会报空针
	if err != nil {
		log.Log.Println("fileStore err: ", err)
		return nil, err
	}
	uploadT := time.Now().Format(cfgComm.TimeLayout)
	// 获取上传文件所带参数
	id, _ := strconv.ParseInt(c.Request.FormValue("id"), 10, 32)
	senderName := c.Request.FormValue("SenderName")
	receiverType, _ := strconv.ParseInt(c.Request.FormValue("ReceiverType"), 10, 64)
	receiverId, _ := strconv.ParseInt(c.Request.FormValue("ReceiverId"), 10, 64)
	receiverName := c.Request.FormValue("ReceiverName")
	sTime := c.Request.FormValue("SendTime")

	//简单做一下数据判断
	if id <= 0 || receiverId <= 0 || receiverType <= 0 {
		return nil, errors.New("file param is cant be nil")
	}

	fParams := &model.ImMsgData{
		Id:           int(id),
		SenderName:   senderName,
		ReceiverType: int(receiverType),
		ReceiverId:   int(receiverId),
		ReceiverName: receiverName,
		SendTime:     sTime,
	}
	log.Log.Printf("file params: %+v", fParams)

	//写入文件
	fName := strconv.FormatInt(int64(fParams.Id), 10) + "_" +
		strconv.FormatInt(time.Now().Unix(), 10) + "_" +
		header.Filename

	fSrc, err := ioutil.ReadAll(file)
	if err != nil {
		log.Log.Println("read file error: ", err)
		return nil, err
	}
	fileType, fExtName := utils.GetImFileType(header.Filename)

	// 先检验文件的hash值，避免重复上传
	md5h := md5.New()
	md5h.Write(fSrc)
	fMd5 := hex.EncodeToString(md5h.Sum([]byte("")))
	fmt.Printf("this file md5: %s\n", hex.EncodeToString(md5h.Sum([]byte("")))) //md5

	// 存储文件到fastdfs

	fileId, err := client.UploadByBuffer(fSrc, fExtName)
	if err != nil {
		log.Log.Println("UploadByBuffer to fastdfs error: ", err)
		return nil, err
	}
	log.Log.Printf("file size: %d ", len(fSrc))

	fContext := &model.FileContext{
		UserId:         fParams.Id,
		FilePath:       cfgWs.FILE_BASE_URL + fileId,
		FileParams:     fParams,
		FileType:       fileType,
		FileName:       fName,
		FileSize:       len(fSrc),
		FileMD5:        fMd5,
		FileFastId:     fileId,
		FileUploadTime: uploadT,
	}

	// 记录存储到mysql
	if err := tfi.AddFileInfo(fContext); err != nil {
		log.Log.Printf("Add file info to mysql error: %s", err.Error())
		return nil, err
	}

	return fContext, nil
}

// websocket与grpc交换数据
func ImPush(c *gin.Context) {
	uidStr := c.Param("accountId")
	uid, _ := strconv.Atoi(uidStr) // TODO 校验用户是否存在

	log.Log.Println("im push uid :", uid)

	// 调用调用GRPC接口，转发数据
	webCliStream, err := pb.NewTalkCloudClient(conn).DataPublish(context.Background())
	if err != nil {
		log.Log.Printf("connect grpc fail with error: %s", err.Error())
		return
	}

	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Log.Printf("connect grpc fail with error: %s", err.Error())
		return
	}
	defer func() {
		ws.Close()
		//conn.Close()
	}()

	imWorker := &worker{
		uId:        int32(uid),
		cliStream:  &webCliStream,
		ws:         ws,
		Data:       make(chan interface{}, 1),
		WorkerDone: make(chan int, 1),
	}

	// 有两种跳出循环的情况：
	// 1、web端主动关闭连接，grpc也就要不再接受数据，
	// 2、web端重复登录，TODO 放在这里判断重复登录有点不妥当，不过如果前面的登录做得好，这里不会出现这种情况，以防万一吧。
	log.Log.Println(strconv.FormatInt(int64(imWorker.uId), 10) + " ws grpc start")

	ctx, cancel := context.WithCancel(context.Background())

	// 接收web端的消息，转发给grpc
	go pushImMessage(imWorker, ctx)

	// 发送ws消息
	go sendImMessage(imWorker, ctx)

	if wd := <-imWorker.WorkerDone; wd == IM_MSG_WORKWRONG {
		_ = imWorker.ws.WriteMessage(websocket.TextMessage,
			[]byte("The connection with id:"+strconv.FormatInt(int64(imWorker.uId), 10)+
				" has been disconnected, please reconnect"))
		log.Log.Println("break******************************")
		cancel()
		return
	} else {
		// TODO grpc服务主动拒绝连接（重复登录）
	}
	//}(imWorker, &wg)
	//wg.Wait()
}

func pushImMessage(imw *worker, ctx context.Context) {
	// 发送给GRPC
	if err := (*imw.cliStream).Send(&pb.StreamRequest{
		Uid:      imw.uId,
		DataType: OFFLINE_IM_MSG,
	}); err != nil {
		imw.WorkerDone <- IM_MSG_WORKWRONG
		log.Log.Println("im message send error: ", err)
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// 读取ws中的数据，这里的数据，默认只有文本数据
				mt, message, err := imw.ws.ReadMessage()
				imw.mt = mt
				if err != nil {
					// 客户端关闭连接时也会进入
					log.Log.Printf("%d WS message read error: %s", imw.uId, err.Error())
					imw.WorkerDone <- IM_MSG_WORKWRONG // TODO
					return
				}

				log.Log.Println("ws receive msg: ", message)
				wsImMsg := &model.ImMsgData{}
				if err := json.Unmarshal(message, wsImMsg); err != nil {
					log.Log.Printf("json unmarshal fail with err :%v", err)
					// TODO  暂时忽略这条消息
					continue
				}

				// 暂时默认发过来的消息都是普通文本
				wsImMsg.MsgType = IM_TEXT_MSG

				// 发送给GRPC
				log.Log.Printf("ws will send to grpc: %+v", wsImMsg)
				if err := sendToGrpc(imw, wsImMsg); err != nil {
					imw.WorkerDone <- IM_MSG_WORKWRONG
					log.Log.Println("grpc im message send error: ", err)
					return
				}
			}
		}
	}()

	// 发送心跳
	tick := time.NewTicker(time.Second * time.Duration(cfgWs.Interval))
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			if err := (*imw.cliStream).Send(&pb.StreamRequest{
				Uid:      int32(imw.uId),
				DataType: KEEP_ALIVE_MSG,
			}); err != nil {
				imw.WorkerDone <- IM_MSG_WORKWRONG
				log.Log.Println("grpc im heartbeat message send error: ", err)
				return
			}
		}
	}
}
func sendToGrpc(imw *worker, wsImMsg *model.ImMsgData) error {
	if err := (*imw.cliStream).Send(&pb.StreamRequest{
		Uid:      int32(imw.uId),
		DataType: IM_MSG_FROM_UPLOAD_OR_WS_OR_APP,
		ImMsg: &pb.ImMsgReqData{
			Id:           int32(wsImMsg.Id),
			SenderName:   wsImMsg.SenderName,
			ReceiverId:   int32(wsImMsg.ReceiverId),
			ReceiverName: wsImMsg.ReceiverName,
			SendTime:     wsImMsg.SendTime,
			ReceiverType: int32(wsImMsg.ReceiverType),
			ResourcePath: wsImMsg.ResourcePath, // 文本消息直接放路劲这个字段
			MsgType:      int32(wsImMsg.MsgType),
		},
	}); err != nil {
		return err
	}
	return nil
}

// 从grpc stream中接收数据通过ws转发给web
func sendImMessage(imw *worker, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			resp, err := (*imw.cliStream).Recv()
			if err != nil {
				imw.WorkerDone <- IM_MSG_WORKWRONG
				log.Log.Printf("%d grpc recv message error: %s", imw.uId, err.Error())
				break
			}
			log.Log.Printf("%d web grpc client receive : %+v", imw.uId, resp)

			// 写入ws数据 二进制返回
			if resp.DataType == IM_MSG_FROM_UPLOAD_OR_WS_OR_APP {
				// 把中文转换为utf-8
				resp.ImMsgData.ResourcePath = utils.ConvertOctonaryUtf8(resp.ImMsgData.ResourcePath)

				log.Log.Printf("web grpc client receive : %+v", resp)

				// 返回JSON字符串
				err = imw.ws.WriteJSON(resp)
				if err != nil {
					imw.WorkerDone <- IM_MSG_WORKWRONG
					log.Log.Println("WS message send error:", err)
					//break
				}
			}

			if resp.DataType == OFFLINE_IM_MSG {
				// 把中文转换为utf-8
				for _, msg := range resp.OfflineImMsgResp.OfflineGroupImMsgs {
					if msg.ImMsgData != nil {
						for _, userMsg := range msg.ImMsgData {
							userMsg.ResourcePath = utils.ConvertOctonaryUtf8(userMsg.ResourcePath)
						}
					}
				}
				for _, msg := range resp.OfflineImMsgResp.OfflineSingleImMsgs {
					if msg.ImMsgData != nil {
						for _, userMsg := range msg.ImMsgData {
							userMsg.ResourcePath = utils.ConvertOctonaryUtf8(userMsg.ResourcePath)
						}
					}
				}

				log.Log.Printf("web grpc client receive : %+v", resp)

				// 返回JSON字符串
				err = imw.ws.WriteJSON(resp)
				if err != nil {
					imw.WorkerDone <- IM_MSG_WORKWRONG
					log.Log.Println("WS message send error:", err)
					//break
				}
			}

			// 掉线通知
			if resp.DataType == LOGOUT_NOTIFY_MSG || resp.DataType == LOGIN_NOTIFY_MSG || resp.DataType == IM_SOS_MSG {
				err = imw.ws.WriteJSON(resp)
				if err != nil {
					imw.WorkerDone <- IM_MSG_WORKWRONG
					log.Log.Println("WS message send error:", err)
					//break
				}
			}
		}
	}
}
