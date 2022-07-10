package server

import "google/uuid"

func (g *Game) appendPlayer(u *User) error {
	return g.PlayerSet.appendUser(u)
}

func (g *Game) deletePlayerWithId(id uuid.UUID) error {
	return g.PlayerSet.deleteUserById(id)
}

func (g *Game) appendLine(newLine Line) {
	g.Lines = append(g.Lines, newLine)
}

func NewGame(u *User, ans []byte) *Game {
	g := &Game{
		Id: uuid.New(),
		PlayerSet: UserSet{
			users: []*User{u},
		},
		Drawer: u,
		Lines:  []Line{},
		Answer: ans,
	}
	return g
}