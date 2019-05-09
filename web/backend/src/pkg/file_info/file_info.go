/*
@Time : 2019/4/28 19:51 
@Author : yanKoo
@File : file_info
@Software: GoLand
@Description:
*/
package file_info

import (
	"model"
	"server/common/src/db"
)

// 增加文件信息
func AddFileInfo(fc *model.FileContext) error {
	stmtIns, err := db.DBHandler.Prepare("INSERT INTO file_info (uid, f_name, f_size, f_upload_t, f_mdf, fdfs_id) VALUES (?, ?, ?, ?, ?, ?) ")
	if err != nil {
		return err
	}

	if _, err := stmtIns.Exec(fc.UserId, fc.FileName, fc.FileSize, fc.FileUploadTime, fc.FileMD5, fc.FileFastId); err != nil {
		return err
	}

	defer stmtIns.Close()
	return nil
}

// 获取文件信息
func GetFileInfo(uId int32) (*model.FileContext, error) {
	stmtOut, err := db.DBHandler.Prepare("SELECT f_name, f_size, f_upload_t, f_mdf, fdfs_id FROM file_info WHERE uid = ? ")
	if err != nil {
		return nil, err
	}

	fc := &model.FileContext{}
	if err = stmtOut.QueryRow(uId).Scan(&fc.FileName, &fc.FileSize, &fc.FileUploadTime, &fc.FileMD5, &fc.FileFastId); err != nil {
		return nil, err
	}

	defer stmtOut.Close()
	return fc, nil
}
