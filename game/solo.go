package game

import (
	"errors"
	"github.com/google/uuid"
)

type SoloBlocklesOptions struct {
	GameOpts GameOptions
}

type SoloBlockles struct {
	Name    string
	player  *Player
	Options SoloBlocklesOptions
	Id      uuid.UUID
}

func (s* SoloBlockles) AddPlayer(p *Player) error {
	if s.player != nil {
		return errors.New("Player already exists")
	}
	s.player = p
	return nil
}
