package game

// TBD: This will probably become database tables, but right now in-memory is fine

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var soloGames map[uuid.UUID]*SoloBlockles
var players map[uuid.UUID]*Player

func InitGameManager() {
	soloGames = make(map[uuid.UUID]*SoloBlockles)
	players = make(map[uuid.UUID]*Player)
}

func CreateNewPlayer(name string, conn *websocket.Conn) uuid.UUID {
	log.Debug("Creating new player")
	p := &Player{
		name: name,
		conn: conn,
		id: uuid.New(),
	}
	players[p.id] = p
	return p.id
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

func NewSoloGame(name string, o SoloBlocklesOptions) uuid.UUID {
	log.Debug("Creating new solo game")
	game := &SoloBlockles{
		Name:    name,
		Options: o,
		Id:      uuid.New(),
	}
	soloGames[game.Id] = game
	return game.Id
}

func SoloGameExists(gameId uuid.UUID) bool {
	return soloGames[gameId] != nil
}

func GetSoloGame(gameId uuid.UUID) *SoloBlockles {
	return soloGames[gameId]
}

func EndSoloGame(gameId uuid.UUID) {
	delete(soloGames, gameId)
}
