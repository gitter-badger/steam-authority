package websockets

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/gorilla/websocket"
)

const (
	CHANGES = "changes"
)

var wsConnections []*websocket.Conn
var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Send(data interface{}) {

	fmt.Println("Sending websocket to " + strconv.Itoa(len(wsConnections)) + " connections")

	for k, v := range wsConnections {

		err := v.WriteJSON(data)
		if err != nil {

			// todo, tidy with https://github.com/gorilla/websocket/issues/104
			if strings.Contains(err.Error(), "broken pipe") {
				v.Close()
				wsConnections = append(wsConnections[:k], wsConnections[k+1:]...) // Remove from slice
			} else {
				logger.Error(err)
			}
		}
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {

	// Upgrade the connection
	connection, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	wsConnections = append(wsConnections, connection)
}
