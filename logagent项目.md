# LogAgent项目

## 项目架构设计：

![image-20220320203222171](C:\Users\DELL\AppData\Roaming\Typora\typora-user-images\image-20220320203222171.png)

## 项目工作流程

1：解析配置文件数据----ini第三方库

2：读日志----tail第三方库

3：往kafka里面写日志

## 项目进展一：

### 学习内容：

#### Kafka

![image-20220320201607086](C:\Users\DELL\AppData\Roaming\Typora\typora-user-images\image-20220320201607086.png)

1：Kafka集群架构：

​		1：broker

​		2：topic

​		3：partition分区，把同一个topic分成不同的分区，高负载

​				1：leader ：分区的主节点

​				2：follower：分区的从节点

​		4：Consumer Group

2.生产者往Kafka发数据的流程（6步）

![image-20220320201747532](C:\Users\DELL\AppData\Roaming\Typora\typora-user-images\image-20220320201747532.png)

3：Kafka选择分区的模式（3种）

​		1：指定往那个分区写

​		2：指定key，kafka根据key做hash然后决定写那个分区

​		3：轮询方式

4：生产者往kafka中发送数据的模式（3种）

​		1：`0` ：把数据发送给leader就成功，效率高，安全性低

​		2： `1`：把数据发送给leader，等待leader回复ack

​		3：`all`：把数据发送给leader，所有的follower从leader中拉去数据回复ack给leader，leader在向生产者发送ack；安全性最高

5：分区存储文件的原理

![image-20220320202758385](C:\Users\DELL\AppData\Roaming\Typora\typora-user-images\image-20220320202758385.png)

![image-20220320202814315](C:\Users\DELL\AppData\Roaming\Typora\typora-user-images\image-20220320202814315.png)

![image-20220320202835915](C:\Users\DELL\AppData\Roaming\Typora\typora-user-images\image-20220320202835915.png)

6：为什么kafka快

![image-20220320202852052](C:\Users\DELL\AppData\Roaming\Typora\typora-user-images\image-20220320202852052.png)

7：消费者组消费数据原理

![image-20220320202902531](C:\Users\DELL\AppData\Roaming\Typora\typora-user-images\image-20220320202902531.png)

![image-20220320202921530](C:\Users\DELL\AppData\Roaming\Typora\typora-user-images\image-20220320202921530.png)



8：kafka使用场景

![image-20220320203005423](C:\Users\DELL\AppData\Roaming\Typora\typora-user-images\image-20220320203005423.png)

![image-20220320203017824](C:\Users\DELL\AppData\Roaming\Typora\typora-user-images\image-20220320203017824.png)

#### Zookeeper

![image-20220320203143963](C:\Users\DELL\AppData\Roaming\Typora\typora-user-images\image-20220320203143963.png)

### 完成内容

程序入口：通过第三方库ini：获取配置文件数据、初始化kfka、tail、编写Run函数将pkg tail中得到的数据写入kafka中

package kafka：编写初始化函数、向kafka内写入数据的函数

package taillog：编写初始化函数、返回值为只读管道的函数 (<-chan *tail.Line)

总结：

​			项目初步实现将日志写入kafka

​			不足之处：

​							当修改config文件的时候，程序需要重新启动，并且不能够对多个日志进行同时操作

​			改进方案：

​							通过添加etcd将tail需要收集的指标存入etcd中，并对其watch从而动态的改变tail对日志收集的策略

