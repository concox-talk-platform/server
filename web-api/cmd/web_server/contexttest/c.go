/*
@Time : 2019/5/15 13:38 
@Author : yanKoo
@File : c
@Software: GoLand
@Description:
*/
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func worker2(ctx context.Context) {
	for {
		fmt.Printf("worker2\n")
		select {
		case <-ctx.Done():
			fmt.Printf("worker2 is done")
			return
		default:
			fmt.Printf("worker2 is alive")
		}
	}

}

func worker(ctx context.Context) {
	go worker2(ctx)
	//LOOP:
	//	for {
	fmt.Printf("worker\n")
	time.Sleep(time.Second * 3)
	//select {
	//case <- ctx.Done():
	//	fmt.Printf("worker will break")
	//	break LOOP
	//default:

	//}
	//}
	//wg.Done()
	return
}

func main() {
	//ctx := context.Background()
	//ctx, cancel := context.WithCancel(ctx)
	//wg.Add(1)
	//go worker(ctx)
	//time.Sleep(time.Second*3)
	//cancel()
	//wg.Wait()
	c := make(chan int, 1)
	go func() {
		time.Sleep(time.Second * 9)
		c <- 1
	}()
	for {
		select {
		case i := <-c:
			fmt.Printf("5555%d\n", i)
			return
		default:
			fmt.Printf("thisthisht\n")
		}

	}
}
