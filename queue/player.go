package queue

import (
	"strconv"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/steam-authority/steam-authority/datastore"
	"github.com/streadway/amqp"
)

func getPlayerQueue() (queue amqp.Queue, err error) {

	err = connect()
	if err != nil {
		return queue, err
	}

	queue, err = channel.QueueDeclare(
		namespace+"Players", // name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)

	return queue, err
}

func PlayerProducer(id int) (err error) {

	logger.Info("Adding player " + strconv.Itoa(id) + " to rabbit")

	queue, err := getPlayerQueue()
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

func playerConsumer() {

	for {

		queue, err := getPlayerQueue()
		if err != nil {
			logger.Error(err)
			continue
		}

		//fmt.Println("Getting player messages from rabbit")
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

				//id := string(msg.Body)
				//logger.Info("Reading player " + id + " from rabbit")

				err := datastore.ConsumePlayer(msg)
				if err != nil {
					logger.Error(err)
				} else {
					msg.Ack(false)
				}
			}
		}
	}
}
