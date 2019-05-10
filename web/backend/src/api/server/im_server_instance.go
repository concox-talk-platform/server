/*
@Time : 2019/4/12 19:29
@Author : yanKoo
@File : talk_cloud_app_login_impl
@Software: GoLand
@Description:
*/
package server

import (
	pb "api/talk_cloud"
	cfgComm "configs/common"
	cfgGs "configs/grpc_server"
	"context"
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"log"
	"model"
	"net/http"
	tfi "pkg/file_info"
	tgc "pkg/group_cache"
	tuc "pkg/user_cache"
	"server/common/src/cache"
	"strconv"
	"time"
	"utils"
)

// 文件上传等im消息
type simpleImClientFuncImpl struct{}

func (simpleImClientFuncImpl) Dispatcher(dc *DataContext, ds DataSource) {
	imMessagePublishDispatcher(dc, ds)
}
func (simpleImClientFuncImpl) DispatcherScheduler(dc *DataContext, longLived bool) {
	dispatcherScheduler(dc, longLived)
}

// sos消息
type sosImImpl struct{}

func (sosImImpl) Dispatcher(dc *DataContext, ds DataSource) {
	notifyToOther(dc, ds.(*pb.ReportDataReq).DeviceInfo.Id, SOS_MSG)
}
func (sosImImpl) DispatcherScheduler(dc *DataContext, longLived bool) {
	dispatcherScheduler(dc, longLived)
}

// Im主要的推送
type imClientFuncImpl struct{}

func (imClientFuncImpl) Dispatcher(dc *DataContext, ds DataSource) {
	pushDataDispatcher(dc, ds)
}
func (imClientFuncImpl) DispatcherScheduler(dc *DataContext, longLived bool) {
	dispatcherScheduler(dc, longLived)
}

// Ptt
type pttImMsgImpl struct{}
type interphoneMsg struct {
	Uid       string `json:"uid"`
	MsgType   string `json:"m_type"`
	Md5       string `json:"md5"`
	GId       string `json:"grp_id"`
	FilePath  string `json:"file_path"`
	Timestamp string `json:"timestamp"`
}

func (pttImMsgImpl) Dispatcher(dc *DataContext, ds DataSource) {
	// TODO redis获取对讲音频信息，分发
	redisCli := cache.GetRedisClient()
	pttD := make(chan string, 1)
	go func() {
		for {
			//log.Printf("Start get ptt msg form redis")
			value, err := redis.Strings(redisCli.Do("blpop", cfgGs.PttMsgKey, cfgGs.PttWaitTime))
			if err != nil {
				//log.Println("blpop failed:", err.Error())
			}
			if value != nil {
				pttD <- value[1]
				log.Printf("Get ptt msg from redis: %s", value[1])
			}
		}
	}()

	var tasks [] string
	var executor = CreatePttDispatcher(dc)
	//tick := time.NewTicker(time.Second * time.Duration(5))
	for {
		var activeExecu chan string
		var activeTask string
		if len(tasks) > 0 {
			activeExecu = executor
			activeTask = tasks[0]
		}
		select {
		case t := <-pttD:
			tasks = append(tasks, t)
		case activeExecu <- activeTask:
			tasks = tasks[1:]
		//case <-tick.C:
		//	log.Printf("now ptt task queue len:%d", len(tasks))
		}
	}
}

func CreatePttDispatcher(dc *DataContext) chan string {
	tc := make(chan string)
	go pttMidHandler(tc, dc)
	return tc
}

func pttMidHandler(c chan string, dc *DataContext) {
	go func() {
		for {
			m := <-c
			log.Printf("Will send Ptt msg%s", m)
			pttMsg := &interphoneMsg{}
			if err := json.Unmarshal([]byte(m), pttMsg); err != nil {
				log.Printf("Interphone ppt msg json unmarshal fail with error :%+v", err)
			}

			if pttMsg != nil && pttMsg.Uid != "" && pttMsg.GId != "" {
				pttMsgDispatcher(dc, pttMsg)
			}
		}
	}()
}

