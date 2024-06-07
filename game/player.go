package game

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Player struct {
	name string
	conn *websocket.Conn
	game Game
	id   uuid.UUID
}
