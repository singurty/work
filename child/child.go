package child

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
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
	go pingRoot(conn)
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

func pingRoot( conn net.Conn) {
	for {
		err := sendMessage(conn, "1\n")
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(5* time.Second)
	}
}
