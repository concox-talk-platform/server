/*
@Time : 2019/4/29 10:37 
@Author : yanKoo
@File : common
@Software: GoLand
@Description:
*/
package common

import (
	"github.com/go-ini/ini"
	"log"
	"os"
)

var (
	TimeLayout        string // 时间模板
)

func init() {
	cfg, err := ini.Load("common.ini") // 编译之后的执行文件所在位置的相对位置
	if err != nil {
		log.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	TimeLayout = cfg.Section("time").Key("layout").String()
}
