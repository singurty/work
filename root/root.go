package root

import (
	"bufio"
	"net"
	"sync"
	"time"
	"log"
	"os"
	"strconv"
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
	handler chan string
}
var children []child
var workload []work

func Initialize(address string, port int, logFile string, wg *sync.WaitGroup) {
	f, err :=  os.OpenFile(logFile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println("test log entry")
	wg.Add(1)
	go listenForChildren(address, port, wg)
	log.Println("listening for children")
	wg.Add(1)
	go pollWorkload(wg)
	log.Println("polling workload")
	addWork(1, "whoami")
}

func addWork(merit int, command string) {
	newWork := work{
		merit: merit,
		status: 0,
		command: command,
	}
	workload = append(workload, newWork)
}

func pollWorkload(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		for index, work := range workload {
			if work.status == 0 && len(children) > 0  {
				wg.Add(1)
				handler := make(chan string)
				go handleWork(&workload[index], index, handler, wg)
				workload[index].status = 1
			}
		}
	}
}

/*
Work status codes (for root)
----------------------------
| ID  | meaning                              |
| --- | ------------------------------------ |
| 0   | this work doesn't have a handler     |
| 1   | there's a handler handling this work |
| 2   | work has been successfully executed  |
*/

func handleWork(work *work, index int, c chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	work.handler = c
	child := children[0]
	conn := child.conn
	message := "2" + strconv.Itoa(index) + work.command + "\n"
	err := sendMessage(conn, message)
	if err != nil {
		log.Println("failed to send work")
	} else {
		select {
		case message := <-c:
			if string(message[0]) == "4" {
				log.Print("work executed successfully")
				log.Print(message[1:])
			}
		}
	}
}

func sendMessage(conn net.Conn, message string) error {
	_, err := conn.Write([]byte(message))
	return err
}

func listenForChildren(address string, port int, wg *sync.WaitGroup) {
	ln, err := net.Listen("tcp", address + ":" + strconv.Itoa(port))
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	defer wg.Done()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleChild(conn)
	}
}

/*
Protocol IDs
-------------
| ID  | meaning                                                          |
| --- | ---------------------------------------------------------------- |
| 1   | Ping                                                             |
| 2   | sending work. followed by index and command ending with "\n"     |
| 3   | work added to workload. followed by work index                   |
| 4   | successfully did the work. index and output followd up to "\n"   |
| 5   | failed to do the work     
*/

func handleChild(conn net.Conn) {
	child := child{address: conn.RemoteAddr().String(), alive: true, conn: conn}
	children = append(children, child)
	log.Println("new child connected:", child.address)
	defer conn.Close()
	vitals := make(chan int)
	go checkChildVitals(&child, vitals)
	for {
		if !child.alive {
			log.Println("child dead:", child.address)
			return
		}
		buffer, _ := bufio.NewReader(conn).ReadBytes('\n')
		if len(buffer) == 0 {
			continue
		} else {
			switch string(buffer[0]) {
			case "1":
				vitals <- 1
			case "3":
			case "4":
				index, err := strconv.Atoi(string(buffer[1]))
				if err != nil {
					log.Println("invalid work id received")
					continue
				}
				work := workload[index]
				handler := work.handler
				handler <- "4" + string(buffer[2:len(buffer)-1]) 
			}
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
