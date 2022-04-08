日志收集系统  logagent
1.先将需要收集的配置项，配置到config.ini中
2.程序入口获取config配置项、初始化kafka、etcd
3.向etcd中写入 key value
4.将etcd的配置项取出进行日志收集，并写入到Kafka中
5.将Kafka中的数据消费到es中