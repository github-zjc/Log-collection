package taillog

import (
	"context"
	"fmt"
	"logagent/kafka"

	"github.com/hpcloud/tail"
)

//将tail做成一个结构体对象
type TailTask struct {
	Topic    string
	Path     string
	instance *tail.Tail
	//对tailtask对象进行停止管理
	context    context.Context
	cancelFunc context.CancelFunc
}

func NewTailTask(topic, path string) (tailObj *TailTask) {
	ctx, cancel := context.WithCancel(context.Background())
	tailObj = &TailTask{
		Topic:      topic,
		Path:       path,
		context:    ctx,
		cancelFunc: cancel,
	}
	tailObj.init() //根据路径打开对应的日志
	return
}

func (t TailTask) init() {
	config := tail.Config{
		ReOpen:    true,                                 //重新打开
		Follow:    true,                                 //是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, //从文件的哪个地方开始读
		MustExist: false,                                //文件不存在报错
		Poll:      true,
	}
	var err error
	t.instance, err = tail.TailFile(t.Path, config)
	if err != nil {
		fmt.Println("tail file failed err", err)
		return
	}
	go t.run() //直接起一个协程去向kafka中写数据
}

//将日志内容写入kafka中
func (t *TailTask) run() {
	for {
		select {
		case <-t.context.Done():
			fmt.Printf("tail task :%s_%s 结束了...\n", t.Path, t.Topic)
			return
		case line := <-t.instance.Lines:
			//将数据写入kafka中
			// kafka.SendToKafka(t.Topic, line.Text) //这一步完成后才可以读下一行，导致效率低
			//优化  先把日志发送到chan中
			kafka.SendToChan(t.Topic, line.Text)
			//kafka哪个包中有单独的goroutine去取日志数据发送到kafka

		}
	}
}
