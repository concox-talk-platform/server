/*
@Time : 2019/4/4 14:18 
@Author : yanKoo
@File : web_device_impl
@Software: GoLand
@Description: 实现web端需要用到的关于设备管理需要用到的GRpc接口
*/
package server

import (
	pb "api/talk_cloud"
	"context"
	"errors"
	"log"
	"model"
	"net/http"
	td "pkg/device"
)

// 批量导入设备
func (wssu *WebServiceServerImpl) ImportDeviceByRoot(ctx context.Context, req *pb.ImportDeviceReq) (*pb.ImportDeviceResp, error) {
	// 设备串号和账户id进行校验
	log.Println("start Import DeviceByRoot")
	for _, v := range req.DeviceImei {
		if v == "" {
			return &pb.ImportDeviceResp{
				Result: &pb.Result{
					Msg:  "Import device Imei can't be empty! please try again later.",
					Code: http.StatusUnprocessableEntity,
				},
			}, errors.New("import device imei is empty")
		}
	}

	// 只有root用户有权限导入设备 TODO 鉴权 应该有一个root独有的签名
	if req.AccountId != 1 {
		return &pb.ImportDeviceResp{
			Result: &pb.Result{
				Msg:  "Only the root account can import devices.",
				Code: http.StatusUnprocessableEntity,
			},
		}, errors.New("account is unauthorized")
	}

	devices := make([]*model.User, 0)
	for _, v := range req.DeviceImei {
		devices = append(devices, &model.User{
			IMei:      v,
			UserName:  string([]byte(v)[7:len(v)]),
			PassWord:  string([]byte(v)[7:len(v)]),
			AccountId: int(req.GetAccountId()),
		})
	}

	if err := td.ImportDevice(devices); err != nil {
		return &pb.ImportDeviceResp{
			Result: &pb.Result{
				Msg:  "Import device error, please try again later.",
				Code: http.StatusInternalServerError,
			},
		}, err
	}

	return &pb.ImportDeviceResp{
		Result: &pb.Result{
			Msg:  "import device successful.",
			Code: http.StatusOK,
		},
	}, nil
}
