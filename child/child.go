package child

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

type work struct {
	status int
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
	wg.Add(1)
 	go pingRoot(&conn, &wg)
	wg.Wait()
}

func poolRoot(conn net.Conn, c chan string) {
	for {
		buffer, _ := bufio.NewReader(conn).ReadBytes('\n')
		if len(buffer) == 0 {
			continue
		}
	}
}

func addWork(command string) {
	newWork := work{
		status: 0,
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
	defer wg.Done()
	for {
		err := sendMessage(*conn, "1\n")
		if err != nil {
			fmt.Println(err)
			break
		}
		time.Sleep(5 * time.Second)
	}
}
