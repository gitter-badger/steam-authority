package datastore

import (
	"strconv"

	"github.com/streadway/amqp"
)

func ConsumeApp(msg amqp.Delivery) (err error) {

	id := string(msg.Body)

	idx, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	_, err = GetArticlesFromSteam(idx)
	if err != nil {
		return err
	}

	return nil
}
