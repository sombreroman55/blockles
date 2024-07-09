// The main game state
package game

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type GameOptions struct {
	ShowGhostPiece bool
	ShowNextQueue  bool
	EnableHolding  bool
}

type Game struct {
	GameOver     bool
	fromBlockles chan Command
	toBlockles   chan []byte
	playerId     uuid.UUID
}

func newGame(blockCh chan []byte, playerId uuid.UUID) *Game {
	return &Game{
		GameOver:     false,
		fromBlockles: make(chan Command),
		toBlockles:   blockCh,
		playerId:     playerId,
	}
}

func (g *Game) run() {
	for !g.GameOver {
		select {
		case command := <-g.fromBlockles:
			log.Infof("Got %s from %s\n", command.Command, command.Sender)
			if command.Command == "quit" {
				g.GameOver = true
			}
		}
	}
}
