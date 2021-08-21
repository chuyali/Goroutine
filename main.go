/*实现两个goroutine通信，第三个goroutine可以设置各自频率，并查找各自发送次数，ctrl+c后，各个goroutine都可以退出*/
package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
)

//根据设置的频率发送数据，并监听是否ctrl+c,若有则退出，并关闭发送chan
func goroutine1(para *paraOption) {
	defer para.wg.Done()
	sendNum1 := 0
	for {
		select {
		case sig := <-para.sigChan:
			fmt.Println("goroutine1 | get", sig)
			close(para.msgSend)
			close(para.receNum1)
			fmt.Println("goroutine1 | close channel, exit")
			return
		default:
		}
		for i := 0; i < para.fre1; i++ {
			para.msgSend <- true
			sendNum1++
		}
		para.receNum1 <- sendNum1
		for i := 0; i < para.fre2; i++ {
			<-para.msgRecv
		}
	}
}

//根据设置的频率发送数据，并判断接收chan是否关闭，若关闭，则退出
func goroutine2(para *paraOption) {
	defer para.wg.Done()
	sendNum2 := 0
	for {
		for i := 0; i < para.fre1; i++ {
			_, ok := <-para.msgSend
			if !ok {
				fmt.Println("goroutine2 | exit")
				close(para.receNum2)
				return
			}
		}
		for i := 0; i < para.fre2; i++ {
			para.msgRecv <- true
			sendNum2++
		}
		para.receNum2 <- sendNum2
	}
}

//设置频率，并获取两个goroutine发送个数，若chan关闭，则退出
func goroutine3(para *paraOption) {
	para.fre1 = 3
	para.fre2 = 4
Loop:
	for {
		sendNum1, ok := <-para.receNum1
		if !ok {
			break Loop
		}
		fmt.Println("goroutine1 RecvNum", sendNum1)
		sendNum2, ok := <-para.receNum2
		if !ok {
			break Loop
		}
		fmt.Println("goroutine2 RecvNum", sendNum2)
	}
	fmt.Println("goroutine3 | exit")
	para.wg.Done()
}
func newPara() *paraOption {
	return &paraOption{
		sigChan:  make(chan os.Signal, 1),
		receNum1: make(chan int),
		receNum2: make(chan int),
		msgSend:  make(chan bool),
		msgRecv:  make(chan bool),
		wg:       sync.WaitGroup{},
		fre1:     1,
		fre2:     2,
	}
}

type paraOption struct {
	sigChan  chan os.Signal
	receNum1 chan int
	receNum2 chan int
	msgSend  chan bool
	msgRecv  chan bool
	wg       sync.WaitGroup
	fre1     int
	fre2     int
}

func main() {
	para := newPara()
	para.wg.Add(3)
	signal.Notify(para.sigChan, os.Interrupt)
	go goroutine3(para)
	go goroutine1(para)
	go goroutine2(para)
	para.wg.Wait()
}
