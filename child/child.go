package child

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
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
	master := node {
		ip: rootIp,
		port: rootPort,
	}
	response, err := sendMessage(master, "init")
	if err != nil {
		panic(err)
	}
	if response == "ack" {
		fmt.Println("connected to root")
	}
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

func sendMessage(dest node, message string) (string, error) {
	conn, err  := net.Dial("tcp", dest.ip + ":" + strconv.Itoa(dest.port))
	if err != nil {
		return "", err
	}
	defer conn.Close()
	conn.Write([]byte(message))
	response := make([]byte, 4096)
	_, err = conn.Read(response)
	if err != nil {
		return "", err
	}
	return string(response), nil
}
