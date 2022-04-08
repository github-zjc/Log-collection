package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

var (
	cli *clientv3.Client
)

//需要收集日志信息的结构体
type LogEntry struct {
	Path  string `json:"path"`  //日志存放的路径
	Topic string `json:"topic"` //日志发往Kafka中哪个topic
}

func Init(addr string, timeout int) (err error) {
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: time.Duration(timeout) * time.Second,
	})
	if err != nil {
		fmt.Println("clientv3.New err =", err)
		return err
	}
	return
}

//从etcd中根据key获取配置项
func GetConf(key string) (logEntryConf []*LogEntry, err error) {
	//get
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		// fmt.Printf("%s:%s\n", ev.Key, ev.Value)
		err = json.Unmarshal(ev.Value, &logEntryConf)
		if err != nil {
			fmt.Println("unmarshal etcd value failed ,err", err)
			return
		}
	}
	return
}

//监控etcd配置项
func WatchConf(key string, newConfCh chan<- []*LogEntry) {
	// watch
	rch := cli.Watch(context.Background(), key) // <-chan WatchResponse
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			//通知tailMgr
			//判断操作类型
			var newconf []*LogEntry
			if ev.Type != clientv3.EventTypeDelete {
				err := json.Unmarshal(ev.Kv.Value, &newconf)
				if err != nil {
					fmt.Println("json.Unmarshal err = ", err)
					continue
				}
			}
			fmt.Printf("get new conf :%v\n", newconf)
			newConfCh <- newconf
		}
	}
}
