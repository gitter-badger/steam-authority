package main

import (
	"io/ioutil"
	"log"

	"github.com/Philipp15b/go-steam"
	"github.com/Philipp15b/go-steam/protocol/steamlang"
)

func steamx() {
	myLoginInfo := new(steam.LogOnDetails)
	myLoginInfo.Username = "Your username"
	myLoginInfo.Password = "Your password"

	client := steam.NewClient()
	client.Connect()
	for event := range client.Events() {
		switch e := event.(type) {
		case *steam.ConnectedEvent:
			client.Auth.LogOn(myLoginInfo)
		case *steam.MachineAuthUpdateEvent:
			ioutil.WriteFile("sentry", e.Hash, 0666)
		case *steam.LoggedOnEvent:
			client.Social.SetPersonaState(steamlang.EPersonaState_Online)
		case steam.FatalErrorEvent:
			log.Print(e)
		case error:
			log.Print(e)
		}
	}
}
