package game

import "DrawAndGuess/user"

type RelativePoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Line struct {
	Points []RelativePoint `json:"points"`
}

type LineWithUser struct {
	NewLine Line      `json:"newLine"`
	From    user.User `json:"from"`
}

type LinesWithUser struct {
	NewLines []Line    `json:"newLines"`
	From     user.User `json:"from"`
}
