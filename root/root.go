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

func AddWork(merit int, command string) {
	client := dialDaemon()
	var resp rootd.Workload
	args := &rootd.AddWorkArgs{
		Merit: merit,
		Command: command,
	}
	err := client.Call("Workload.AddWork", args, &resp)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

func ShowWorkload() {
	client := dialDaemon()
	var resp rootd.Workload
	err := client.Call("Workload.GetWorkload", "", &resp)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

func ShowChildren() {
	client := dialDaemon()
	var resp rootd.Children
	err := client.Call("Children.GetChildren", "", &resp)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
