package main

import (
	"github.com/c-bata/go-prompt"
	"github.com/singurty/fakework/child"
	"github.com/singurty/fakework/root"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	mode = kingpin.Arg("mode", "run either root or chld node").String()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()
	switch *mode {
	case "root":
		rootShell()
	case "child":
		child.Initialize("127.0.0.1", 8000)
	}
}

func rootShell() {
	prompt := prompt.New(root.Executor,
		root.Completer,
		prompt.OptionTitle("root control center"),
		prompt.OptionInputTextColor(prompt.Yellow),
	)
	prompt.Run()
}
