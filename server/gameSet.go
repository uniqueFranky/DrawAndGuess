package server

import (
	"errors"
	"google/uuid"
)

func (gs *GameSet) findGameByIdStr(idStr string) (*Game, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	return gs.findGameById(id)
}

func (gs *GameSet) findGameById(id uuid.UUID) (*Game, error) {
	for _, g := range gs.games {
		if g.Id == id {
			return g, nil
		}
	}
	return nil, errors.New("Game Not Found for Uuid: " + id.String())
}

func (gs *GameSet) appendGame(g *Game) error {
	gs.games = append(gs.games, g)
	return nil
}

func (gs *GameSet) deletePlayerInGame(u *User) {
	gameId := u.GameId
	g, err := gs.findGameById(gameId)
	if err != nil {
		return
	}
	_ = g.deletePlayerWithId(u.UserId)
}
