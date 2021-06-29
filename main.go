package main

import (
	"fmt"
	"strconv"
	"sync"

	//"strconv"
	"github.com/singurty/fakework/child"
	"github.com/singurty/fakework/root"
	"github.com/spf13/cobra"
)

var wg sync.WaitGroup

func main() {
	defer wg.Wait()
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
				root.Initialize("0.0.0.0", port, &wg)
			} else {
				port, err := strconv.Atoi(args[1])
				if err != nil {
					fmt.Println("invalid port number")
					return
				}
				root.Initialize(args[0], port, &wg)

			}
			fmt.Println("starting root control panel")
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
	var rootCmd = &cobra.Command{Use: "fakeroot"}
	rootCmd.AddCommand(cmdRoot, cmdChild)
	rootCmd.Execute()
}
