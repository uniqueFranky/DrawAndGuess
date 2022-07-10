package server

import "google/uuid"

func (g *Game) appendPlayer(u *User) error {
	return g.PlayerSet.appendUser(u)
}

func (g *Game) deletePlayerWithId(id uuid.UUID) error {
	err := g.PlayerSet.deleteUserById(id)
	if err != nil {
		return err
	}
	if len(g.PlayerSet.users) == 0 {
		g.HasEnded = true
	}
	return nil
}

func (g *Game) appendLine(newLine Line) {
	g.Lines = append(g.Lines, newLine)
}

func NewGame(u *User, ans string) *Game {
	g := &Game{
		Id: uuid.New(),
		PlayerSet: UserSet{
			users: []*User{u},
		},
		Drawer:   u,
		Lines:    []Line{},
		Answer:   ans,
		Messages: []Message{},
		HasEnded: false,
	}
	return g
}
