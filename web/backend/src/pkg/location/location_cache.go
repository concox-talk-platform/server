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
	"configs"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	USR_DATA_KEY_FMT = "usr:%d:data"
)

func MakeUserDataKey(uid int32) string {
	return fmt.Sprintf(USR_DATA_KEY_FMT, uid)
}

func UpdateUserLocationInCache(req *pb.ReportDataReq, redisCli redis.Conn) error {
	// 向缓存添加用户信息数据
	if redisCli == nil {
		return errors.New("redis conn is nil")
	}
	defer redisCli.Close()

	if _, err := redisCli.Do("HMSET", MakeUserDataKey(req.DeviceInfo.Id),
		"local_time", convertTimeUnix(req.LocationInfo.GpsInfo.LocalTime),
		"lon", req.LocationInfo.GpsInfo.Longitude,
		"lat", req.LocationInfo.GpsInfo.Latitude,
		"speed", req.LocationInfo.GpsInfo.Speed,
		"course", req.LocationInfo.GpsInfo.Course,
	); err != nil {
		return errors.New("hSet failed with error: " + err.Error())
	}

	return nil
}

// TODO 给web端推送数据
func GetUserLocationInCache(uId int32, rd redis.Conn) (*pb.GPSHttpResp, error) {
	defer rd.Close()

	value, err := redis.Values(rd.Do("HMGET", MakeUserDataKey(uId),
		"lon", "lat", "speed", "course", "local_time"))
	if err != nil {
		fmt.Println("hmget failed", err.Error())
	}
	//log.Printf("Get group %d user info value string  : %s from cache ", gid, value)

	var valueStr string
	var gpsData *pb.GPSHttpResp
	resStr := make([]string, 0)
	for _, v := range value {
		if v != nil {
			valueStr = string(v.([]byte))
			resStr = append(resStr, valueStr)
		} else {
			break // redis找不到，去数据库加载
		}
	}
	log.Printf("Get user %d  gps info : %v from cache", uId, resStr)
	if value[0] != nil { // 只要任意一个字段为空就是没有这个数据
		lon, err := strconv.ParseFloat(resStr[0], 128)
		if err != nil {
			log.Printf("convent lon error:%v", err)
		}
		lat, err := strconv.ParseFloat(resStr[1], 128)
		if err != nil {
			log.Printf("convent lat error:%v", err)
		}

		speed, err := strconv.ParseFloat(resStr[2], 64)
		if err != nil {
			log.Printf("convent speed error:%v", err)
		}

		course, err := strconv.ParseFloat(resStr[3], 64)
		if err != nil {
			log.Printf("convent course error:%v", err)
		}

		lTimeT, err := time.Parse(configs.TimeLayout, resStr[4])
		if err != nil {
			log.Printf("convent lTime error:%v", err)
		}

		gpsData = &pb.GPSHttpResp{
			Uid: uId,
			Res: &pb.Result{Msg: "", Code: http.StatusOK},
			GpsInfo: &pb.GPS{
				LocalTime: uint64(lTimeT.Unix()),
				Longitude: lon,
				Latitude:  lat,
				Speed:     float32(speed),
				Course:    float32(course),
			},
		}
	} else {
		// 去数据库查找，返回空
	}
	return gpsData, nil
}
