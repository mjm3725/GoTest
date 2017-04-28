package tcpserver

import (
	"fmt"
	"net"
)

type Server interface {
	Start(string) bool
}

func CreateServer(sessionHandler SessionHandler, protocolFilter ProtocolFilter) Server {
	return server{sessionHandler: sessionHandler, protocolFilter: protocolFilter}
}

//-------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------

type server struct {
	sessionHandler SessionHandler
	protocolFilter ProtocolFilter
}

func (server server) Start(port string) bool {
	listener, err := net.Listen("tcp", ":"+port)

	if err != nil {
		fmt.Println(err)
		return false
	}

	go func() {
		defer listener.Close()

		for {
			conn, err := listener.Accept()

			if err != nil {
				fmt.Println(err)
				return
			}

			session := createSession(conn, &server)
			server.sessionHandler.OnConnect(session)
			session.Run()
		}
	}()

	return true
}
