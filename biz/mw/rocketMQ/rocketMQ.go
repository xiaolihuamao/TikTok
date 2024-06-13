package rocketmq

import (
	"fmt"
	"log"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"

)



type RabbitMQ struct {
	conn  *.Connection
	mqurl string
}

var Rmq *RabbitMQ

// InitRabbitMQ 初始化RabbitMQ的连接和通道。
func InitRabbitMQ() {

	Rmq = &RabbitMQ{
		mqurl: MQURL,
	}
	dial, err := amqp.Dial(Rmq.mqurl)
	Rmq.failOnErr(err, "创建连接失败")
	Rmq.conn = dial

}

// 连接出错时，输出错误信息。
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s\n", err, message)
		panic(fmt.Sprintf("%s:%s\n", err, message))
	}
}

// 关闭mq通道和mq的连接。
func (r *RabbitMQ) destroy() {
	r.conn.Close()
}
