package main

import (
	"strings"
	"fmt"
	"os"
	"sync"
	//"strconv"
	"github.com/c-bata/go-prompt"
	"github.com/singurty/fakework/child"
	"github.com/singurty/fakework/root"
	"gopkg.in/alecthomas/kingpin.v2"
)

var wg sync.WaitGroup
var (
	mode = kingpin.Arg("mode", "run either root or chld node").String()
)

func main() {
	defer wg.Wait()
	kingpin.Parse()
	switch *mode {
	case "root":
		rootShell()
	case "child":
		child.Initialize("127.0.0.1", 8001)
	}
}

func rootShell() {
	shell := prompt.New(
		executer,
		root.Completer,
		prompt.OptionPrefix(">> "),
		prompt.OptionTitle("root control center"),
	)
	shell.Run()
}

func executer(s string) {
	s = strings.TrimSpace(s)
	values := strings.Fields(s)
	switch strings.TrimSpace(values[0]) {
	case "quit":
	case "exit":
		fmt.Println("exiting..")
		os.Exit(0)
		return
	case "init":
		// commented for easier debugging
		//port, err := strconv.Atoi(values[2])
		//if err != nil {
		//	fmt.Println("invalid port number")
		//	return
		//}
		//root.Initialize(values[1], port, &wg)
		root.Initialize("0.0.0.0", 8001, &wg)
	default:
		fmt.Println("invalid command", s)
	}
}
