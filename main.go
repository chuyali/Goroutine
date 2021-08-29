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
func goroutine(wg *sync.WaitGroup, sendNumList []chan int, require chan int, msgRecv, msgSend chan bool, fre int) {
	defer wg.Done()
	var sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	var num int
	t := time.Tick(time.Duration(fre) * time.Second)
	for {
		select {
		case sig := <-sigChan:
			fmt.Println("goroutine | get, exit", sig)
			fmt.Println("goroutine exit")
			return
		case <-t:
			msgRecv <- true
			num++
		case <-msgSend:
		case iD := <-require:
			fmt.Println("goroutine iD", iD)
			sendNumList[iD] <- num
		default:
		}
	}
}

func getSetPara(para *paraOption, requireA, requireB chan int, sendNumA, sendNumB chan int, iD int) {
	var sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	fmt.Println(`please enter your demand:
		set-ping-pong A 3s : set goroutine frequency A 3s and B default
		set-ping-pong B 3s : set goroutine frequency B 3s and A default
		set-ping-pong 3s : set goroutine A frequency 3s and B 3s
		get-ping-pongA :get sendNum and frequency of goroutine A 
		get-ping-pongB :get sendNum and frequency of goroutine B
		`)
	for {
		fmt.Scanf("%s %s %s\n", &para.pingPongOpt, &para.setGoroutine, &para.setFre)
		re := regexp.MustCompile("[0-9]+")
		freNum, _ := strconv.Atoi(re.FindString(para.setFre))
		switch {
		case para.pingPongOpt == "set" && para.setGoroutine == "A":
			para.freA = freNum
		case para.pingPongOpt == "set" && para.setGoroutine == "B":
			para.freB = freNum
		case para.pingPongOpt == "set" && freNum != 0:
			para.freA = freNum
			para.freB = freNum
		case para.pingPongOpt == "getA":
			requireA <- iD
			getSendNum(para.freA, sendNumA, iD)
		case para.pingPongOpt == "getB":
			requireB <- iD
			getSendNum(para.freB, sendNumB, iD)
		default:
		}
		fmt.Printf("frequency A : %s B : %s:\n",
			strconv.Itoa(para.freA)+"s", strconv.Itoa(para.freB)+"s")
		select {
		case sig := <-sigChan:
			fmt.Println("for | get", sig)
			fmt.Println("for, exit")
			return
		default:
		}
	}
}
func getSendNum(fre int, sendNum chan int, iD int) {
	select {
	case num := <-sendNum:
		fmt.Printf("iD : %d ,goroutineA sendNum %d sendFre %s\n",
			iD, num, strconv.Itoa(fre)+"s")
	case <-time.After(1 * time.Second):
		fmt.Println("超时")
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
	var MAX_CLIENT = 4
	para := newPara()
	var wg = sync.WaitGroup{}
	wg.Add(2)

	var msgSend = make(chan bool, 1)
	var msgRecv = make(chan bool, 1)
	sendNumAList := []chan int{}
	sendNumBList := []chan int{}
	for i := 0; i < MAX_CLIENT; i++ {
		tmp := make(chan int, 1)
		sendNumAList = append(sendNumAList, tmp)
		sendNumBList = append(sendNumBList, tmp)
	}
	var requireA = make(chan int, 1)
	var requireB = make(chan int, 1)
	go goroutine(&wg, sendNumAList, requireA, msgSend, msgRecv, para.freA)
	go goroutine(&wg, sendNumBList, requireB, msgRecv, msgSend, para.freB)
	for i := 0; i < MAX_CLIENT; i++ {
		go getSetPara(para, requireA, requireB, sendNumAList[i], sendNumBList[i], i)
	}

	wg.Wait()
}
