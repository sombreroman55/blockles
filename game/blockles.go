package game

import (
	"encoding/json"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Command struct {
	Sender  string `json:"sender"`
	Command string `json:"command"`
}

type Blockles struct {
	Name         string
	maxPlayers   int
	players      map[uuid.UUID]*Player
	Options      GameOptions
	Id           uuid.UUID
	commands     chan []byte
	stateUpdates chan []byte
	register     chan *Player
	unregister   chan *Player
}

func (b *Blockles) AddNewPlayer(player *Player) {
	b.register <- player
	player.game = newGame(b.stateUpdates, player.Id)
	go player.readCommands()
	go player.writeState()
	go player.game.run()
}

func (b *Blockles) Run() {
	for {
		select {
		case player := <-b.register:
			log.Trace("Registering new player into the hub...")
			if len(b.players) < b.maxPlayers {
				b.players[player.Id] = player
			}
		case player := <-b.unregister:
			log.Trace("Unregistering player from the hub...")
			delete(b.players, player.Id)
		case commsg := <-b.commands:
			log.Infof("commsg: %s", string(commsg))
			var command Command
			err := json.Unmarshal(commsg, &command)
			if err != nil {
				log.Error("Error, could not deserialize command:", err)
				continue
			}

			senderId := uuid.MustParse(command.Sender)
			target, tok := b.players[senderId]
			if !tok {
				log.Error("Failed to get target for action, dropping")
				continue
			}
			target.game.fromBlockles <- command
		case message := <-b.stateUpdates:
			log.Trace("Sending update to players...")
			for _, player := range b.players {
				select {
				case player.send <- message:
				default:
					close(player.send)
					delete(b.players, player.Id)
				}
			}
		}
	}
}
