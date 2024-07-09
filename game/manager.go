package game

// TBD: This will probably become database tables, but right now in-memory is fine

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var games map[uuid.UUID]*Blockles
var players map[uuid.UUID]*Player

func InitGameManager() {
	games = make(map[uuid.UUID]*Blockles)
	players = make(map[uuid.UUID]*Player)
}

func PlayerExists(id uuid.UUID) bool {
	return players[id] != nil
}

func GetPlayer(id uuid.UUID) *Player {
	return players[id]
}

func DeletePlayer(id uuid.UUID) {
	delete(players, id)
}

func NewGame(name string, maxPlayers int, o GameOptions) uuid.UUID {
	log.Debug("Creating new solo game")
	game := &Blockles{
		Name:         name,
		maxPlayers:   maxPlayers,
		Options:      o,
		Id:           uuid.New(),
		players:      make(map[uuid.UUID]*Player),
		commands:     make(chan []byte, 256),
		stateUpdates: make(chan []byte, 256),
		register:     make(chan *Player, 2),
		unregister:   make(chan *Player, 2),
	}
	games[game.Id] = game
	return game.Id
}

func GameExists(gameId uuid.UUID) bool {
	return games[gameId] != nil
}

func GetGame(gameId uuid.UUID) *Blockles {
	return games[gameId]
}

func EndGame(gameId uuid.UUID) {
	delete(games, gameId)
}
