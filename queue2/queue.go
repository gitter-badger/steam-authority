package queue2

import (
	"fmt"
	"os"
	"time"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/streadway/amqp"
)

const (
	namespace = "STEAM_"
)

const (
	queueChanges  = "Changes"
	queueApps     = "Apps"
	queuePackages = "Packages"
	queuePlayers  = "Players"
)

var (
	queues map[string]queue
)

func Produce(queue string, data []byte) (err error) {
	return queues[queue].send(data)
}

func RunConsumers() {

	queues = map[string]queue{
		queueChanges: {Name: queueChanges, Callback: processChange},
		//queueApps:     {Name: queueApps, Callback: },
		//queuePackages: {Name: queuePackages, Callback: },
		//queuePlayers:  {Name: queuePlayers, Callback: },
	}

	for _, v := range queues {
		go v.receive()
	}
}

type queue struct {
	Name     string
	Callback func(msg amqp.Delivery) (err error)
}

func (s queue) getConnection() (conn *amqp.Connection, ch *amqp.Channel, q amqp.Queue, closeChannel chan *amqp.Error, err error) {

	closeChannel = make(chan *amqp.Error)

	conn, err = amqp.Dial(os.Getenv("STEAM_RABBIT"))
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

func (s queue) send(data []byte) (err error) {

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

func (s queue) receive() (err error) {

	for {
		fmt.Println("Getting change messages")

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
