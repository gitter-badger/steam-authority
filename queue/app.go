package queue

import (
	"strconv"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/steam-authority/steam-authority/datastore"
	"github.com/steam-authority/steam-authority/mysql"
	"github.com/streadway/amqp"
)

func getAppQueue() (queue amqp.Queue, err error) {

	err = connect()
	if err != nil {
		return queue, err
	}

	queue, err = channel.QueueDeclare(
		namespace+"Apps", // name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)

	return queue, err
}

func AppProducer(id int, change int) (err error) {

	logger.Info("Adding app " + strconv.Itoa(id) + " to rabbit")

	queue, err := getAppQueue()
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

func appConsumer() {

	for {

		queue, err := getAppQueue()
		if err != nil {
			logger.Error(err)
			continue
		}

		//fmt.Println("Getting app messages from rabbit")
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
				//logger.Info("Reading app " + id + " from rabbit")

				dsErr := datastore.ConsumeApp(msg)
				if err != nil {
					logger.Error(err)
				}

				sqlErr := mysql.ConsumeApp(msg)
				if err != nil && err.Error() != "no app with id" {
					logger.Error(err)
					err = nil
				}

				if dsErr == nil && sqlErr == nil {
					msg.Ack(false)
				}
			}
		}
	}
}
