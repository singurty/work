package main

import (
	"github.com/singurty/fakework/root"
	"github.com/singurty/fakework/child"
)

func main() {
	if os.Args[1] == "child" {
		child.Initialize()
	} else if os.Args[1] == "root" {
		root.Initialize()
	}
}
