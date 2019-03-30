/*
@Time : 2019/3/28 15:33
@Author : yanKoo
@File : TalkCloudRegisterImpl
@Software: GoLand
@Description: 目前主要供web端调用
*/
package api

import (
	pb "api/talk_cloud"
	"context"
	"errors"
	"log"
	"model"
	"net/http"
	td "pkg/device"
	"strconv"
)

type WebServiceServerImpl struct {
}

func (wssu *WebServiceServerImpl) ImportDeviceByRoot(ctx context.Context, req *pb.ImportDeviceReq) (*pb.ImportDeviceResp, error) {
	// 设备串号和账户id进行校验
	log.Println("start Import DeviceByRoot")
	for _, v := range req.DeviceImei {
		if v == "" {
			return &pb.ImportDeviceResp{
				Result: &pb.Result{
					Msg:       "Import device Imei can't be empty! please try again later.",
					StateCode: http.StatusUnprocessableEntity,
				},
			}, errors.New("import device imei is empty")
		}
	}

	// 只有root用户有权限导入设备 TODO 鉴权 应该有一个root独有的签名
	if req.AccountId != 1 {
		return &pb.ImportDeviceResp{
			Result: &pb.Result{
				Msg:       "Only the root account can import devices.",
				StateCode: http.StatusUnprocessableEntity,
			},
		}, errors.New("account is unauthorized")
	}

	devices := make([]*model.User, 0)
	for _, v := range req.DeviceImei {
		devices = append(devices, &model.User{
			IMei:  v,
			UserName: string([]byte(v)[7:len(v)]),
			PassWord: string([]byte(v)[7:len(v)]),
			AccountId: strconv.FormatInt(int64(req.GetAccountId()), 10),
		})
	}

	if err := td.ImportDevice(devices); err != nil {
		return &pb.ImportDeviceResp{
			Result: &pb.Result{
				Msg:       "Import device error, please try again later.",
				StateCode: http.StatusInternalServerError,
			},
		}, err
	}

	return &pb.ImportDeviceResp{
		Result: &pb.Result{
			Msg:       "import device successful.",
			StateCode: http.StatusOK,
		},
	}, nil
}

// TODO 获取经销商的下级用户
func (wssi *WebServiceServerImpl) GetAccountClass(ctx context.Context, req *pb.AccountClassReq) (*pb.AccountClassResp, error) {

	if req.Name == "bob" {
		return &pb.AccountClassResp{Id: 1}, nil
	} else {
		return &pb.AccountClassResp{Id: -1}, nil
	}
}
