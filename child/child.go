package child

import (
	"os/exec"
	"fmt"
	"net"
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
	sendMessage(master, "init")
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

func sendMessage(dest node, message string) (response string, err error) {
	conn, err  := net.Dial("tcp", dest.ip + ":" + strconv.Itoa(dest.port))
	conn.Write([]byte(message))
}
