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

func NewGame(u *User, ans string) *Game {
	g := &Game{
		Id: uuid.New(),
		PlayerSet: UserSet{
			users: []*User{u},
		},
		DrawerName: u.UserName,
		Lines:      []Line{},
		Answer:     ans,
		Messages:   []Message{},
	}
	return g
}
