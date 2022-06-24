package conn

import (
	"encoding/json"
)

// Message interface
type Message interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

type jsonMessage struct {
}

func (*jsonMessage) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (*jsonMessage) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// chat message struct
type chatMessage struct {
	To      string
	ToRoom  string
	Unix    int64
	Content string
}
