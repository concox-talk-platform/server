package talk_cloud

import (
	"context"
)

type TalkCloudService struct {}

func (serv *TalkCloudService) AppRegister(context.Context, *AppRegReq) (*AppRegRsp, error) {
	return nil, nil
}

func (serv *TalkCloudService) DeviceRegister(context.Context, *DeviceRegReq) (*DeviceRegRsp, error) {
	return nil, nil
}

func (serv *TalkCloudService) Login(context.Context, *LoginReq) (*LoginRsp, error) {
	return nil, nil
}

func (serv *TalkCloudService) Logout(context.Context, *LogoutReq) (*LogoutRsp, error) {
	return nil, nil
}