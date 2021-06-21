package main

import (
	"fmt"
	"os"
	"fakework/root"
	"fakework/chld"
)

func main() {
	if os.Args[1] == "child" {
		child.Initialize()
	} else if os.Args[1] == "root" {
		root.Initialize()
	}
}
