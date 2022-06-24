package conn

import (
	"context"
	"net"
)

type ConnectionManger interface {

	//Register 注册
	Register(context.Context, *connection)

	//UnRegister 取消注册
	UnRegister(context.Context, *connection)

	// Send send data
	Send([]byte) error
}

func NewConnectionManager() ConnectionManger {
	return &connectionManager{
		conns:            make(map[string][]*connection),
		chatMessageQueue: NewChatMessageQueue(),
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
}

func (cm *connectionManager) Register(ctx context.Context, conn *connection) {
	cm.conns[conn.uniqueId] = append(cm.conns[conn.uniqueId], conn)
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

}

// Send send data
func (cm *connectionManager) Send(data []byte) error {
	return cm.chatMessageQueue.Send(data)
}
