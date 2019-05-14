/*
@Time : 2019/5/6 17:25 
@Author : yanKoo
@File : app_user_apk
@Software: GoLand
@Description:
*/
package server

import (
	"context"
	"net/http"
	pb "server/grpc-server/api/talk_cloud"
	cfgGs "server/grpc-server/configs/grpc_server"
	tfi "server/grpc-server/dao/file_info"
)

// 返回apk消息
func (tcs *TalkCloudServiceImpl) GetApkInfo(ctx context.Context, req *pb.ApkInfoReq) (*pb.ApkInfoResp, error) {
	resp := &pb.ApkInfoResp{Res: &pb.Result{Msg: "Get apk info successful", Code: http.StatusOK}}
	apkInfo, err := tfi.GetFileInfo(req.Uid)
	if err != nil {
		resp.Res.Code = http.StatusInternalServerError
		resp.Res.Msg = "Get apk info fail please try again later"
		return resp, nil
	}
	resp.ApkPath = cfgGs.FILE_BASE_URL + apkInfo.FileMD5
	resp.ApkVersion = apkInfo.FileName
	return resp, nil
}
