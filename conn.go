package letuschat

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
)

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		conns:            make(map[string][]*Connection),
		chatMessageQueue: NewChatMessageQueue(),
		dis:              Newdispatch(),
	}
}

// net.Conn wrapper
type Connection struct {
	conn     *net.Conn
	uniqueId string
}

// ConnectionManger implement
type ConnectionManager struct {
	conns            map[string][]*Connection // one 2 many
	chatMessageQueue ChatMessageQueue
	dis              *Dispatch
}

func (cm *ConnectionManager) Register(ctx context.Context, conn *Connection) {
	cm.conns[conn.uniqueId] = append(cm.conns[conn.uniqueId], conn)
	ipAddr, port, _ := GetOutboundIP()
	cm.dis.Register(conn.uniqueId, fmt.Sprintf("%s:%d", ipAddr.String(), port))
}

func (cm *ConnectionManager) Bind(ctx context.Context, uniqueId string, conn *net.Conn) {

}

func (cm *ConnectionManager) UnRegister(ctx context.Context, conn *Connection) {

	conns := cm.conns[conn.uniqueId]

	findIndex := -1
	for index := range conns {
		if conns[index] == conn {
			findIndex = index
			break
		}
	}

	if findIndex > -1 {
		cm.conns[conn.uniqueId] = append(conns[:findIndex], conns[findIndex+1:]...)
	}

	ipAddr, port, _ := GetOutboundIP()
	cm.dis.Unregister(conn.uniqueId, fmt.Sprintf("%s:%d", ipAddr.String(), port))
}

// Send send data
func (cm *ConnectionManager) Send(data []byte) error {
	return cm.chatMessageQueue.Send(data)
}

func (cm *ConnectionManager) LocalSend(date []byte) error {
	var m map[string]interface{}
	json.Unmarshal(date, &m)
	for _, v := range cm.conns[m["to"].(string)] {
		_, err := (*v.conn).Write([]byte(m["content"].(string)))
		if err != nil {
			fmt.Println("发送数据异常", err)
		}
	}
	return nil
}

func (cm *ConnectionManager) HandleConn(conn *net.Conn) {

	data := make([]byte, 1024)
	for {
		select {
		default:
			len, err := (*conn).Read(data)
			if err != nil {
				// cm.UnRegister(context.TODO(), conn)
				break
			}

			if err = cm.Send(data[:len]); err != nil {

				wrapperConn := &Connection{
					uniqueId: string(data[:len]),
					conn:     conn,
				}
				cm.Register(context.TODO(), wrapperConn)

				fmt.Printf("send data error:%s\n", err.Error())
			}

		}

	}
}
