package queue

import (
	"encoding/json"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/steam-authority/steam-authority/mysql"
	"github.com/streadway/amqp"
)

func processPackage(msg amqp.Delivery) (err error) {

	// Get message
	message := new(PackageMessage)

	err = json.Unmarshal(msg.Body, message)
	if err != nil {
		logger.Error(err)
		msg.Nack(false, false)
		return
	}

	// Update package
	db, err := mysql.GetDB()
	if err != nil {
		logger.Error(err)
	}

	pack := new(mysql.Package)

	db.Attrs(mysql.GetDefaultPackageJSON()).FirstOrCreate(pack, mysql.Package{ID: message.PackageID})

	if message.ChangeID != 0 {
		pack.ChangeID = message.ChangeID
	}

	// Move all the stuff in here to queue?
	pack.Fill()

	db.Save(pack)
	if db.Error != nil {
		logger.Error(db.Error)
	}

	msg.Ack(false)
	return nil
}

type PackageMessage struct {
	PackageID int
	ChangeID  int
}
