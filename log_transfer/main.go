package main

import (
	"fmt"
	"log_transfer/conf"
	"log_transfer/es"
	"log_transfer/kafka"
	"sync"

	"gopkg.in/ini.v1"
)

var (
	wg sync.WaitGroup
)

func main() {
	//获取conf信息
	var p conf.Conf
	fmt.Printf("%p\n", &p)
	err := ini.MapTo(&p, "./conf/conf.ini")
	if err != nil {
		fmt.Println("ini.MapTo err=", err)
		return
	}
	fmt.Printf("%p~~~\n", &p)
	fmt.Println(p)
	fmt.Printf("es address=%v kafka address=%v kafka topic=%v \n", p.EsConf.Address, p.KafkaConf.Address, p.KafkaConf.Topic)

	//初始化es
	err = es.InitES(p.EsConf.Address, p.EsConf.GoNums)
	if err != nil {
		fmt.Println("es.Init err=", err)
		return
	}
	//初始化kafka
	wg.Add(1)
	err = kafka.InitKafka([]string{p.KafkaConf.Address}, p.KafkaConf.Topic)
	if err != nil {
		fmt.Println("kafka.Init err=", err)
		return
	}
	wg.Wait()

}
