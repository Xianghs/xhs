package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	Start(os.Args[1])
}

func Start(tcpAddrStr string) {
	//创建
	tcpAddr, err := net.ResolveTCPAddr("tcp4", tcpAddrStr)
	if err != nil {
		log.Printf("Resolve tcp addr failed:%v\n", err)
		return
	}

	//向服务器拨号
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Printf("Dial to server failed:%v\n", err)
		return
	}

	//向服务器发送消息
	go SendMsg(conn)

	//接收服务器消息
	buf := make([]byte, 1024)
	for {
		length, err := conn.Read(buf)
		if err != nil {
			log.Printf("接收服务器消息错误%v\n", err)
			conn.Close()
			os.Exit(0)
			break
		}
		fmt.Println(string(buf[0:length]))
	}
}

func SendMsg(conn net.Conn) {
	// username := conn.LocalAddr().String()
	for {
		var input string
		fmt.Scanln(&input)
		if input == "/q" || input == "/quit" {
			fmt.Println("quit chat...")
			conn.Close()
			os.Exit(0)
		}
		if len(input) > 0 {
			msg := "baba" + ":" + input
			_, err := conn.Write([]byte(msg))
			if err != nil {
				conn.Close()
				break
			}
		}
	}
}
