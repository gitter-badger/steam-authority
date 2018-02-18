package queue

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/steam-authority/steam-authority/datastore"
	"github.com/streadway/amqp"
)

func getChangeQueue() (queue amqp.Queue, err error) {

	err = connect()
	if err != nil {
		return queue, err
	}

	queue, err = channel.QueueDeclare(
		namespace+"Changes", // name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)

	return queue, err
}

func ChangeProducer(change *datastore.Change) (err error) {

	logger.Info("Adding change " + strconv.Itoa(change.ChangeID) + " to rabbit")

	queue, err := getChangeQueue()
	if err != nil {
		return err
	}

	changeJSON, err := json.Marshal(*change)
	if err != nil {
		return err
	}

	err = channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         changeJSON,
		})
	if err != nil {
		return err
	}

	return nil
}

func changeConsumer() {

	for {

		queue, err := getChangeQueue()
		if err != nil {
			logger.Error(err)
			continue
		}

		fmt.Println("Getting change messages from rabbit")
		messages, err := channel.Consume(
			queue.Name, // queue
			"",         // consumer
			false,      // auto-ack
			false,      // exclusive
			false,      // no-local
			false,      // no-wait
			nil,        // args
		)
		if err != nil {
			logger.Error(err)
			continue
		}

		for {
			select {
			case err = <-closeChannel:
				break
			case msg := <-messages:

				err := datastore.ConsumeChange(msg)
				if err != nil {
					logger.Error(err)
				} else {
					msg.Ack(false)
				}
			}
		}
	}
}
