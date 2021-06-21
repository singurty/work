package root

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

func Initialize(address string, port int) {
	listenForChildren(address, port)
}

func listenForChildren(address string, port int) {
	fmt.Println("listening for children")
	ln, err := net.Listen("tcp", address + ":" + strconv.Itoa(port))
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleChild(conn)
	}
}

func handleChild(conn net.Conn) {
	fmt.Println("new child connected:", conn.RemoteAddr().String())
	for {
		buffer, _ := bufio.NewReader(conn).ReadBytes('\n')
		if len(buffer) > 0 {
			fmt.Println(string(buffer))
		}
	}
}
