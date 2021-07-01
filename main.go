package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/singurty/fakework/child"
	"github.com/singurty/fakework/root"
	"github.com/singurty/fakework/rootd"
	"github.com/spf13/cobra"
)

var wg sync.WaitGroup

func main() {
	defer wg.Wait()
	var logFile *os.File
	defer logFile.Close()
	var logFileName string
	var follow bool
	var cmdRoot = &cobra.Command{
		Use: "root",
		Short: "run a root node",
		Long: "Run a root node listening on specified IP address and port. If IP address is not specified it'll listen on all interfaces",
		Args: cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				port, err := strconv.Atoi(args[0])
				if err != nil {
					fmt.Println("invalid port number")
					return
				}
				logFile, err :=  os.OpenFile(logFileName, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
				if err != nil {
					fmt.Println("error opening file:", err)
					return
				}
				rootd.Initialize("0.0.0.0", port, logFile, &wg)
			} else {
				port, err := strconv.Atoi(args[1])
				if err != nil {
					fmt.Println("invalid port number")
					return
				}
				logFile, err :=  os.OpenFile(logFileName, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
				if err != nil {
					fmt.Println("error opening file:", err)
					return
				}
				rootd.Initialize(args[0], port, logFile, &wg)
			}
		},
	}
	var cmdChild = &cobra.Command{
		Use: "child",
		Short: "run a child node",
		Long: "Run a child node and connect to a root node listening on specified IP address and port",
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			port, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("invalid port number")
				return
			}
			child.Initialize(args[0], port)
		},
	}
	var cmdLog = &cobra.Command{
		Use: "log",
		Short: "view root daemon logs",
		Long: "Read logs produced by the root daemon. Reads from 'root.log' file by default",
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			root.ViewLog(logFileName, follow)
		},
	}
	var cmdAdd = &cobra.Command{
		Use: "add",
		Short: "add work",
		Long: "add a command to be executed by child",
		Args: cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				root.AddWork(0, args[0])
			} else if len(args) == 2 {
				merit, err := strconv.Atoi(args[1])
				if err != nil {
					fmt.Println("invalid merit")
					return
				}
				root.AddWork(merit, args[0])
			}
		},
	}
	var cmdShow = &cobra.Command{
		Use: "show",
		Short: "show something",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "workload":
				root.ShowWorkload()
			}
		},
	}
	var rootCmd = &cobra.Command{Use: "fakeroot"}

	cmdRoot.Flags().StringVarP(&logFileName, "log", "l", "root.log", "file to write logs to (Default: root.log)")
	cmdLog.Flags().BoolVarP(&follow, "follow", "f", false, "keep polling for logs")
	cmdLog.Flags().StringVarP(&logFileName, "log", "l", "root.log", "file to write logs to (Default: root.log)")
	rootCmd.AddCommand(cmdRoot, cmdChild, cmdLog, cmdAdd, cmdShow)
	rootCmd.Execute()
}
