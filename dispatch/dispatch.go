package dispatch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

func Newdispatch() *Dispatch {
	return &Dispatch{
		redisCli: redis.NewClient(&redis.Options{
			Addr:         ":49155",
			Password:     "redispw",
			MaxRetries:   3,
			DB:           0,
			ReadTimeout:  time.Second * 3,
			WriteTimeout: time.Second * 3,
			MaxConnAge:   time.Hour * 24 * 365,
		}),
	}
}

type Dispatch struct {
	redisCli *redis.Client
}

func (d *Dispatch) StoreConf(uniqueId, addr string) {
	d.redisCli.SAdd(context.TODO(), uniqueId, addr)
}

func (d *Dispatch) DeleteConf(uniqueId, addr string) {
	d.redisCli.SRem(context.TODO(), uniqueId, addr)
}

func (d *Dispatch) LoadConf(uniqueId string) []string {
	s, _ := d.redisCli.SMembers(context.TODO(), uniqueId).Result()
	return s
}

type message struct {
	From, To string
	// Content  string
}

func (d *Dispatch) invokeDispatch(data interface{}) error {

	fmt.Println("接收到的消息", data)
	dataField := data.(map[string]interface{})["data"]
	fmt.Println("dataField", dataField)

	var m message
	if err := json.Unmarshal([]byte(dataField.(string)), &m); err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("m", m.To)

	for _, v := range d.LoadConf(m.To) {
		fmt.Println("请求", v)
		by, _ := json.Marshal(m)
		res, err := http.Post("http://"+v+"/dispatch", "applicaton/json", bytes.NewReader(by))
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		rc := res.Body
		b, err := ioutil.ReadAll(rc)
		if err != nil {
			fmt.Println(err.Error())
			rc.Close()
			continue
		}

		fmt.Println(string(b))
		rc.Close()
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
			Count:    1,
			NoAck:    true,
		}).Result()

		if err != nil {
			fmt.Println("read group error", err.Error())
			continue
		}

		data := readGroup[0].Messages[0]
		d.invokeDispatch(data.Values)
	}

}
