package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	c "linuxea.com/letuschat/connect"
	"linuxea.com/letuschat/dispatch"
)

func main() {

	cm := c.NewConnectionManager()
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

	// dispath callback
	go func() {
		http.HandleFunc("/dispatch", func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("hello")
			b, _ := ioutil.ReadAll(r.Body)
			fmt.Println("request body", string(b))
			cm.LocalSend(b)

		})
		http.ListenAndServe(fmt.Sprintf(":%d", 9090), nil)
	}()

	go func() {
		dispatch.Newdispatch().Listen()
	}()

	<-make(chan struct{})
}
