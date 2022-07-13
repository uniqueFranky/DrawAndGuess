package game

import (
	"errors"
	"google/uuid"
)

type CurrentGameSet struct {
	Games []*CurrentGame
}

func newCurrentGameSet() *CurrentGameSet {
	return &CurrentGameSet{Games: []*CurrentGame{}}
}

func (gs *CurrentGameSet) findGameByIdStr(idStr string) (*CurrentGame, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	return gs.findGameById(id)
}

func (gs *CurrentGameSet) findGameById(id uuid.UUID) (*CurrentGame, error) {
	for _, g := range gs.Games {
		if g.Id == id {
			return g, nil
		}
	}
	return nil, errors.New("CurrentGame Not Found for Uuid: " + id.String())
}

func (gs *CurrentGameSet) appendGame(g *CurrentGame) error {
	gs.Games = append(gs.Games, g)
	return nil
}

func (gs *CurrentGameSet) endGame(g *CurrentGame) {
	g.end()
}

func (gs *CurrentGameSet) endGameById(id uuid.UUID) error {
	g, err := gs.findGameById(id)
	if err != nil {
		return err
	}
	g.end()
	return nil
}

func (gs *CurrentGameSet) endGameByIdStr(idStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	return gs.endGameById(id)
}

func (gs *CurrentGameSet) getGames() []*CurrentGame {
	return gs.Games
}

//func (gs *CurrentGameSet) deletePlayerInGame(u *User) {
//	gameId := u.GameId
//	g, err := gs.findGameById(gameId)
//	if err != nil {
//		return
//	}
//	_ = g.deletePlayerWithId(u.UserId)
//}
