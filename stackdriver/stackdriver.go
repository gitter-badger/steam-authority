package stackdriver

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/logging"
)

var (
	ctx          = context.Background()
	logMain      *logging.Logger
	logPics      *logging.Logger
	logConsumers *logging.Logger
)

func init() {
	client, err := logging.NewClient(ctx, os.Getenv("STEAM_GOOGLE_PROJECT"))
	if err != nil {
		fmt.Println("Error creating stackdriver client")
	}

	logMain = client.Logger("main")
	logPics = client.Logger("pics")
	logConsumers = client.Logger("consumers")
}

func PicsError(err error) {
	logPics.Log(logging.Entry{Payload: err.Error(), Severity: logging.Error})
}

func PicsMessage(payload string) {
	logPics.Log(logging.Entry{Payload: payload, Severity: logging.Info})
}

func PicsCrit(err error) {
	logPics.LogSync(ctx, logging.Entry{Payload: err.Error(), Severity: logging.Critical})
}
