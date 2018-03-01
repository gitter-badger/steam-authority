package queue

import (
	"strconv"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/steam-authority/steam-authority/datastore"
	"github.com/streadway/amqp"
)

func getFriendsQueue() (queue amqp.Queue, err error) {

	err = connect()
	if err != nil {
		return queue, err
	}

	queue, err = channel.QueueDeclare(
		namespace+"Friends", // name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)

	return queue, err
}

func FriendsProducer(id int) (err error) {

	//logger.Info("Adding friend " + strconv.Itoa(id) + " to rabbit")

	queue, err := getFriendsQueue()
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

func friendsConsumer() {

	for {

		queue, err := getFriendsQueue()
		if err != nil {
			logger.Error(err)
			continue
		}

		//fmt.Println("Getting friends messages from rabbit")
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

				id := string(msg.Body)
				idx, _ := strconv.Atoi(id)

				//logger.Info("Reading friend " + id + " from rabbit")

				player, err := datastore.GetPlayer(idx)
				for _, v := range player.Friends {
					vv, _ := strconv.Atoi(v.SteamID)
					PlayerProducer(vv)
				}

				if err == nil {
					msg.Ack(false)
				} else {
					logger.Error(err)
				}
			}
		}
	}
}
