package letuschat

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func Newdispatch(cm *ConnectionManager) *Dispatch {
	return &Dispatch{
		redisCli: redis.NewClient(&redis.Options{
			Addr:         ":49155",
			Password:     "redispw",
			MaxRetries:   3,
			DB:           0,
			ReadTimeout:  time.Second * 1,
			WriteTimeout: time.Second * 1,
		}),
		cm: cm,
	}
}

type Dispatch struct {
	redisCli *redis.Client
	cm       *ConnectionManager
}

func (d *Dispatch) Register(uniqueId, address string) {
	d.redisCli.SAdd(context.TODO(), uniqueId, address)
}

func (d *Dispatch) Unregister(uniqueId, address string) {
	d.redisCli.SRem(context.TODO(), uniqueId, address)
}

func (d *Dispatch) loadConf(uniqueId string) []string {
	s, _ := d.redisCli.SMembers(context.TODO(), uniqueId).Result()
	return s
}

func (d *Dispatch) invokeDispatch(data interface{}) error {
	dataInter := data.(map[string]interface{})
	toField := dataInter["to"].(string)
	for _, v := range d.loadConf(toField) {
		b, _ := GetBytes(v)
		_ = d.cm.LocalSend(b)
	}

	return nil

}

func (d *Dispatch) Listen() {
	_ = d.redisCli.XGroupCreate(context.TODO(), "chat", "messageGroup", "0")
	for {
		readGroup, err := d.redisCli.XReadGroup(context.TODO(), &redis.XReadGroupArgs{
			Group:    "messageGroup",
			Consumer: "little boy",
			Streams:  []string{"chat", ">"},
			Count:    1000,
			NoAck:    true,
		}).Result()

		if err != nil {
			fmt.Println("read group error", err.Error())
			continue
		}

		for idx := range readGroup[0].Messages {
			d.invokeDispatch(readGroup[0].Messages[idx].Values)
		}

	}

}
