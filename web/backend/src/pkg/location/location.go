/*
@Time : 2019/4/12 16:24 
@Author : yanKoo
@File : location
@Software: GoLand
@Description: 存储设备发过来的的数据
*/
package location

import (
	pb "api/talk_cloud"
	"database/sql"
	"github.com/smartwalle/dbs"
	"log"
	cfgComm "configs/common"
	"strconv"
	"strings"
	"time"
)

/*
id bigint(12) NULL记录id
uid bigint(12) NOT NULL设备/用户的id
local_time timestamp NULLGPS数据定位时间
lng varchar(64) NOT NULL经度
lat varchar(64) NOT NULL纬度
cse_sp varchar(128) NULL航向，速度

country varchar(255) NULL基站所在国家
operator varchar(255) NULL基站的运营商所在地区
region varchar(255) NULL基站所在的地区
bs_sth int(12) NULL基站的信号强度
wifi_sth int(12) NULLwifi的信号强度
bt_sth int(12) NULL蓝牙的信号强度
create_time timestamp NOT NULL记录存入数据库的时间
*/
// 插入GPS数据

func InsertLocationData(req *pb.ReportDataReq, db *sql.DB) error {
	log.Printf("receive gps data from app: %+v", req)
	ib := dbs.NewInsertBuilder()
	ib.Table("location")

	ib.SET("uid", req.DeviceInfo.Id)
	ib.SET("local_time", convertTimeUnix(req.LocationInfo.GpsInfo.LocalTime))
	ib.SET("lng", req.LocationInfo.GpsInfo.Longitude)
	ib.SET("lat", req.LocationInfo.GpsInfo.Latitude)
	ib.SET("cse_sp", PackCourseSpeed(req.LocationInfo.GpsInfo.Course, req.LocationInfo.GpsInfo.Speed))

	if _, err := ib.Exec(db); err != nil {
		return err
	}
	return nil
}

// 打包航向和速度
func PackCourseSpeed(course, speed float32) string {
	return strconv.FormatFloat(float64(course), 'f', -1, 32) + "," +
		strconv.FormatFloat(float64(speed), 'f', -1, 32)
}

// 解析出航向和速度
func ParseCourseSpeed(cseSpeed string) (int32, float32) {
	strs := strings.Split(cseSpeed, ",")

	course, err := strconv.ParseInt(strs[0], 10, 32)
	if err != nil {
		log.Println("parse course fail with error: ", err)
	}

	speed, err := strconv.ParseFloat(strs[1], 64)
	if err != nil {
		log.Println("parse speed fail with error: ", err)
	}

	return int32(course), float32(speed)
}

func convertTimeUnix(t uint64) string {
	return time.Unix(int64(t), 0).Format(cfgComm.TimeLayout)
}
