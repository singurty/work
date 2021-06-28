package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	//"strconv"
	"github.com/jroimartin/gocui"
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
		rootControlPanel()
	case "child":
		child.Initialize("127.0.0.1", 8001)
	}
}

func rootControlPanel() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		fmt.Println("could not strat root control panel")
		return
	}
	defer g.Close()
	g.SetManagerFunc(rootLayout)
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		fmt.Println("could not start main loop")
		return
	}
}

func rootLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	logView, err := g.SetView("root control logs", maxX, maxY/2, maxX, maxY/2)
	if err != nil {
		return err
	}
	logView.Editable = false
	logView.Frame = true
	logView.Title = "logs"
	logView.Wrap = true
	logView.Autoscroll = true
	return nil
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
