package conf

type Conf struct {
	KafkaConf `ini:"kafka"`
	EsConf    `ini:"es"`
}

type KafkaConf struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}

type EsConf struct {
	Address string `ini:"address"`
	GoNums  int    `ini:"gonums"`
}
