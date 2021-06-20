package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	priority, err := strconv.Atoi(os.Args[1])
	if priority > 100 || priority < 1 {
		panic("priority must be between 1 and 100")
	}
	command := os.Args[2]
	fmt.Println(priority)
	output, err := exec.Command(command).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Println(string(output[:]))
}
