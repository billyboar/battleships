package models

import (
	"github.com/gofrs/uuid"
)

// Session contains each board for computer and player
type Session struct {
	Player   *Board `json:"player"`
	Computer *Board `json:"computer"`
	ID       string `json:"id"`
}

// NewSession creates new session with boards
// initialized
func NewSession() (*Session, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	playerBoard, err := GenerateBoard(false)
	if err != nil {
		return nil, err
	}

	computerBoard, err := GenerateBoard(true)
	if err != nil {
		return nil, err
	}
	return &Session{
		Player:   playerBoard,
		Computer: computerBoard,
		ID:       id.String(),
	}, nil
}
