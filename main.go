package main

import (
	"fmt"
	"sync"
	"time"
)

var flag = 0

func main() {
	mu := sync.Mutex{}
	go printx(&mu)
	go printy(&mu)
	time.Sleep(time.Second * 100)
}

func printx(m *sync.Mutex) {
	for i := 0; i < 10000000; i++ {
		m.Lock()
		//fmt.Println("x lock")
		if flag == 0 {
			fmt.Println("x")
			flag = 1
		}
		//fmt.Println("x unlock")
		m.Unlock()
	}
}

func printy(m *sync.Mutex) {
	for i := 0; i < 10000000; i++ {
		m.Lock()
		//fmt.Println("y lock")
		if flag == 1 {
			fmt.Println("y")
			flag = 0
		}
		//fmt.Println("y unlock")
		m.Unlock()
	}
}
