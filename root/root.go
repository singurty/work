package root

import (
	"fmt"
	"github.com/hpcloud/tail"
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
