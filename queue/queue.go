package queue

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

var (
	connection   *amqp.Connection
	channel      *amqp.Channel
	closeChannel chan *amqp.Error
	namespace    = "STEAM_"
)

func init() {
	closeChannel = make(chan *amqp.Error)
}

func RunConsumers() {

	//go playerConsumer()
	go appConsumer()
	go changeConsumer()
	go packageConsumer()
}

func connect() (err error) {

	if connection != nil && channel != nil {
		return nil
	}

	connection, err = amqp.Dial(os.Getenv("STEAM_RABBIT"))
	connection.NotifyClose(closeChannel)
	if err != nil {
		return err
	}
	//defer amqpConn.Close()

	channel, err = connection.Channel()
	if err != nil {
		return err
	}
	//defer ch.Close()

	fmt.Println("connected")

	return nil
}

//// todo, Have consume queue and produce queue!!
//
//const (
//	QApp     = "app"
//	QPackage = "package"
//)
//
//var (
//	queues map[string]*queue
//)
//
//func setup() {
//
//	queues = make(map[string]*queue, 0)
//
//	newQueue(QApp)
//	newQueue(QPackage)
//}
//
//type queue struct {
//	connection   *amqp.Connection
//	channel      *amqp.Channel
//	queue        amqp.Queue
//	name         string
//	closeChannel chan *amqp.Error
//}
//
//func newQueue(key string) (q *queue) {
//
//	q = new(queue)
//	q.name = key
//	q.init()
//
//	queues[key] = q
//
//	return q
//}
//
//func (q *queue) runConsumer() {
//
//	go func() {
//
//	}()
//}
//
//func (q *queue) init() (err error) {
//
//	if q.connection != nil {
//		return nil
//	}
//
//	q.closeChannel = make(chan *amqp.Error)
//
//	q.connection, err = amqp.Dial(os.Getenv("STEAM_RABBIT"))
//	q.connection.NotifyClose(q.closeChannel)
//	if err != nil {
//		return err
//	}
//
//	q.channel, err = q.connection.Channel()
//	if err != nil {
//		return err
//	}
//
//	q.queue, err = q.channel.QueueDeclare(
//		"Steam_"+strings.Title(q.name), // name
//		true,                           // durable
//		false,                          // delete when unused
//		false,                          // exclusive
//		false,                          // no-wait
//		nil,                            // arguments
//	)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (q *queue) getQueue() (r amqp.Queue, err error) {
//
//	return r, nil
//}
//
//func (q *queue) produce(bytes []byte) (err error) {
//
//	if q.queue.Name == "" {
//		q.init()
//	}
//
//	err = q.channel.Publish(
//		"",           // exchange
//		q.queue.Name, // routing key
//		false,        // mandatory
//		false,        // immediate
//		amqp.Publishing{
//			DeliveryMode: amqp.Persistent,
//			ContentType:  "text/plain",
//			Body:         bytes,
//		})
//	if err != nil {
//		return err
//	}
//
//	return nil
//
//}
//
//func (q *queue) consume(func(msg amqp.Delivery) (err error)) {
//
//}
