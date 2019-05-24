/*
@Time : 2019/5/11 14:40
@Author : yanKoo
@File : mian.go
@Software: GoLand
@Description:
*/
package main

import (
	"fmt"
)

type interphoneMsg struct {
	Uid       string `json:"uid"`
	MsgType   string `json:"m_type"`
	Md5       string `json:"md5"`
	GId       string `json:"grp_id"`
	FilePath  string `json:"file_path"`
	Timestamp string `json:"timestamp"`
}

type node struct {
	value int
	next *node
}

func mai9n()  {
	head := &node{value:0}
	node1 := &node{value:1}
	node2 := &node{value:2}
	node3 := &node{value:3}
	node4 := &node{value:4}
	node5 := &node{value:5}
	head.next = node1
	node1.next = node2
	node2.next = node3
	node3.next = node4
	node4.next = node5

	printNode(head)
	printNode(rever(head))

}


func printNode(head *node) {
	for head != nil {
		//fmt.Print(head.value, "\t")
		fmt.Println(head)
		head = head.next
	}
	fmt.Println()
}


func rever(head *node) *node {
	var  (
		now =  &node{}
		nextN = &node{}
	)
	for head != nil {
		// 保存头节点的下一个节点，
		nextN = head.next

		head.next = now
		now = head
		head = nextN
	}

	return now
}




//
//func reverseNode(head *node) *node {
//	//  先声明两个变量
//	//  前一个节点
//	var preNode *node
//	preNode = nil
//	//  后一个节点
//	nextNode := new(node)
//	nextNode = nil
//	for head != nil {
//		//  保存头节点的下一个节点，
//		nextNode = head.nextNode
//		//  将头节点指向前一个节点
//		head.nextNode = preNode
//		//  更新前一个节点
//		preNode = head
//		//  更新头节点
//		head = nextNode
//	}
//	return preNode
//}


func DeepCopy(value interface{}) interface{} {
	if valueMap, ok := value.(map[interface{}]interface{}); ok {
		newMap := make(map[interface{}]interface{})
		for k, v := range valueMap {
			newMap[k] = DeepCopy(v)
		}
		return newMap
	}
	return value
}


//func main() {
//	r := cache.GetRedisClient()
//	for i := 0; i < 3; i++ {
//		//go func(i int) {
//		m := &interphoneMsg{
//			Uid:       "8", //strconv.FormatInt(int64(i),10),
//			MsgType:   "ptt",
//			Md5:       "555555555555",
//			GId:       "6",
//			FilePath:  "hahaahahahahha.ptt",
//			Timestamp: time.Now().Format("2006-01-02 15:04:05"),
//		}
//		v, e := json.Marshal(m)
//		fmt.Printf("%s", v)
//		if e != nil {
//
//		}
//		_, err := r.Do("lpush", "janusImUrlList", v)
//		if err != nil {
//			fmt.Printf("push redis data with error: %+v", err)
//			return
//		}
//		//}(i)
//	}
//	//select {}
//}

type person struct {
	name string
	age int
	wt string
}
