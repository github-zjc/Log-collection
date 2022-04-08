package kafka

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

type logData struct {
	topic string
	data  string
}

var (
	client      sarama.SyncProducer
	logDataChan chan *logData
)

func Init(addr []string, maxSize int) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	// 连接kafka
	client, err = sarama.NewSyncProducer(addr, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}
	logDataChan = make(chan *logData, maxSize)
	//开启后台goroutine从通道中取数据写入到kafka中
	go SendToKafka()
	return
}

//给外部暴露一个函数，该函数只把日志数据放入到内部的一个chan中
func SendToChan(topic, data string) {
	msg := &logData{
		topic: topic,
		data:  data,
	}
	logDataChan <- msg
}
//真正忘kafka写入数据
func SendToKafka() {
	for {
		select {
		case ld := <-logDataChan:
			// 构造一个消息
			msg := &sarama.ProducerMessage{}
			msg.Topic = ld.topic
			msg.Value = sarama.StringEncoder(ld.data)
			// 发送消息
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				fmt.Println("send msg failed, err:", err)
				return
			}
			fmt.Printf("pid:%v offset:%v msg:%s \n ", pid, offset,msg.Value)
		default:
			time.Sleep(time.Microsecond*50)

		}
	}
}
