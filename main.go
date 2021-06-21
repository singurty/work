package main

import (
	"fmt"
	"os"
	"fakework/master"
	"fakework/slave"
)

func main() {
	if os.Args[1] == "slave" {
		slave.Initialize()
	} else if os.Args[1] == "master" {
		master.Initialize()
	}
}
