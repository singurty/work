package main

import (
	"os"
	"github.com/abiosoft/ishell/v2"
	"github.com/singurty/fakework/child"
	"github.com/singurty/fakework/root"
)

func main() {
	if os.Args[1] == "child" {
		child.Initialize("127.0.0.1", 8000)
	} else if os.Args[1] == "root" {
		rootShell()
		//root.Initialize("0.0.0.0", 8000)
	}
}

func rootShell() {
	shell := ishell.New()
	shell.AddCmd(&ishell.Cmd{
		Name: "init",
		Help: "Initialize root node",
		Func: func(c *ishell.Context) {
			c.Println("Initializing root node")
			root.Initialize("0.0.0.0", 8000)
		},
	})
	shell.Run()
}
