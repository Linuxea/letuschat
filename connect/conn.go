package conn

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"linuxea.com/letuschat/dispatch"
)

type ConnectionManger interface {

	//Register 注册
	Register(context.Context, *connection)

	//UnRegister 取消注册
	UnRegister(context.Context, *connection)

	// Send send data
	Send([]byte) error

	// LocalSend local send data
	LocalSend([]byte) error

	HandleConn(conn *net.Conn)
}

func NewConnectionManager() ConnectionManger {
	return &connectionManager{
		conns:            make(map[string][]*connection),
		chatMessageQueue: NewChatMessageQueue(),
		dis:              dispatch.Newdispatch(),
	}
}

// net.Conn wrapper
type connection struct {
	conn     *net.Conn
	uniqueId string
}

// ConnectionManger implement
type connectionManager struct {
	conns            map[string][]*connection // one 2 many
	chatMessageQueue ChatMessageQueue
	dis              *dispatch.Dispatch
}

func (cm *connectionManager) Register(ctx context.Context, conn *connection) {
	fmt.Println("注册", conn.uniqueId, (*conn.conn).LocalAddr().String())
	cm.conns[conn.uniqueId] = append(cm.conns[conn.uniqueId], conn)
	cm.dis.StoreConf(conn.uniqueId, "127.0.0.1:9090")
}

func (cm *connectionManager) UnRegister(ctx context.Context, conn *connection) {

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

	cm.dis.DeleteConf(conn.uniqueId, (*conn.conn).LocalAddr().String())

}

// Send send data
func (cm *connectionManager) Send(data []byte) error {
	return cm.chatMessageQueue.Send(data)
}

func (cm *connectionManager) LocalSend(date []byte) error {
	var m map[string]string
	json.Unmarshal(date, &m)
	fmt.Println("本地数据", m, cm.conns)
	for _, v := range cm.conns[m["To"]] {
		_, err := (*v.conn).Write([]byte("你好啊，现在是几点了"))
		if err != nil {
			fmt.Println("发送数据异常", err)
		}
	}
	return nil
}

func (cm *connectionManager) HandleConn(conn *net.Conn) {

	uniqueId := "linuxea"
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
