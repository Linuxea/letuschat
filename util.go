package letuschat

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

// Get preferred outbound ip of this machine
func GetOutboundIP() (net.IP, int, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
		return nil, 0, fmt.Errorf("获取ip异常:%s", err.Error())
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, localAddr.Port, nil
}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
