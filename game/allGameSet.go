package game

import (
	"errors"
	"google/uuid"
)

type AllGameSet struct {
	currentGameSet *CurrentGameSet
	endedGameSet   *EndedGameSet
}

func NewAllGameSet() AllGameSet {
	return AllGameSet{
		currentGameSet: newCurrentGameSet(),
		endedGameSet:   newEndedGameSet(),
	}
}

//Functions of CurrentGameSet

func (ags *AllGameSet) FindCurrentGameById(id uuid.UUID) (*CurrentGame, error) {
	g, err := ags.currentGameSet.findGameById(id)
	if err != nil {
		return nil, err
	}
	if g.IsEnd {
		return nil, errors.New("the game has ended")
	}

	return g, nil
}

func (ags *AllGameSet) FindCurrentGameByIdStr(idStr string) (*CurrentGame, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	return ags.FindCurrentGameById(id)
}

func (ags *AllGameSet) FindGameInCurrentGamesById(id uuid.UUID) (*CurrentGame, error) {
	return ags.currentGameSet.findGameById(id)
}

func (ags *AllGameSet) FindGameInCurrentGamesByIdStr(idStr string) (*CurrentGame, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	return ags.FindGameInCurrentGamesById(id)
}

func (ags *AllGameSet) AppendCurrentGame(g *CurrentGame) error {
	return ags.currentGameSet.appendGame(g)
}

func (ags *AllGameSet) EndCurrentGame(g *CurrentGame) {
	ags.currentGameSet.endGame(g)
}

func (ags *AllGameSet) GetCurrentGames() []*CurrentGame {
	return ags.currentGameSet.getGames()
}

//Functions of EndedGameSet

func (ags *AllGameSet) AppendEndedGame(g *CurrentGame, winner string) {
	ags.endedGameSet.appendEndedGame(g, winner)
}

func (ags *AllGameSet) GetEndedGames() []*EndedGame {
	return ags.endedGameSet.getEndedGames()
}

func (ags *AllGameSet) FindEndedGameByIdStr(idStr string) (*EndedGame, error) {
	return ags.endedGameSet.findGameByIdStr(idStr)
}

//Functions of AllGameSet

func (ags *AllGameSet) EndGame(g *CurrentGame, winner string) {
	ags.EndCurrentGame(g)
	ags.AppendEndedGame(g, winner)
}
