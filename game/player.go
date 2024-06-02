package game

import "github.com/gorilla/websocket"

type Player struct {
	name string
	conn *websocket.Conn
}
