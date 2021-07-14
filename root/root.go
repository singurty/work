package root

import (
	"fmt"
	"net/rpc"
	"github.com/hpcloud/tail"
	"github.com/singurty/fakework/rootd"
)

func ViewLog(logFile string, follow bool) {
	t, err := tail.TailFile(logFile, tail.Config{Follow: follow})
	if err != nil {
		panic(err)
	}
	for line := range t.Lines {
		fmt.Println(line.Text)
	}
}

func dialDaemon() *rpc.Client {
	client, err := rpc.Dial("tcp", "127.0.0.1:9002")
	if err != nil {
		fmt.Println("could not connect to the root daemon. is it running?")
		fmt.Println()
		panic(err)
	}
	return client
}

func AddWork(command string, each bool) {
	client := dialDaemon()
	var resp rootd.Workload
	args := &rootd.AddWorkArgs{
		Command: command,
		Each: each,
	}
	err := client.Call("Workload.AddWork", args, &resp)
	if err != nil {
		panic(err)
	}
	fmt.Println("successfully added. new workload:")
	ShowWorkload()
}

func ShowWorkload() {
	client := dialDaemon()
	var resp rootd.Workload
	err := client.Call("Workload.GetWorkload", "", &resp)
	if err != nil {
		panic(err)
	}
	for index, work := range resp {
		var status string
		switch work.Status {
		case 0:
			status = "new work. nothing done yet"
		case 1:
			status = "work sent to a child"
		case 2:
			status = "work successfully executed"
		case 4:
			status = "work was sent to a child but error'd"
		}
		if len(work.Output) == 0 {
			fmt.Printf("%v. Command: %v Status: %v\n", index + 1, work.Command, status)
		} else {
			fmt.Printf("%v. Command: %v Status: %v Output: %v\n", index + 1, work.Command, status, work.Output)
		}
	}
}

func ShowChildren() {
	client := dialDaemon()
	var resp []rootd.Child
	err := client.Call("Children.GetChildren", "", &resp)
	if err != nil {
		panic(err)
	}
	for index, child := range resp {
		fmt.Printf("%v. Address: %v Alive: %t\n", index + 1, child.Address, child.Alive)
	}
}
