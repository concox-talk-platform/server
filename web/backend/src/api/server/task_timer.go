/*
@Time : 2019/4/16 14:48 
@Author : yanKoo
@File : task_timer
@Software: GoLand
@Description:
*/
package server

//
//import (
//	"log"
//	"sync"
//	"time"
//	//"log"
//)
//
//type Worker struct {
//	ticker *time.Timer
//	uid    int32
//	// 添加一个修改定时器时间的方法
//}
//
//var timerMap sync.Map
//
//func newWorker(interval time.Duration, uid int32) *Worker {
//	return &Worker{
//		ticker: time.NewTimer(interval * time.Second),
//		uid:    uid,
//	}
//}
//
//func (w *Worker) startWorker() {
//	//for {
//	select {
//	case <-w.ticker.C:
//		w.deleteStreamMap(w.uid)
//	}
//	//}
//}
//
//func StartTimer(interval time.Duration, uid int32) {
//	// Start map cleaning
//	w := newWorker(interval, uid)
//	go w.startWorker()
//}
//
//func (w *Worker) deleteStreamMap(uid int32) {
//	log.Println("now is: ", time.Now().String(), " will delete: ", uid)
//	StreamMap.Delete(uid)
//}

//func ResetTime(interval time.Duration, uid int32) {
//	if v, ok := timerMap.Load(uid); !ok {
//		log.Println("no this timer")
//	} else {
//		v.(Worker).ticker.Reset(interval)
//	}
//}
//
//func RunTimer(interval int) {
//	//创建定时器并设置定时时间
//	TimerDemo := time.NewTimer(time.Duration(interval) * time.Second)
//
//	//循环监听定时器
//	for {
//		select {
//		case <-TimerDemo.C:
//			clearStreamMap()
//		}
//	}
//}
//
//func clearStreamMap() {
//	StreamMap.Range(func(key, value interface{}) bool {
//		StreamMap.Delete(key)
//		return true
//	})
//	log.Println("this is timer, map : % +v", StreamMap)
//}

