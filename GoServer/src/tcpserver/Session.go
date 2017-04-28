package tcpserver

import (
	"container/list"
	"fmt"
	"net"
	"sync"
)

type SessionHandler interface {
	OnConnect(session Session)
	OnClose(session Session)
	OnRecv(session Session, packet []byte)
}

type Session interface {
	Close()
	Send(packet []byte, packetSize int)
}

//-------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------
type session struct {
	conn      net.Conn
	server    *server
	sendQueue *list.List
	lock      *sync.Mutex
	cond      *sync.Cond
}

func createSession(conn net.Conn, server *server) *session {
	s := &session{conn: conn, server: server}
	s.sendQueue = list.New()
	s.lock = new(sync.Mutex)
	s.cond = sync.NewCond(s.lock)

	return s
}

func (session *session) Run() {
	go session.DoRecv()
	go session.DoSend()
}

func (session *session) DoRecv() {
	defer session.conn.Close()

	buffer := make([]byte, 4096)
	var totalRecvSize int

	for {
		recvBytes, err := session.conn.Read(buffer[totalRecvSize:])

		if err != nil {
			session.server.sessionHandler.OnClose(session)
			fmt.Println(err)
			return
		}

		if recvBytes == 0 {
			fmt.Println("recv bytes 0")
			session.server.sessionHandler.OnClose(session)
			return
		}

		totalRecvSize += recvBytes

		packetSize := session.server.protocolFilter.Parse(buffer, totalRecvSize)

		if packetSize > 0 {
			newBuffer := make([]byte, 4096)

			copy(newBuffer, buffer[packetSize:totalRecvSize])

			packet := buffer[:packetSize]
			session.server.sessionHandler.OnRecv(session, packet)

			buffer = newBuffer
			totalRecvSize = totalRecvSize - packetSize
		}
	}
}

var a int

func (session *session) DoSend() {
	for {
		session.lock.Lock()

		if session.sendQueue.Len() == 0 {
			session.cond.Wait()
		}

		e := session.sendQueue.Front()
		packet := e.Value
		session.sendQueue.Remove(e)

		session.lock.Unlock()

		session.conn.Write(packet.([]byte))
	}
}

func (session session) Close() {
	session.conn.Close()
}

func (session session) Send(packet []byte, packetSize int) {
	session.lock.Lock()

	session.sendQueue.PushBack(packet[:packetSize])

	session.cond.Signal()

	session.lock.Unlock()
}
