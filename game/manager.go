package game

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var soloGames map[uuid.UUID]*SoloGame

func InitGameManager() {
	soloGames = make(map[uuid.UUID]*SoloGame)
}

func NewSoloGame(p Player, o SoloGameOptions, id uuid.UUID) {
	log.Debug("Creating new game")
	game := &SoloGame{
		player:  p,
		options: o,
		id:      id,
	}
	soloGames[game.id] = game
}

func EndSoloGame(gameId uuid.UUID) {
	delete(soloGames, gameId)
}
