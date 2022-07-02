package letuschat

import (
	"flag"
	"fmt"
	"net"
)

func main() {

	cm := NewConnectionManager()
	// accept connection
	go func() {
		port := flag.Int("port", 8080, "tcp listen port")
		listen, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
		if err != nil {
			panic(err.Error())
		}
		for {
			conn, _ := listen.Accept()
			go cm.HandleConn(&conn)
		}
	}()

	go func() {
		Newdispatch(cm).Listen()
	}()

	<-make(chan struct{})
}
