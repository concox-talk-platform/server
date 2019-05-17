/*
@Time : 2019/5/17 10:55 
@Author : yanKoo
@File : validation
@Software: GoLand
@Description:
*/
package model

import (
	"gopkg.in/go-playground/validator.v8"
	"reflect"
	"server/web-api/utils"
)

func DeviceInfoAble(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	if d, ok := field.Interface().([]*Device); ok {
		// 校验device的imei号码
		for _, v := range d {
			if !utils.CheckIMei(v.IMei) {
				return false
			}
			return true
		}
	}
	return false
}
