/*
@Time : 2019/5/9 20:57 
@Author : yanKoo
@File : misc
@Software: GoLand
@Description:
*/
package utils

import (
	cfgComm "configs/common"
	"strconv"
	"time"
)

func ConvertTimeUnix(t uint64) string {
	return time.Unix(int64(t), 0).Format(cfgComm.TimeLayout)
}

func UnixStrToTimeFormat(tStr string) string {
	t, _ := strconv.ParseInt(tStr, 10, 64)
	return time.Unix(t, 0).Format(cfgComm.TimeLayout)
}

func FormatStrength(first, second, third, fourth int32) string {
	return strconv.FormatInt(int64(first), 10) + "," +
		strconv.FormatInt(int64(second), 10) + "," +
		strconv.FormatInt(int64(third), 10) + "," +
		strconv.FormatInt(int64(fourth), 10)
}
