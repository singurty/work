package child

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"sync"
	"time"
	"github.com/google/shlex"
	"github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
	multiaddr "github.com/multiformats/go-multiaddr"
)

type work struct {
	status int
	index int
	command string
}
// node is used for both peers and root
type node struct {
	address string
	accepting bool
	conn net.Conn
}

var wg sync.WaitGroup
var root node
var peers []node
var workload []work

func Initialize(rootIp string, rootPort int) {
	root = node{address: rootIp + ":" + strconv.Itoa(rootPort)}
	conn, err  := net.Dial("tcp", root.address)
	root.conn = conn
	if err != nil {
		panic(err)
	}
	wg.Add(1)
	go pollRoot(&wg)
	wg.Add(1)
	go pingRoot(&wg)
	wg.Add(1)
	go pollWorkload(&wg)
	wg.Wait()
}

func StartPeerNode() {
	wg.Add(1)
	ctx := context.Background()
	peerNode, err := libp2p.New(ctx)
	if err != nil {
		panic(err)
	}
	peerInfo := peerstore.AddrInfo{
		ID: peerNode.ID(),
		Addrs: peerNode.Addrs(),
	}
	addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("node address:", addrs[0])
	wg.Wait()
	if err := peerNode.Close(); err != nil {
		panic(err)
	}
	wg.Done()
}

func ConnectToPeer(address string) {
	fmt.Println("connecting to", address)
	addr, err := multiaddr.NewMultiaddr(address)
	if err != nil {
		panic(err)
	}
	peer, err := peerstore.AddrInfoFromP2pAddr(addr)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	peerNode, err := libp2p.New(ctx)
	if err != nil {
		panic(err)
	}
	if err := peerNode.Connect(ctx, *peer); err != nil {
		panic(err)
	}
}

func pollRoot(wg *sync.WaitGroup) {
	for {
		buffer, err := bufio.NewReader(root.conn).ReadBytes('\n')
		if err != nil {
			fmt.Println("error polling root")
			break
		}
		if len(buffer) == 0 {
			continue
		} else {
			if string(buffer[0]) == "2" {
				index, err := strconv.Atoi(string(buffer[1]))
				if err != nil {
					fmt.Println("invalid work index received")
					continue
				}
				addWork(index, string(buffer[2:len(buffer)-1]))
			}
		}
	}
	wg.Done()
}

func pollWorkload(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		for index, work := range workload {
			if work.status == 0 {
				handleWork(&workload[index])
			}
		}
	}
}

func handleWork(work *work) {
	// check if command is added yet. this prevents race condition in which an entry has been created but command has not been added.
	if len(work.command) == 0 {
		return
	}
	command, err := shlex.Split(work.command)
	output, err := exec.Command(command[0], command[1:]...).Output()
	if err != nil {
		fmt.Println(err)
		work.status = 2
		sendMessage(root.conn, "5" + strconv.Itoa(work.index) + err.Error() + "\n")
		return
	}
	work.status = 1
	sendMessage(root.conn, "4" + strconv.Itoa(work.index) + string(output) + "\n")
}

func addWork(index int, command string) {
	newWork := work{
		status: 0,
		command: command,
		index: index,
	}
	workload = append(workload, newWork)
	sendMessage(root.conn, "3" + strconv.Itoa(newWork.index) + "\n")
}

func sendMessage(conn net.Conn, message string) (error) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		return err
	}
	return nil
}

func requestPeersList() {
	sendMessage(root.conn, "6\n")
}

func pingRoot(wg *sync.WaitGroup) {
	for {
		err := sendMessage(root.conn, "1\n")
		if err != nil {
			fmt.Println(err)
			break
		}
		time.Sleep(10 * time.Second)
	}
	wg.Done()
}
