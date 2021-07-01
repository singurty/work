package rootd

import (
	"bufio"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"sync"
	"time"
)

type child struct {
	address string
	alive bool
	conn net.Conn
}
type Work struct {
	Merit int
	Status int
	Command string
	handler chan string
}
type AddWorkArgs struct {
	Merit int
	Command string
}
type Workload []Work
var workload Workload
var children []child

func Initialize(address string, port int, f *os.File, wg *sync.WaitGroup) {
	log.SetOutput(f)
	wg.Add(1)
	log.Println("starting rpc server")
	go startRpc()
	wg.Add(1)
	go listenForChildren(address, port, wg)
	log.Println("listening for children")
	wg.Add(1)
	go pollWorkload(wg)
	log.Println("polling workload")
}

func (w *Workload) AddWork(args AddWorkArgs, resp *Workload) error {
	newWork := Work{
		Merit: args.Merit,
		Status: 0,
		Command: args.Command,
	}
	workload = append(workload, newWork)
	log.Println("work added:", newWork)
	*resp = workload
	return nil
}

func (w *Workload) GetWorkload(args string,  resp *Workload) error {
	*resp = workload
	return nil
}

func startRpc() {
	server := rpc.NewServer()
	server.Register(&workload)
	listener, err := net.Listen("tcp", "127.0.0.1:9002")
	if err != nil {
		panic(err)
	}
	for {
		conn, _ := listener.Accept()
		go server.ServeConn(conn)
	}
}

func pollWorkload(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		for index, work := range workload {
			if work.Status == 0 && len(children) > 0  {
				wg.Add(1)
				handler := make(chan string)
				go handleWork(&workload[index], index, handler, wg)
				workload[index].Status = 1
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

func handleWork(work *Work, index int, c chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	work.handler = c
	child := children[0]
	conn := child.conn
	message := "2" + strconv.Itoa(index) + work.Command + "\n"
	err := sendMessage(conn, message)
	if err != nil {
		log.Println("failed to send work")
		work.Status = 0
		return
	}
	for {
		select {
		case message := <-c:
			switch string(message[0]) {
			case "3":
				log.Println("child ack'd the work")
			case "4":
				work.Status = 2
				log.Println("work executed successfully")
				log.Println(message[1:])
			case "5":
				work.Status = 3
				log.Println("child failed to do the work:", work.Command)
				log.Println(message[1:])
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
		}
		switch string(buffer[0]) {
		case "1":
			vitals <- 1
		case "3", "4", "5":
			index, err := strconv.Atoi(string(buffer[1]))
			if err != nil {
				log.Println("invalid work index received")
				continue
			}
			work := workload[index]
			handler := work.handler
			handler <- string(buffer[0]) + string(buffer[2:len(buffer)-1]) 
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
