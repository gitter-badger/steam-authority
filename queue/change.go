package queue

import (
	"encoding/json"
	"strconv"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/steam-authority/steam-authority/datastore"
	"github.com/steam-authority/steam-authority/mysql"
	"github.com/steam-authority/steam-authority/websockets"
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

		//fmt.Println("Getting change messages from rabbit")
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

				change, err := datastore.ConsumeChange(msg)
				if err != nil {
					logger.Error(err)
				} else {
					msg.Ack(false)
				}

				// Send websocket
				if websockets.HasConnections() {

					// Get apps for change
					var apps []changeAppWebsocketPayload
					appsResp, err := mysql.GetApps(change.Apps, []string{"id", "name"})
					if err != nil {
						logger.Error(err)
					}

					for _, v := range appsResp {
						apps = append(apps, changeAppWebsocketPayload{
							ID:   v.ID,
							Name: v.GetName(),
						})
					}

					// Get packages for change
					var packages []changePackageWebsocketPayload
					packagesResp, err := mysql.GetPackages(change.Packages, []string{"id", "name"})
					if err != nil {
						logger.Error(err)
					}

					for _, v := range packagesResp {
						packages = append(packages, changePackageWebsocketPayload{
							ID:   v.ID,
							Name: v.GetName(),
						})
					}

					payload := changeWebsocketPayload{
						ID:        change.ChangeID,
						CreatedAt: change.CreatedAt.Unix(),
						Apps:      apps,
						Packages:  packages,
					}
					websockets.Send(websockets.CHANGES, payload)
				}
			}
		}
	}
}

type changeWebsocketPayload struct {
	ID        int                             `json:"id"`
	CreatedAt int64                           `json:"created_at"`
	Apps      []changeAppWebsocketPayload     `json:"apps"`
	Packages  []changePackageWebsocketPayload `json:"packages"`
}

type changeAppWebsocketPayload struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type changePackageWebsocketPayload struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
