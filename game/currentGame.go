package game

import (
	"DrawAndGuess/user"
	"errors"
	"google/uuid"
)

type CurrentGame struct {
	Id         uuid.UUID    `json:"id"`
	PlayerSet  user.UserSet `json:"-"`
	DrawerName string       `json:"drawerName"`
	Lines      []Line       `json:"lines"`
	Answer     string       `json:"answer"`
	Messages   []Message    `json:"messages"`
	IsEnd      bool         `json:"IsEnd"`
}

func (g *CurrentGame) AppendPlayer(u *user.User) error {
	return g.PlayerSet.AppendUser(u)
}

func (g *CurrentGame) DeletePlayerWithId(id uuid.UUID) error {
	return g.PlayerSet.DeleteUserById(id)
}

func (g *CurrentGame) AppendLine(newLine Line) {
	g.Lines = append(g.Lines, newLine)
}

func NewGame(u *user.User, ans string) *CurrentGame {
	g := &CurrentGame{
		Id: uuid.New(),
		PlayerSet: user.UserSet{
			Users: []*user.User{u},
		},
		DrawerName: u.UserName,
		Lines:      []Line{},
		Answer:     ans,
		Messages:   []Message{},
		IsEnd:      false,
	}
	return g
}

func (g *CurrentGame) end() {
	g.IsEnd = true
}

func (g *CurrentGame) HasEnded() error {
	if g.IsEnd {
		return errors.New("the games has ended")
	} else {
		return nil
	}
}
