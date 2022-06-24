package main

import (
	"flag"
	"fmt"
	"net"

	c "linuxea.com/letuschat/connect/conn"
)

func main() {

	port := flag.Int("port", 8080, "tcp listen port")
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err.Error())
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("accept error:%s\n", err.Error())
		}
		go c.NewConnectionManager().HandleConn(&conn)
	}

}
