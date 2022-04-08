package main

import (
	"fmt"
	"logagent/config"
	"logagent/etcd"
	"logagent/kafka"
	taillog "logagent/tail_log"
	"logagent/utils"
	"sync"
	_ "time"

	"gopkg.in/ini.v1"
)

var (
	wg sync.WaitGroup
	p  = new(config.Config)
)

func main() {
	//获取Config配置中的信息
	err := ini.MapTo(p, "./config/config.ini")
	if err != nil {
		fmt.Println("load init faild err=", err)
		return
	}
	//kafka初始化
	err = kafka.Init(p.KafkaConf.Address, p.KafkaConf.ChanMaxSize)
	if err != nil {
		fmt.Println("kafka init faild err=", err)
		return
	}
	fmt.Println("kafka 初始化成功")
	//etcd 初始化
	err = etcd.Init(p.EtcdConf.Address, p.EtcdConf.TimeOut)
	if err != nil {
		fmt.Println("etcd init faild err=", err)
		return
	}
	fmt.Println("etcd 初始化成功")

	//获取IP、
	ip, err := utils.GetOutboundIP()
	if err != nil {
		fmt.Println("get ip faild err=", err)
		return
	}
	key := fmt.Sprintf(p.EtcdConf.Key, ip)
	fmt.Println(key)
	//从etcd中获取配置项
	logEntryConf, err := etcd.GetConf(key)
	if err != nil {
		fmt.Println("etcd.GetConf err=", err)
	}
	for index, v := range logEntryConf {
		fmt.Printf("index=%v logconf=%s ", index, v)
	}

	//将etcd中的配置向取出并进行日志收集，发送到kafka中
	taillog.InitMgr(logEntryConf)

	//派一个goroutine去监听配置项的改变
	newconf := taillog.NewConfChan()
	wg.Add(1)
	go etcd.WatchConf(p.EtcdConf.Key, newconf)
	wg.Wait()
}
