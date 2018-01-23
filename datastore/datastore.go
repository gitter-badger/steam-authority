package datastore

import (
	"context"
	"os"

	"cloud.google.com/go/datastore"
	"github.com/Jleagle/go-helpers/logger"
)

const (
	// CHANGE datastore kind
	CHANGE = "Change"

	// APP atastore kind
	APP = "App"

	// PACKAGE atastore kind
	PACKAGE = "Package"

	// PLAYER atastore kind
	PLAYER = "Player"
)

// todo, return error
func saveKind(key *datastore.Key, data interface{}) (newKey *datastore.Key) {

	client, context := getDSClient()
	newKey, err := client.Put(context, key, data)
	if err != nil {
		logger.Error(err)
	}

	return newKey
}

// todo, cache in global variable
func getDSClient() (*datastore.Client, context.Context) {

	context := context.Background()
	client, err := datastore.NewClient(
		context,
		os.Getenv("STEAM_GOOGLE_PROJECT"),
	)
	if err != nil {
		logger.Error(err)
	}

	return client, context
}
