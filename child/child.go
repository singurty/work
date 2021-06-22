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
	wg.Add(1)
	go pollRoot(conn, &wg)
	wg.Add(1)
	go pollWorkload(&wg)
	wg.Wait()
}

func pollRoot(conn net.Conn, wg *sync.WaitGroup) {
	for {
		buffer, _ := bufio.NewReader(conn).ReadBytes('\n')
		if len(buffer) == 0 {
			continue
		} else {
			fmt.Println(string(buffer)[1])
		}
	}
}

func pollWorkload(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		for _, work := range workload {
			if work.status == 0 {
				executeCommand(work.command)
			}
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

func executeCommand(command string) {
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
