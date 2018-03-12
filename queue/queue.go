package queue

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/streadway/amqp"
)

const (
	namespace = "STEAM_"

	QueueChanges  = "Changes"
	QueueApps     = "Apps"
	QueuePackages = "Packages"
	QueuePlayers  = "Players"
)

var (
	queues map[string]queue
)

func init() {

	qs := []queue{
		{Name: QueueChanges, Callback: processChange},
		{Name: QueueApps, Callback: processApp},
		{Name: QueuePackages, Callback: processPackage},
		{Name: QueuePlayers, Callback: processPlayer},
	}

	queues = make(map[string]queue)
	for _, v := range qs {
		queues[v.Name] = v
	}
}

func RunConsumers() {

	for _, v := range queues {
		go v.consume()
	}
}

func Produce(queue string, data []byte) (err error) {

	if val, ok := queues[queue]; ok {
		return val.produce(data)
	}

	return errors.New("no such queue")
}

type queue struct {
	Name     string
	Callback func(msg amqp.Delivery) (err error)
}

func (s queue) getConnection() (conn *amqp.Connection, ch *amqp.Channel, q amqp.Queue, closeChannel chan *amqp.Error, err error) {

	closeChannel = make(chan *amqp.Error)

	conn, err = amqp.Dial(os.Getenv("STEAM_AMQP"))
	conn.NotifyClose(closeChannel)
	if err != nil {
		logger.Error(err)
	}

	ch, err = conn.Channel()
	if err != nil {
		logger.Error(err)
	}

	q, err = ch.QueueDeclare(
		namespace+s.Name, // name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		logger.Error(err)
	}

	return conn, ch, q, closeChannel, err
}

func (s queue) produce(data []byte) (err error) {

	conn, ch, q, _, err := s.getConnection()
	defer conn.Close()
	defer ch.Close()
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         data,
		})
	if err != nil {
		logger.Error(err)
	}

	return nil

}

func (s queue) consume() (err error) {

	for {
		fmt.Println("Getting " + s.Name + " messages")

		conn, ch, q, closeChan, err := s.getConnection()

		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			false,  // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		if err != nil {
			logger.Error(err)
		}

		var breakFor = false

		for {
			select {
			case err = <-closeChan:

				logger.Info("change channel closed")
				time.Sleep(time.Second * 10)

				breakFor = true
				break

			case msg := <-msgs:
				err := s.Callback(msg)
				if err != nil {

					logger.Info("change process error:")
					logger.Error(err)

					breakFor = true
					break

				}
			}

			if breakFor {
				break
			}
		}

		conn.Close()
		ch.Close()
	}
}
