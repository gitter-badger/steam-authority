package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/gorilla/websocket"
	"github.com/kr/pretty"
)

const (
	changes = "changes"
)

var ws webSocket

func sendWebsocket(data interface{}) {

	pretty.Print(ws.section)

	if ws.connection == nil {
		fmt.Println("connection is nil")
		return
	}

	err := ws.connection.WriteJSON(data)
	if err != nil {
		logger.Error(err)
	}
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {

	// Make a connection
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		logger.Error(err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	// Save the connection
	newSocket := webSocket{}
	newSocket.connection = conn
	newSocket.time = time.Now().Unix()
	newSocket.section = "changes"

	fmt.Println("Saving websocket to memory")
	ws = webSocket{}
	ws = newSocket
	pretty.Print(ws.section)
}

type webSocket struct {
	section    string
	time       int64
	connection *websocket.Conn
}
