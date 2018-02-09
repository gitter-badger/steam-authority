package queue

import (
	"fmt"
	"os"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/streadway/amqp"
)

var (
	closeChannel chan *amqp.Error
	namespace    = "STEAM_"
)

func init() {
	go playerConsumer()
}

func getChannel() (conn *amqp.Connection, channel *amqp.Channel, err error) {

	//time.Sleep(1 * time.Second)

	closeChannel = make(chan *amqp.Error)

	conn, err = amqp.Dial(os.Getenv("STEAM_RABBIT"))
	conn.NotifyClose(closeChannel)
	if err != nil {
		logger.Error(err)
	}
	//defer amqpConn.Close()

	channel, err = conn.Channel()
	if err != nil {
		logger.Error(err)
	}
	//defer ch.Close()

	fmt.Println("connected")

	return conn, channel, err
}
