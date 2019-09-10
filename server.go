package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	port := "9090"
	Start(port)
}

func Start(port string) {
	host := ":" + port

	//获取tcp地址
	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		log.Printf("resolve tcp addr failed:%v", err)
		return
	}
	//监听
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Printf("listen tcp port failed:%v\n", err)
		return
	}

	//建立连接池
	conns := make(map[string]net.Conn)

	//消息通道
	messageChan := make(chan string, 10)

	//广播信息
	go BroadMessages(&conns, messageChan)

	//启动服务器发消息
	go SendMsg1(messageChan)

	//启动
	for {
		fmt.Printf("listening port %s ...\n", port)
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("Acccept failed:%v\n", err)
			continue
		}
		//把连接上的用户放到连接池
		conns[conn.RemoteAddr().String()] = conn
		fmt.Println(conns)

		//处理消息
		go Handler(conn, &conns, messageChan)
	}

}

func Handler(conn net.Conn, conns *map[string]net.Conn, messages chan string) {
	buf := make([]byte, 1024)
	for {
		length, err := conn.Read(buf)
		if err != nil {
			log.Printf("read client message failed:%v", err)
			delete(*conns, conn.RemoteAddr().String())
			conn.Close()
			break
		}
		//把客户端发过来的信息写到通道中
		recvStr := string(buf[0:length])
		messages <- recvStr
	}
}

//广播信息
func BroadMessages(conns *map[string]net.Conn, messages chan string) {
	for {
		//从通道里读取信息
		msg := <-messages
		fmt.Println(msg)

		//将读出来的信息发送给所有客户端
		for k, conn := range *conns {
			_, err := conn.Write([]byte(msg))
			if err != nil {
				log.Printf("broad message to %s failed:%v", err)
				delete(*conns, k)
			}
		}
	}
}

func SendMsg1(messageChan chan string) {
	for {
		var s string
		fmt.Scanln(&s)
		s = "服务器：" + s
		messageChan <- s
	}
}
