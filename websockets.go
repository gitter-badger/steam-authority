package main

import "github.com/gorilla/websocket"

const (
	changes = "changes"
)

type webSocket struct {
	section    string
	time       int64
	connection *websocket.Conn
}
