package main

import (
	"github.com/singurty/fakework/root"
	"github.com/singurty/fakework/child"
	"os"
)

func main() {
	if os.Args[1] == "child" {
		child.Initialize("127.0.0.1", 8000)
	} else if os.Args[1] == "root" {
		root.Initialize("0.0.0.0", 8000)
	}
}
