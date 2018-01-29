package websockets

import (
	"math/rand"
	"net/http"
	"strings"

	"github.com/Jleagle/go-helpers/logger"
	"github.com/gorilla/websocket"
)

const (
	CHANGES = "changes"
	CHAT = "chat"
)

var wsConnections map[int]websocketStruct
var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Send(page string, data interface{}) {

	// count := len(wsConnections)
	// if count > 0 {
	// 	fmt.Println("Sending websocket to " + strconv.Itoa(count) + " connections")
	// }

	ws := websocketPayload{}
	ws.Page = page
	ws.Data = data

	for k, v := range wsConnections {
		err := v.connection.WriteJSON(ws)
		if err != nil {

			// Clean up old connections
			// todo, tidy with https://github.com/gorilla/websocket/issues/104
			if strings.Contains(err.Error(), "broken pipe") {
				v.connection.Close()
				_, ok := wsConnections[k]
				if ok {
					delete(wsConnections, k)
				}

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

	structx := websocketStruct{}
	structx.connection = connection

	if wsConnections == nil {
		wsConnections = make(map[int]websocketStruct)
	}

	wsConnections[rand.Int()] = structx
}

// Properties must be exported so websocket can read them.
type websocketPayload struct {
	Data interface{}
	Page string
}

type websocketStruct struct {
	id         int
	connection *websocket.Conn
}
