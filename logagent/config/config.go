package config

type Config struct {
	KafkaConf   `ini:"kafka"`
	TailLogConf `ini:"taillog"`
	EtcdConf    `ini:"etcd"`
}

type KafkaConf struct {
	//Topic   string   `ini:"topic"`
	Address []string `ini:"address"`
	ChanMaxSize int      `ini:"chan_max_siza"`
}

type TailLogConf struct {
	FileName string `ini:"filename"`
}

type EtcdConf struct {
	Address string `ini:"address"`
	TimeOut int    `ini:"timeout"`
	Key     string `ini:"collect_log_key"`
}
