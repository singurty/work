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
	index int
	command string
}
type node struct {
	address string
	conn net.Conn
}

var root node
var workload []work

func Initialize(rootIp string, rootPort int) {
	var wg sync.WaitGroup
	root = node{address: rootIp + ":" + strconv.Itoa(rootPort)}
	conn, err  := net.Dial("tcp", root.address)
	root.conn = conn
	if err != nil {
		panic(err)
	}
	wg.Add(1)
	go pollRoot(&wg)
	wg.Add(1)
 	go pingRoot(&wg)
	wg.Add(1)
	go pollWorkload(&wg)
	wg.Wait()
}

func pollRoot(wg *sync.WaitGroup) {
	for {
		buffer, err := bufio.NewReader(root.conn).ReadBytes('\n')
		if err != nil {
			fmt.Println("error polling root")
			break
		}
		if len(buffer) == 0 {
			continue
		} else {
			if string(buffer[0]) == "2" {
				fmt.Println("work received")
				index, err := strconv.Atoi(string(buffer[1]))
				if err != nil {
					fmt.Println("invalid work index received")
					continue
				}
				addWork(index, string(buffer[2:len(buffer)-1]))
			}
		}
	}
	wg.Done()
}

func pollWorkload(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		for index, work := range workload {
			if work.status == 0 {
				handleWork(&workload[index])
			}
		}
	}
}

func handleWork(work *work) {
	// check if command is added yet. this prevents race condition in which an entry has been created but command has not been added.
	if len(work.command) == 0 {
		return
	}
	output, err := executeCommand(work.command)
	if err != nil {
		fmt.Println(err)
		work.status = 2
		sendMessage(root.conn, "5" + strconv.Itoa(work.index) + err.Error() + "\n")
		return
	}
	work.status = 1
	sendMessage(root.conn, "4" + strconv.Itoa(work.index) + output + "\n")
}

func addWork(index int, command string) {
	newWork := work{
		status: 0,
		command: command,
		index: index,
	}
	workload = append(workload, newWork)
	fmt.Println("work added", workload)
	sendMessage(root.conn, "3" + strconv.Itoa(newWork.index) + "\n")
	fmt.Println("sent work ack")
}

func executeCommand(command string) (string, error) {
	output, err := exec.Command(command).Output()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(output), nil
}

func sendMessage(conn net.Conn, message string) (error) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		return err
	}
	return nil
}

func pingRoot(wg *sync.WaitGroup) {
	for {
		err := sendMessage(root.conn, "1\n")
		if err != nil {
			fmt.Println(err)
			break
		}
		time.Sleep(10 * time.Second)
	}
	wg.Done()
}
