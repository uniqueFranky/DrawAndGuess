package user

import "google/uuid"

type User struct {
	UserName string    `json:"userName"`
	UserId   uuid.UUID `json:"userId"`
	GameId   uuid.UUID `json:"-"`
}
