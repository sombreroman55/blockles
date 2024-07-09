package game

import (
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Player struct {
	Name string
	Id   uuid.UUID
	conn *websocket.Conn
	send chan []byte
	game *Game
	hub  *Blockles
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func CreateNewPlayer(name string, hub *Blockles) uuid.UUID {
	log.Debug("Creating new player")
	p := &Player{
		Name: name,
		conn: nil,
		hub:  hub,
		Id:   uuid.New(),
		send: make(chan []byte, 256),
	}
	players[p.Id] = p
	return p.Id
}

func (p *Player) SetName(name string) {
	p.Name = name
}

func (p *Player) AttachWebsocket(conn *websocket.Conn) {
	p.conn = conn
}

func (p *Player) readCommands() {
	defer func() {
		p.hub.unregister <- p
		p.conn.Close()
	}()
	p.conn.SetReadLimit(maxMessageSize)
	p.conn.SetReadDeadline(time.Now().Add(pongWait))
	p.conn.SetPongHandler(func(string) error { p.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := p.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Errorf("error: %v\n", err)
			}
			break
		}
		log.Info(message)
		p.hub.commands <- message
	}
}

func (p *Player) writeState() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		p.conn.Close()
	}()
	for {
		select {
		case message, ok := <-p.send:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				log.Error("Hub closed the player write state channel")
				p.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := p.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Error("writeState could not write state")
				return
			}
		case <-ticker.C:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Error("writeState could not write ping message to websocket")
				return
			}
		}
	}
}
