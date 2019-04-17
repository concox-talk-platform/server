/*
@Time : 2019/4/16 13:57 
@Author : yanKoo
@File : pushdata_test
@Software: GoLand
@Description:
*/
package server

import (
	"testing"
)

func TestPushDataExecutor(t *testing.T) {
	re := []int32{333}
	dc := &DataContent{DataChan:make(chan interface{}, 1)}
	e := make(chan interface{}, 1)

	dc.DataChan <- "test"
	dc.Receiver = re
	if err := pushDataExecutor(dc, e); err != nil {
		t.Logf("gen error: %v", err)
	}
}
