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
	"log"
	"regexp"
)

// 校验手机号码
func CheckPhone (phone string) bool {
	reg, err := regexp.Compile("^((13[0-9])|(15[^4,\\D])|(18[0,5-9]))\\d{8}$")
	if err != nil {
		log.Println(err)
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

	reg, err := regexp.Compile("^(\\w){6,20}$")
	if err != nil {
		fmt.Println(err)
	}
	return reg.MatchString(pwd)
}

// 校验登录名：只能输入5-20个以字母开头、可带数字、“_”的字串
func CheckName(name string) bool {
	reg, err := regexp.Compile("^[a-zA-Z]{1}([a-zA-Z0-9]|[_]){4,19}$")
	if err != nil {
		fmt.Println(err)
	}
	return reg.MatchString(name)
}

func maint() {
	log.Println(CheckPhone("13769569111"))
}