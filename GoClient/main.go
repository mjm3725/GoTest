package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	startTime := time.Now()

	for i := 0; i < 100000; i++ {
		var s string
		//		_, err := fmt.Scanln(&s)

		//		if err != nil {
		//			fmt.Println(err)
		//			return
		//		}

		s = strconv.Itoa(i)
		sendBuf := []byte(s)
		sendBuf = append(sendBuf, 0)

		conn.Write(sendBuf)

		//fmt.Println("send: " + s)

		recvBuf := make([]byte, 4096)

		_, err := conn.Read(recvBuf)

		if err != nil {
			fmt.Println(err)
			return
		}

		//fmt.Println("recv: " + string(recvBuf[:l]))
	}

	fmt.Println(time.Since(startTime))

	for {
		time.Sleep(time.Second * 2)
	}
}
