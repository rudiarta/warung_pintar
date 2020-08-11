package model

import "github.com/gorilla/websocket"

type Connection struct {
	Ws   *websocket.Conn
	Send chan Message
}
