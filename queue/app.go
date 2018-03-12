package queue

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/steam-authority/steam-authority/datastore"
	"github.com/steam-authority/steam-authority/mysql"
	"github.com/streadway/amqp"
)

func processApp(msg amqp.Delivery) (err error) {

	// Get message payload
	message := new(AppMessage)

	err = json.Unmarshal(msg.Body, message)
	if err != nil {
		msg.Nack(false, false)
		return nil
	}

	// Get news
	_, err = datastore.GetArticlesFromSteam(message.AppID)
	if err != nil {
		logger.Error(err)
	}

	// Update app
	app := new(mysql.App)

	db, err := mysql.GetDB()
	if err != nil {
		logger.Error(err)
	}

	db.Attrs(mysql.GetDefaultAppJSON()).FirstOrCreate(app, mysql.App{ID: message.AppID})

	if message.ChangeID != 0 {
		app.ChangeNumber = message.ChangeID
	}

	err = app.Fill()
	if err != nil {

		if strings.HasSuffix(err.Error(), "connect: connection refused") {
			time.Sleep(time.Second * 1)
			msg.Nack(false, true)
			return nil
		}

		logger.Error(err)
	}

	db.Save(app)
	if db.Error != nil {
		logger.Error(err)
	}

	// Save price change
	err = datastore.CreatePrice(app.ID, app.PriceFinal, app.PriceDiscount)
	if err != nil {
		logger.Error(err)
	}

	// Ack
	msg.Ack(false)
	return nil
}

type AppMessage struct {
	AppID    int
	ChangeID int
}
