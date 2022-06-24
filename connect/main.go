package conn

import (
	"context"
	"flag"
	"fmt"
	"net"
	"time"
)

var cm = NewConnectionManager()

func main() {

	port := flag.Int("port", 8080, "tcp listen port")
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err.Error())
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("accept error:%s\n", err.Error())
		}
		go handleConn(&conn)
	}

}

func handleConn(conn *net.Conn) {

	uniqueId := fmt.Sprintf("%d", time.Now().UnixNano())
	wrapperConn := &connection{
		uniqueId: uniqueId,
		conn:     conn,
	}
	cm.Register(context.TODO(), wrapperConn)

	data := make([]byte, 1024)
	for {
		select {
		case <-time.After(time.Duration(30) * time.Minute):
			cm.UnRegister(context.TODO(), wrapperConn)
		default:
			len, err := (*conn).Read(data)
			if err != nil {
				cm.UnRegister(context.TODO(), wrapperConn)
				break
			}

			if err = cm.Send(data[:len]); err != nil {
				fmt.Printf("send data error:%s\n", err.Error())
			}

		}

	}
}
