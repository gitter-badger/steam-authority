package queue
//
//import (
//	"fmt"
//	"strconv"
//
//	"github.com/Jleagle/go-helpers/logger"
//	"github.com/streadway/amqp"
//)
//
//func getPlayerQueue() (conn *amqp.Connection, chann *amqp.Channel, queue amqp.Queue, err error) {
//
//	conn, chann, err = connect()
//
//	queue, err = chann.QueueDeclare(
//		namespace+"player_updates", // name
//		true,                       // durable
//		false,                       // delete when unused
//		false,                      // exclusive
//		false,                      // no-wait
//		nil,                        // arguments
//	)
//
//	return conn, chann, queue, err
//}
//
//func PlayerProducer(id int) (err error) {
//
//	logger.Info("Adding player " + strconv.Itoa(id) + " to rabbit")
//
//	_, chann, queue, err := getPlayerQueue()
//	if err != nil {
//		return err
//	}
//
//	err = chann.Publish(
//		"",         // exchange
//		queue.Name, // routing key
//		false,      // mandatory
//		false,      // immediate
//		amqp.Publishing{
//			DeliveryMode: amqp.Persistent,
//			ContentType:  "text/plain",
//			Body:         []byte(strconv.Itoa(id)),
//		})
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func playerConsumer() {
//
//	for {
//		_, chann, queue, err := getPlayerQueue()
//		if err != nil {
//			logger.Error(err)
//		}
//
//		fmt.Println("Getting player messages from rabbit")
//		messages, err := chann.Consume(
//			queue.Name, // queue
//			"",         // consumer
//			false,      // auto-ack
//			false,      // exclusive
//			false,      // no-local
//			false,      // no-wait
//			nil,        // args
//		)
//		if err != nil {
//			logger.Error(err)
//		}
//
//		for {
//			select {
//			case err = <-closeChannel:
//				break
//			case msg := <-messages:
//				logger.Info("Received a message: " + string(msg.Body))
//				err := playerWork(msg)
//				if err != nil {
//					logger.Error(err)
//				} else {
//					msg.Ack(false)
//				}
//			}
//		}
//	}
//}
//
//func playerWork(messagex amqp.Delivery) (err error) {
//
//	id := string(messagex.Body)
//	logger.Info("Reading player " + id + " from rabbit")
//	//idx := strconv.Atoi()
//
//	// todo, check if player updated in the last 24 hours and return no errors.
//
//	return nil
//}
