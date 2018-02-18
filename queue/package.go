package queue

import (
	"fmt"
	"strconv"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/steam-authority/steam-authority/mysql"
	"github.com/streadway/amqp"
)

func getPackageQueue() (queue amqp.Queue, err error) {

	err = connect()
	if err != nil {
		return queue, err
	}

	queue, err = channel.QueueDeclare(
		namespace+"Packages", // name
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)

	return queue, err
}

func PackageProducer(id int, change int) (err error) {

	logger.Info("Adding package " + strconv.Itoa(id) + " to rabbit")

	queue, err := getPackageQueue()
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
			ContentType:  "text/plain",
			Body:         []byte(strconv.Itoa(id)),
		})
	if err != nil {
		return err
	}

	return nil
}

func packageConsumer() {

	for {

		queue, err := getPackageQueue()
		if err != nil {
			logger.Error(err)
			continue
		}

		fmt.Println("Getting package messages from rabbit")
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

				err := mysql.ConsumePackage(msg)
				if err != nil {
					logger.Error(err)
				} else {
					msg.Ack(false)
				}
			}
		}
	}
}
