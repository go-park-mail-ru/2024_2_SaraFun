package models

import (
	"time"
)

//go:generate mockgen -source=*.go -destination=*_mock.go -package=*

//type User struct {
//	UUID        string    `json:"id" validate:"required"`
//	Name      string    `json:"name" validate:"required,min=2,max=50"`
//	Age       int       `json:"age" validate:"required,min=18,max=100"`
//	Gender    string    `json:"gender" validate:"required,oneof=male female"`
//	Email     string    `json:"email" validate:"required,regexp=^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"`
//	Phone     string    `json:"phone" validate:"required,regexp=^(\\+|[0-9])([0-9]*)$, min=11, max=12"`
//	Bio       string    `json:"bio,omitempty" validate:"max=150"`
//	Interests []string  `json:"interests,omitempty"`
//	Location  string    `json:"location,omitempty"`
//	CreatedAt time.Time `json:"created_at" validate:"required"`
//	UpdatedAt time.Time `json:"updated_at,omitempty"`
//}

type User struct {
	ID       int    `json:"id" validate:"required"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
	Profile  int    `json:"profile" validate:"required"`
}

type Profile struct {
	ID        int    `json:"id" validate:"required"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Age       int    `json:"age,omitempty"`
	Gender    string `json:"gender,omitempty"`
	Target    string `json:"target,omitempty"`
	About     string `json:"about,omitempty"`
}

type Image struct {
	Id   int    `json:"id"`
	Link string `json:"link"`
}

type Reaction struct {
	Id       int  `json:"id"`
	Author   int  `json:"author"`
	Receiver int  `json:"receiver"`
	Type     bool `json:"type"`
}

type PersonCard struct {
	UserId   int     `json:"user"`
	Username string  `json:"username"`
	Profile  Profile `json:"profile"`
	Images   []Image `json:"images"`
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
	SessionID string    `json:"session_id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
