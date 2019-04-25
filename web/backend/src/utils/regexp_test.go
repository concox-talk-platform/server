/*
@Time : 2019/4/1 13:40 
@Author : yanKoo
@File : regexp_test
@Software: GoLand
@Description:
*/
package utils

import (
	"io/ioutil"
	"os"
	"testing"
)

func testCheckPwd(t *testing.T) {
	t.Log(CheckPwd("gagdfh"))
}

func testCheckNickName(t *testing.T) {
	t.Log(CheckNickName("中"))
}

func testCheckUserName(t *testing.T) {
	t.Log(CheckUserName("safs"))
}

func testCheckId(t *testing.T) {
	t.Log(CheckId(-1))
}

func TestGetFileType(t *testing.T) {
	//f, err := os.Open("C:\\Users\\Administrator\\Desktop\\api.html")
	//f, err := os.Open("C:\\Users\\Administrator\\Desktop\\123.mp4")
	f, err := os.Open("C:\\Users\\Public\\Music\\Sample Music\\Kalimba.mp3")
	if err != nil {
		t.Logf("open error: %v", err)
	}

	fSrc, err := ioutil.ReadAll(f)
	t.Log(GetFileType(fSrc[:10]))
}