package root

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

type child struct {
	address string
	alive bool
	conn net.Conn
}
type work struct {
	merit int
	status int
	command string
}
var children []child
var workload []work

func Initialize(address string, port int) {
	var wg sync.WaitGroup
	wg.Add(1)
	go listenForChildren(address, port, &wg)
	wg.Add(1)
	go pollWorkload(&wg)
	AddWork(1, "neofetch")
	wg.Wait()
}

func AddWork(merit int, command string) {
	newWork := work{
		merit: merit,
		status: 0,
		command: command,
	}
	workload = append(workload, newWork)
}

func pollWorkload(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("polling workload")
	for {
		for index, work := range workload {
			if work.status == 0 {
				handleWork(&workload[index])
			}
		}
	}
}

func handleWork(work *work) {
	if len(children) == 0 {
		return
	}
	child := children[0]
	conn := child.conn
	message := "2" + work.command + "\n"
	fmt.Println("sending work to child")
	err := sendMessage(conn, message)
	if err != nil {
		fmt.Println("failed to send work")
	} else {
		fmt.Println("work sent to child")
		work.status = 1
	}
}

func sendMessage(conn net.Conn, message string) error {
	_, err := conn.Write([]byte(message))
	return err
}

func listenForChildren(address string, port int, wg *sync.WaitGroup) {
	fmt.Println("listening for children")
	ln, err := net.Listen("tcp", address + ":" + strconv.Itoa(port))
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	defer wg.Done()
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
	child := child{address: conn.RemoteAddr().String(), alive: true, conn: conn}
	children = append(children, child)
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
		case <-time.After(30 * time.Second):
			child.alive = false
			break
		}
	}
}
