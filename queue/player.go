package queue

import (
	"encoding/json"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/steam-authority/steam-authority/datastore"
	"github.com/steam-authority/steam-authority/steam"
	"github.com/streadway/amqp"
)

func processPlayer(msg amqp.Delivery) (err error) {

	// Get message
	message := new(PlayerMessage)

	err = json.Unmarshal(msg.Body, message)
	if err != nil {
		logger.Error(err)
		msg.Nack(false, false)
		return
	}

	// Update player
	player, err := datastore.GetPlayer(message.PlayerID)
	if err != nil {
		logger.Error(err)
	}

	err = player.UpdateIfNeeded()
	if err != nil {
		if err.Error() == steam.ErrorInvalidJson {
			// API is probably down
			msg.Nack(false, true)
			return nil
		}
		logger.Error(err)
	}

	msg.Ack(false)
	return nil
}

type PlayerMessage struct {
	PlayerID int
}
