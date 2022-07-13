package game

import "google/uuid"

type EndedGame struct {
	Id         uuid.UUID `json:"id"`
	Answer     string    `json:"answer"`
	WinnerName string    `json:"winnerName"`
}
