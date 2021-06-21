package slave

import (
	"os/exec"
	"fmt"
	"net"
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

func Initialize(masterIp string, masterPort int) {
	master := node {
		ip: masterIp,
		port: masterPort,
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

func sendMessage(dest node, message string) {
	conn, err := net.Dial("tcp", dest.ip + ':' + dest.port)
}
