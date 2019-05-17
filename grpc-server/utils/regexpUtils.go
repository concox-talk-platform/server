/*
@Time : 2019/3/30 11:41
@Author : yanKoo
@File : regexpUtils
@Software: GoLand
@Description:
*/
package utils

import (
	"fmt"
	"regexp"
	"server/grpc-server/log"
	"strconv"
)

// 校验手机号码
func CheckPhone(phone string) bool {
	reg, err := regexp.Compile("^((13[0-9])|(15[^4,\\D])|(18[0,5-9]))\\d{8}$")
	if err != nil {
		log.Log.Println(err)
	}
	return reg.MatchString(phone)
}

// 校验邮箱
func CheckEmail(email string) bool {
	reg, err := regexp.Compile("\\w+(\\.\\w)*@\\w+(\\.\\w{2,3}){1,3}")
	if err != nil {
		fmt.Println(err)
	}
	return reg.MatchString(email)
}

// 校验密码：只能输入6-20个字母、数字、下划线
func CheckPwd(pwd string) bool {

	reg, err := regexp.Compile("^([a-zA-Z0-9]|[_]){6,20}$")
	if err != nil {
		fmt.Println(err)
	}
	return reg.MatchString(pwd)
}

// 校验昵称：只能输入1-20个以字母或者数字开头、可以含中文、下划线的字串。
func CheckNickName(name string) bool {
	reg, err := regexp.Compile("^([a-zA-Z0-9]|[\u4e00-\u9fa5]){1}([a-zA-Z0-9]|[_]|[\u4e00-\u9fa5]){0,20}$")
	if err != nil {
		fmt.Println(err)
	}
	return reg.MatchString(name)
}

//func CheackIMei(name string) bool {}

// 校验用户名：只能输入5-20个包含字母、数字或下划线的字串。
func CheckUserName(name string) bool {
	reg, err := regexp.Compile("^[a-zA-Z]{1}([a-zA-Z0-9]|[_]){4,19}$")
	if err != nil {
		fmt.Println(err)
	}
	return reg.MatchString(name)
}

func CheckId(id int) bool {
	if id == 0 {
		return false
	}
	reg, err := regexp.Compile("^[0-9]{1,10}$")
	if err != nil {
		fmt.Println(err)
	}
	return reg.MatchString(strconv.Itoa(id))
}
