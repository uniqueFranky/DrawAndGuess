package game

import (
	"errors"
	"google/uuid"
)

type EndedGameSet struct {
	games []*EndedGame
}

func newEndedGameSet() *EndedGameSet {
	return &EndedGameSet{games: []*EndedGame{}}
}

func (gs *EndedGameSet) appendEndedGame(g *CurrentGame, winner string) {
	eg := &EndedGame{
		Id:         g.Id,
		Answer:     g.Answer,
		WinnerName: winner,
	}
	gs.games = append(gs.games, eg)
}

func (gs *EndedGameSet) getEndedGames() []*EndedGame {
	return gs.games
}

func (gs *EndedGameSet) findGameById(id uuid.UUID) (*EndedGame, error) {
	for _, g := range gs.games {
		if g.Id == id {
			return g, nil
		}
	}
	return nil, errors.New("game with given uuid not found")
}

func (gs *EndedGameSet) findGameByIdStr(idStr string) (*EndedGame, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	return gs.findGameById(id)
}
