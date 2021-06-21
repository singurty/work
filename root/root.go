package root

import (
	"fmt"
	"net"
)

func Initialize() {

}

func listenForChildren(address string, port int) {
	fmt.Println("listening for children")
	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go handleChildren(conn)
	}
}

func handleChildren(conn net.Conn) {
	fmt.Println("new child connected:", conn.RemoteAddr().String())
}
