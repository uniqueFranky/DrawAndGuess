package game

import "DrawAndGuess/user"

type Message struct {
	From    string `json:"from"`
	Content string `json:"content"`
}

type MessageWithUser struct {
	From    *user.User `json:"from"`
	Content string     `json:"content"`
}
