package game

import (
	"github.com/google/uuid"
)

type SoloGameOptions struct {
	ShowGhostPiece bool
	ShowNextQueue  bool
	EnableHolding  bool
}

type SoloGame struct {
	name    string
	player  Player
	options SoloGameOptions
	id      uuid.UUID
}


