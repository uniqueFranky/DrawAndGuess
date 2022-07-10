package server

import (
	"google/uuid"
	"gorilla/mux"
)

type MyBytes []byte

type User struct {
	UserName string    `json:"userName"`
	UserId   uuid.UUID `json:"userId"`
	GameId   uuid.UUID `json:"gameId"`
}

type UserSet struct {
	users []*User
}

type Game struct {
	Id        uuid.UUID `json:"id"`
	PlayerSet UserSet   `json:"playerSet"`
	Drawer    *User     `json:"drawer"`
	Lines     []Line    `json:"lines"`
	Answer    MyBytes   `json:"answer"`
}

type GameSet struct {
	games []*Game
}

type RelativePoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Line struct {
	Points []RelativePoint `json:"points"`
}

type GameCreationBatch struct {
	User   User   `json:"user"`
	Answer []byte `json:"answer"`
}

type Server struct {
	*mux.Router
	gameSet GameSet
	userSet UserSet
}
