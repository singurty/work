package root

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"time"
)

type child struct {
	address string
	alive bool
}

func Initialize(address string, port int) {
	listenForChildren(address, port)
}

func listenForChildren(address string, port int) {
	fmt.Println("listening for children")
	ln, err := net.Listen("tcp", address + ":" + strconv.Itoa(port))
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleChild(conn)
	}
}

func handleChild(conn net.Conn) {
	child := child{address: conn.RemoteAddr().String(), alive: true}
	fmt.Println("new child connected:", child.address)
	defer conn.Close()
	vitals := make(chan int)
	go checkChildVitals(&child, vitals)
	for {
		if !child.alive {
			fmt.Println("child dead:", child.address)
			return
		}
		buffer, _ := bufio.NewReader(conn).ReadBytes('\n')
		if len(buffer) == 0 {
			continue
		} else if string(buffer) == "1\n" {
			vitals <- 1
		}
	}
}

func checkChildVitals(child *child, c chan int) {
	for {
		select {
		case <-c:
			continue
		case <-time.After(10 * time.Second):
			child.alive = false
			break
		}
	}
}
