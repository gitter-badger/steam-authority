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
	fmt.Println("## Running consumers")

	forever := make(chan bool)

	//go playerConsumer()
	go appConsumer()
	go changeConsumer()

	<-forever
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
