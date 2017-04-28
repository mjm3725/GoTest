package main // GoTest project main.go

import (
	"fmt"
	"runtime"
	"tcpserver"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println("Go server start~~~~~~~~~~~")

	server := tcpserver.CreateServer(MyHandler{}, tcpserver.NullTerminatedProtocolFilter{})
	server.Start("8080")

	for {
		time.Sleep(time.Second * 2)
	}
}

type MyHandler struct {
}

func (handler MyHandler) OnConnect(session tcpserver.Session) {
	fmt.Println("OnConnect")
}

func (handler MyHandler) OnClose(session tcpserver.Session) {
	fmt.Println("OnClose")
}

func (handler MyHandler) OnRecv(session tcpserver.Session, packet []byte) {
	//fmt.Println("OnRecv: " + string(packet))
	session.Send(packet, len(packet))
}
