package conn

import (
	"context"
	"flag"
	"fmt"
	"io"
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
			panic(err.Error())
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

	size := make([]byte, 1024)
	for {
		select {
		case <-time.After(time.Duration(30) * time.Minute):
			cm.UnRegister(context.TODO(), wrapperConn)
			break
		default:
			len, err := (*conn).Read(size)
			if err == io.EOF {
				cm.UnRegister(context.TODO(), wrapperConn)
				break
			}

			if err != nil {
				break
			}

		}

	}
}
