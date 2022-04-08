package es

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/olivere/elastic/v7"
)

type ESData struct {
	Topic string `json:"topic"`
	Data  string `json:"data"`
}

var (
	client *elastic.Client
	ch     = make(chan *ESData, 100000)
)

func InitES(addr string,nums int) (err error) {
	if !strings.HasPrefix(addr, "http://") {
		addr = "http://" + addr
	}
	sniffOpt := elastic.SetSniff(false) //关闭客户端的嗅探
	client, err = elastic.NewClient(elastic.SetURL(addr), sniffOpt)
	if err != nil {
		// Handle error
		return err
	}

	fmt.Println("connect to es success")
	for i:=0;i<= nums;i++ {
		go SendToES()
	}
	return
}

//将数据存入ES的chan中
func SendToESChan(data *ESData) {
	ch <- data
}

//发送数据到ES中
func SendToES() (err error) {
	for {
		select {
		case msg := <-ch:
			put1, err := client.Index().
				Index(msg.Topic).
				BodyJson(msg).
				Do(context.Background())
			if err != nil {
				// Handle error
				fmt.Println("client.Index err=",err)
			}
			fmt.Printf("Indexed user %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
		default:
			time.Sleep(time.Microsecond)
		}
	}
}
