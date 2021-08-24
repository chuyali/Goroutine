/*实现两个goroutine通信，main可以设置各自频率，并查找各自发送次数，ctrl+c后，各个goroutine都可以退出*/
package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"sync"
	"time"
)

//根据设置的频率发送数据，并监听Ctrl+c，若输入Ctrl+c则退出
func goroutine(wg *sync.WaitGroup, sendNum *int, msgRecv, msgSend chan bool, fre int) {
	defer wg.Done()
	var sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	t := time.Tick(time.Duration(fre) * time.Second)
	for {
		select {
		case sig := <-sigChan:
			fmt.Println("goroutine | get, exit", sig)
			fmt.Println("goroutine exit")
			return
		case <-t:
			msgRecv <- true
			*sendNum++
		case <-msgSend:

		default:
		}
	}
}

func getSetPara(para *paraOption, sendNumA, sendNumB *int) {
	var sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	fmt.Println(`please enter your demand:
		set-ping-pong A 3s : set goroutine frequency A 3s and B default
		set-ping-pong B 3s : set goroutine frequency B 3s and A default
		set-ping-pong 3s : set goroutine A frequency 3s and B 3s
		get-ping-pong:get sendNum and frequency of goroutine A and goroutine B
		`)
	for {
		select {
		case sig := <-sigChan:
			fmt.Println("for | get", sig)
			fmt.Println("for, exit")
			return
		default:
		}
		fmt.Scanf("%s %s %s\n", &para.pingPongOpt, &para.setGoroutine, &para.setFre)
		re := regexp.MustCompile("[0-9]+")
		freNum, _ := strconv.Atoi(re.FindString(para.setFre))
		switch {
		case para.pingPongOpt == "set-ping-pong" && para.setGoroutine == "A" && freNum != 0:
			para.freA = freNum
		case para.pingPongOpt == "set-ping-pong" && para.setGoroutine == "B" && freNum != 0:
			para.freB = freNum
		case para.pingPongOpt == "set-ping-pong" && freNum != 0:
			para.freA = freNum
			para.freB = freNum
		case para.pingPongOpt == "get-ping-pong":
			fmt.Printf("goroutineA sendNum %d sendFre %s,goroutineB sendNum %d sendFre %s\n",
				*sendNumA, strconv.Itoa(para.freA)+"s", *sendNumB, strconv.Itoa(para.freB)+"s")
		}
		fmt.Printf("frequency A : %s B : %s:\n",
			strconv.Itoa(para.freA)+"s", strconv.Itoa(para.freB)+"s")
	}
}
func newPara() *paraOption {
	return &paraOption{
		freA: 1, //default frequency
		freB: 2, //default frequency
	}
}

type paraOption struct {
	freA         int
	freB         int
	pingPongOpt  string
	setGoroutine string
	setFre       string
}

func main() {
	var sendNumA int
	var sendNumB int
	para := newPara()
	var wg = sync.WaitGroup{}
	wg.Add(2)

	var msgSend = make(chan bool, 1)
	var msgRecv = make(chan bool, 1)
	go goroutine(&wg, &sendNumA, msgSend, msgRecv, para.freA)
	go goroutine(&wg, &sendNumB, msgRecv, msgSend, para.freB)

	getSetPara(para, &sendNumA, &sendNumB)
	wg.Wait()
}
