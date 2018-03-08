package queue2

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/steam-authority/steam-authority/datastore"
	"github.com/streadway/amqp"
)

var (
	changeCloseChannel chan *amqp.Error
)

func getQueue() (conn *amqp.Connection, ch *amqp.Channel, q amqp.Queue, err error) {

	changeCloseChannel = make(chan *amqp.Error)

	conn, err = amqp.Dial(os.Getenv("STEAM_RABBIT"))
	conn.NotifyClose(changeCloseChannel)
	if err != nil {
		logger.Error(err)
	}

	ch, err = conn.Channel()
	if err != nil {
		logger.Error(err)
	}

	q, err = ch.QueueDeclare(
		namespace+"Changes", // name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		logger.Error(err)
	}

	return conn, ch, q, err
}

func SendChange(change *datastore.Change) (err error) {

	conn, ch, q, err := getQueue()
	defer conn.Close()
	defer ch.Close()
	if err != nil {
		return err
	}

	changeJSON, err := json.Marshal(change)
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
			Body:         changeJSON,
		})
	if err != nil {
		logger.Error(err)
	}

	return nil
}

func receiveChanges() {

	for {
		fmt.Println("Getting change messages")

		conn, ch, q, err := getQueue()

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

		for {
			// todo, these breaks need to break out twice
			select {
			case err = <-changeCloseChannel:
				fmt.Println("change channel closed")
				break
			case msg := <-msgs:
				err := processChange(msg)
				if err != nil {
					fmt.Println("change process error:")
					fmt.Println(err.Error())
					break
				}
			}
		}

		conn.Close()
		ch.Close()
	}
}

func processChange(msg amqp.Delivery) (err error) {
	fmt.Println(string(msg.Body))
}

type changeWebsocketPayload struct {
	ID            int                             `json:"id"`
	CreatedAt     int64                           `json:"created_at"`
	CreatedAtNice string                          `json:"created_at_nice"`
	Apps          []changeAppWebsocketPayload     `json:"apps"`
	Packages      []changePackageWebsocketPayload `json:"packages"`
}

type changeAppWebsocketPayload struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type changePackageWebsocketPayload struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
