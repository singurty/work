package child

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

type work struct {
	merit int
	command string
}
type node struct {
	ip string
	port int
}

var workload []work

func Initialize(rootIp string, rootPort int) {
	var wg sync.WaitGroup
	conn, err  := net.Dial("tcp", rootIp + ":" + strconv.Itoa(rootPort))
	if err != nil {
		panic(err)
	}
//	response, err := sendMessage(conn, "init\n")
//	if err != nil {
//		panic(err)
//	}
//	if response == "ack" {
//		fmt.Println("connected to root")
//	}
	wg.Add(1)
 	go pingRoot(&conn, &wg)
	wg.Wait()
}

func AddWork(merit int, command string) {
	newWork := work{
		merit: merit,
		command: command,
	}
	workload = append(workload, newWork)
}

func executeCommand(command string, priority int) {
	output, err := exec.Command(command).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Print(string(output[:]))
}

func sendMessage(conn net.Conn, message string) (error) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		return err
	}
	return nil
}

func pingRoot(conn *net.Conn, wg *sync.WaitGroup) {
	for {
		err := sendMessage(*conn, "1\n")
		if err != nil {
			fmt.Println(err)
			break
		}
		time.Sleep(5 * time.Second)
	}
	wg.Done()
}
