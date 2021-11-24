package rabbitmq

import (
	"awesomeProject1/datamodels"
	"awesomeProject1/services"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"sync"
)

//url 格式amqp://账号:密码@rabbitmq服务器地址:端口号/vhost
const MQURL = "amqp://lpxlpx:lpxlpx@127.0.0.1:5672/lpx"

type RabbitMQ struct {
	conn *amqp.Connection
	channel *amqp.Channel
	//队列名称
	QueueName string
	//交换机
	Exchange string
	//key
	Key string
	//连接信息
	Mqurl string
	sync.Mutex
}

func NewRabbitMQ(QueueName string ,Exchange string, Key string) *RabbitMQ{
	rabbitmq := &RabbitMQ{QueueName: QueueName,Exchange: Exchange,Key:Key,Mqurl: MQURL}
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err ,"创建链接错误")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "获取channel失败")
	return rabbitmq
}

//断开channel和connection
func (r *RabbitMQ) Destory(){
	r.channel.Close()
	r.conn.Close()
}

func (r *RabbitMQ) failOnErr(err error, message string){
	if err != nil {
		log.Fatalln("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s",message,err))
	}
}

//
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	return NewRabbitMQ(queueName,"","")
}

//简单模式下生产代码
func (r *RabbitMQ) PublishSimple(message string) error {
	r.Lock()
	defer r.Unlock()
	//1.申请队列,如果队列不存在，则创建
	//保证队列存在
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否为自动删除
		false,
		//是否具有排他性
		false,
		//是否堵塞
		false,
		nil,
		)
	if err != nil{
		fmt.Println(err)
	}
	//2.发送消息
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		//如果为true，如果无法找到，会返回给发送者
		false,
		//当不存在消费者，返回给消费者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(message),
		},
		)
	return nil
}

func (r *RabbitMQ)ConsumeSimple(orderService services.IOrderService,productService services.IProductService){
	//1.申请队列,如果队列不存在，则创建
	//保证队列存在
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否为自动删除
		false,
		//是否具有排他性
		false,
		//是否堵塞
		false,
		nil,
	)
	if err != nil{
		fmt.Println(err)
	}
	//消费者控量
	r.channel.Qos(
		1,//当前消费者一次能接受最大消息数量
		0,//服务器传递的最大容量
		false,//如果设置为true，对channel可用
		)
	//接受消息
	msgs, err := r.channel.Consume(
		r.QueueName,
		//用来区分多个消费者
		"",
		//是否自动应答
		false,
		//是否排他
		false,
		//如果设置为true，表示不能将同一个connection中的消息传递给这个connection消费者
		false,
		false,
		nil,
		)
	if err != nil{
		fmt.Println(err)
	}
	forever := make(chan bool)
	//启动协程处理消息
	go func() {
		for d := range msgs{
			//实现我们要处理的逻辑函数
			log.Println("received a message:%s",d.Body)
			message := &datamodels.Message{}
			err := json.Unmarshal([]byte(d.Body),message)
			if err != nil{
				fmt.Println(err)
			}
			//插入订单
			_, err = orderService.InsertOrderByMessage(message)
			if err != nil{
				fmt.Println(err)
			}
			//扣除商品数量
			err = productService.SubNumberOne(message.ProductID)
			if err != nil{
				fmt.Println(err)
			}
			//如果为true 表示确认所有未确认的消息，为false表示确认当前消息
			d.Ack(false)
		}
	}()
	log.Println("[*] Waiting for messages, To exit press CTRL+C")
	<-forever
}