func pttMsgDispatcher(dc *DataContext, pttMsg *interphoneMsg) {
	uId, _ := strconv.ParseInt(pttMsg.Uid, 10, 64)
	imU, err := tuc.GetUserFromCache(int32(uId))
	if err != nil {
		log.Printf("pttImMsgImpl Dispatcher GetUserFromCache error: %+v", err)
	}

	gId, _ := strconv.ParseInt(pttMsg.GId, 10, 64)
	imG, err := tgc.GetGroupInfoFromCache(int32(gId), cache.GetRedisClient())
	if err != nil {
		log.Printf("pttImMsgImpl Dispatcher GetGroupInfoFromCache error: %+v", err)
	}

	if imU != nil && imG != nil && imU.Name != "" && imG.GroupName != "" && pttMsg.MsgType == "ptt" {
		fType, fTStr := utils.GetImFileType(pttMsg.FilePath)
		log.Printf("Get ptt file type:%d, %s", fType, fTStr)
		fContext := &model.FileContext{
			UserId:         int(uId),
			FilePath:       cfgGs.FILE_BASE_URL + pttMsg.FilePath,
			FileType:       fType,
			FileName:       pttMsg.FilePath,
			FileSize:       0, // 文件大小 TODO janus没有返回
			FileMD5:        pttMsg.Md5,
			FileFastId:     pttMsg.FilePath,
			FileUploadTime: time.Now().Format(cfgComm.TimeLayout),
		}
		// 记录存储到mysql
		if err := tfi.AddFileInfo(fContext); err != nil {
			log.Printf("pttImMsgImpl Dispatcher Add file info to mysql error: %s", err.Error())
		}

		imMessagePublishDispatcher(dc, &pb.ImMsgReqData{
			Id:           imU.Uid,
			SenderName:   imU.Name,
			ReceiverType: IM_MSG_FROM_UPLOAD_RECEIVER_IS_GROUP,
			ReceiverId:   imG.Gid,
			ResourcePath: pttMsg.FilePath,
			MsgType:      fType,
			ReceiverName: imG.GroupName,
			SendTime:     pttMsg.Timestamp, // TODO 时间戳转换
		})
	}
}

func (pttImMsgImpl) DispatcherScheduler(dc *DataContext, longLived bool) {
	dispatcherScheduler(dc, longLived)
}

// 分发登录返回数据、IM离线数据、IM离线数据、Heartbeat
func (tcs *TalkCloudServiceImpl) DataPublish(srv pb.TalkCloud_DataPublishServer) error {
	c := Client{
		WorkType:  WORK_BY_GORONTINE,
		Dc:        NewDataContent(),
		Ds:        srv,
		LongLived: true,
		Cf:        imClientFuncImpl{},
	}
	c.Run()

	// 重复登录就直接返回
	uid := <-c.Dc.ExceptionalLogin
	err := srv.Send(&pb.StreamResponse{
		Res: &pb.Result{
			Msg:  "The user with id " + strconv.FormatInt(int64(uid), 10) + " is login already. please try again",
			Code: http.StatusUnauthorized,
		},
	})
	return err
}

// 上传文件方式产生的IM数据推送
func (tcs *TalkCloudServiceImpl) ImMessagePublish(ctx context.Context, req *pb.ImMsgReqData) (*pb.ImMsgRespData, error) {
	c := &Client{
		WorkType:  WORK_BY_NORMAL,
		Dc:        NewDataContent(),
		Ds:        req,
		LongLived: false,
		Cf:        simpleImClientFuncImpl{},
	}
	c.Run()
	log.Printf("# %d im once done", req.Id)
	return &pb.ImMsgRespData{Result: &pb.Result{Msg: "push data done", Code: 200}, MsgCode: req.MsgCode}, nil
}

// 分发sos消息
func (tcs *TalkCloudServiceImpl) ImSosPublish(ctx context.Context, req *pb.ReportDataReq) (*pb.ImMsgRespData, error) {
	c := &Client{
		WorkType:  WORK_BY_NORMAL,
		Dc:        NewDataContent(),
		Ds:        req,
		LongLived: true,
		Cf:        sosImImpl{},
	}
	c.Run()
	log.Printf("# %d im sos once done", req.Uid)
	return &pb.ImMsgRespData{Result: &pb.Result{Msg: "push data done", Code: 200}}, nil
}

// 分发对讲音频消息 TODO
func JanusPttMsgPublish() {
	c := &Client{
		WorkType:  WORK_BY_GORONTINE, // 持续获取数据分发，所以用一个协程挂起来分发
		Dc:        NewDataContent(),
		Ds:        nil, // 因为是去redis获取
		LongLived: true,
		Cf:        pttImMsgImpl{},
	}
	c.Run()
}
