package models

import (
	"time"
)

//go:generate mockgen -source=*.go -destination=*_mock.go -package=*

type User struct {
	ID        string    `json:"id" validate:"required"`
	Name      string    `json:"name" validate:"required,min=2,max=50"`
	Age       int       `json:"age" validate:"required,min=18,max=100"`
	Gender    string    `json:"gender" validate:"required,oneof=male female"`
	Email     string    `json:"email" validate:"required,regexp=^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"`
	Phone     string    `json:"phone" validate:"required,regexp=^(\\+|[0-9])([0-9]*)$, min=11, max=12"`
	Bio       string    `json:"bio,omitempty" validate:"max=150"`
	Interests []string  `json:"interests,omitempty"`
	Location  string    `json:"location,omitempty"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type Match struct {
	UserID1   string    `json:"user_id_1" validate:"required"`
	UserID2   string    `json:"user_id_2" validate:"required"`
	MatchedAt time.Time `json:"matched_at" validate:"required"`
}

type Message struct {
	ID       string    `json:"id" validate:"required"`
	MatchID  string    `json:"match_id" validate:"required"`
	SenderID string    `json:"sender_id" validate:"required"`
	Content  string    `json:"content" validate:"required,min=1,max=500"`
	SentAt   time.Time `json:"sent_at" validate:"required"`
}

type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
