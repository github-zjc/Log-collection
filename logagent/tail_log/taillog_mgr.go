package taillog

import (
	"fmt"
	"logagent/etcd"
	"time"
)

var tskMgr *taillogMgr

type taillogMgr struct {
	logEntry    []*etcd.LogEntry
	tskMap      map[string]*TailTask
	newConfChan chan []*etcd.LogEntry
}

func InitMgr(logEntryConf []*etcd.LogEntry) {
	tskMgr = &taillogMgr{ //把当前的所有的配置信息保存起来
		logEntry:    logEntryConf,
		tskMap:      make(map[string]*TailTask, 16),
		newConfChan: make(chan []*etcd.LogEntry),
	}
	for _, v := range tskMgr.logEntry { //分别执行每一个配置信息
		//起了多少个tailtask对象都要记录，为了后续于新的配置比较
		tailobj := NewTailTask(v.Topic, v.Path)
		mk := fmt.Sprintf("%s_%s", v.Path, v.Topic)
		tskMgr.tskMap[mk] = tailobj
	}
	go tskMgr.run() //一直等最新的配置
}

//监听自己的newconfchan，有了新的配置之后就做对应的处理
func (t *taillogMgr) run() {
	for {
		select {
		case newConf := <-t.newConfChan:
			fmt.Println("新的配置来了", newConf)
			for _, conf := range newConf {
				mk := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
				_, ok := tskMgr.tskMap[mk]
				if ok {
					//原来有这个配置
					continue
				} else {
					tailobj := NewTailTask(conf.Path, conf.Topic)
					tskMgr.tskMap[mk] = tailobj
				}
			}
			//找出原t.logEntry来有,但是newConf没有的，要删掉
			for _, c1 := range t.logEntry {
				isDelete := true
				for _, c2 := range newConf {
					if c1.Path == c2.Path && c1.Topic == c2.Topic {
						//原来有的现在也有
						isDelete = false
						continue
					}
				}
				if isDelete {
					//停止对应的tailtask对象
					mk := fmt.Sprintf("%s_%s", c1.Path, c1.Topic)
					t.tskMap[mk].cancelFunc()
				}
			}

		default:
			time.Sleep(time.Millisecond * 50)
		}
	}
}

//向外面的包暴露tskMgr的newConfChan
func NewConfChan() chan<- []*etcd.LogEntry {
	return tskMgr.newConfChan
}
