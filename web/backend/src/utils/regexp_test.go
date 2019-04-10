/*
@Time : 2019/4/1 13:40 
@Author : yanKoo
@File : regexp_test
@Software: GoLand
@Description:
*/
package utils

import "testing"

func testCheckPwd(t *testing.T) {
	t.Log(CheckPwd("gagdfh"))
}

func TestCheckNickName(t *testing.T) {
	t.Log(CheckNickName("中"))
}

func testCheckUserName(t *testing.T) {
	t.Log(CheckUserName("safs"))
}