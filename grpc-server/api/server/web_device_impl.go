/*
@Time : 2019/4/4 14:18
@Author : yanKoo
@File : web_device_impl
@Software: GoLand
@Description: 实现web端需要用到的关于设备管理需要用到的GRpc接口
*/
package server

import (
	"context"
	"errors"
	"net/http"
	pb "server/grpc-server/api/talk_cloud"
	td "server/grpc-server/dao/device"
	"server/grpc-server/log"
	"server/web-api/model"
)

// 批量导入设备
func (wssu *WebServiceServerImpl) ImportDeviceByRoot(ctx context.Context, req *pb.ImportDeviceReq) (*pb.ImportDeviceResp, error) {
	// 设备串号和账户id进行校验
	log.Log.Println("start Import DeviceByRoot")
	for _, v := range req.Devices {
		if v == nil || v.IMei == "" {
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
	for _, v := range req.Devices {
		devices = append(devices, &model.User{
			IMei:       v.IMei,
			UserName:   v.IMei,
			NickName:   string([]byte(v.IMei)[12:len(v.IMei)]),
			PassWord:   string([]byte(v.IMei)[9:len(v.IMei)]),
			AccountId:  int(req.GetAccountId()),
			ParentId:   "1",
			DeviceType: v.DeviceType,
			ActiveTime: v.ActiveTime,
			SaleTime:   v.SaleTime,
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

func (wssu *WebServiceServerImpl) UpdateDeviceInfo(ctx context.Context, req *pb.UpdDInfoReq) (*pb.UpdDInfoResp, error) {
	if err := td.UpdateDeviceInfo(&model.User{IMei: req.DeviceInfo.IMei, NickName: req.DeviceInfo.NickName}); err != nil {
		return &pb.UpdDInfoResp{
			Res: &pb.Result{
				Msg:  "Update DeviceInfo device error, please try again later.",
				Code: http.StatusInternalServerError,
			},
		}, err
	}
	return &pb.UpdDInfoResp{
		Res: &pb.Result{
			Msg:  "Update DeviceInfo device successful!!!!!",
			Code: http.StatusOK,
		},
	}, nil
}
func (wssu *WebServiceServerImpl) SelectDeviceByImei(ctx context.Context, req *pb.ImeiReq) (*pb.ImeiResp, error) {
	id, err := td.SelectDeviceByImei(&model.User{IMei: req.Imei})
	resp := &pb.ImeiResp{
		Res: &pb.Result{
			Msg:  "SelectDeviceByImei error, please try again later.",
			Code: http.StatusInternalServerError,
		},
	}
	resp.Id = id
	if err != nil {
		return resp, err
	}
	return resp, nil
}