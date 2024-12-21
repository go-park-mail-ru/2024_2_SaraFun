package models

import (
	"html"
	"time"
)

//go:generate mockgen -source=*.go -destination=*_mock.go -package=*
//go:generate easyjson -all user.go

type User struct {
	ID       int    `json:"id" validate:"required"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
	Profile  int    `json:"profile" validate:"required"`
}

func (user *User) Sanitize() {
	user.Username = html.EscapeString(user.Username)
	user.Email = html.EscapeString(user.Email)
	user.Password = html.EscapeString(user.Password)
}

type Profile struct {
	ID           int    `json:"id" validate:"required"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	Age          int    `json:"age,omitempty"`
	BirthdayDate string `json:"birthday_date,omitempty"`
	Gender       string `json:"gender,omitempty"`
	Target       string `json:"target,omitempty"`
	About        string `json:"about,omitempty"`
}

func (profile *Profile) Sanitize() {
	profile.FirstName = html.EscapeString(profile.FirstName)
	profile.LastName = html.EscapeString(profile.LastName)
	profile.Gender = html.EscapeString(profile.Gender)
	profile.Target = html.EscapeString(profile.Target)
	profile.About = html.EscapeString(profile.About)
	profile.BirthdayDate = html.EscapeString(profile.BirthdayDate)
}

type Image struct {
	Id     int    `json:"id"`
	Link   string `json:"link"`
	Number int    `json:"number"`
}

func (image *Image) Sanitize() {
	image.Link = html.EscapeString(image.Link)
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

func (personCard *PersonCard) Sanitize() {
	personCard.Username = html.EscapeString(personCard.Username)
}

type Session struct {
	SessionID string    `json:"session_id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (session *Session) Sanitize() {
	session.SessionID = html.EscapeString(session.SessionID)
}

type Report struct {
	ID       int    `json:"id"`
	Author   int    `json:"author"`
	Receiver int    `json:"receiver"`
	Reason   string `json:"reason"`
	Body     string `json:"body"`
}

func (report *Report) Sanitize() {
	report.Body = html.EscapeString(report.Body)
	report.Reason = html.EscapeString(report.Reason)
}

type Message struct {
	ID       int    `json:"id"`
	Author   int    `json:"author"`
	Receiver int    `json:"receiver"`
	Body     string `json:"body"`
	Time     string `json:"time"`
}

func (message *Message) Sanitize() {
	message.Body = html.EscapeString(message.Body)
}

type Survey struct {
	ID       int    `json:"id"`
	Author   int    `json:"author"`
	Question string `json:"question"`
	Comment  string `json:"comment"`
	Rating   int    `json:"rating"`
	Grade    int    `json:"grade"`
}

type SurveyStat struct {
	Question string  `json:"question"`
	Grade    int     `json:"grade"`
	Rating   float32 `json:"rating"`
	Sum      int     `json:"sum"`
	Count    int     `json:"count"`
}

type AdminQuestion struct {
	Content string `json:"content"`
	Grade   int    `json:"grade"`
}

type Product struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageLink   string `json:"image_link"`
	Price       int    `json:"price"`
	Count       int    `json:"count"`
}

func (product *Product) Sanitize() {
	product.Description = html.EscapeString(product.Description)
	product.ImageLink = html.EscapeString(product.ImageLink)
	product.Title = html.EscapeString(product.Title)
}

type Award struct {
	DayNumber int    `json:"day_number"`
	Type      string `json:"type"`
	Count     int    `json:"count"`
}

func (award *Award) Sanitize() {
	award.Type = html.EscapeString(award.Type)
}

type Activity struct {
	Last_Login       string `json:"last_login"`
	Consecutive_days int    `json:"consecutive_days"`
	UserID           int    `json:"user_id"`
}

func (activity *Activity) Sanitize() {
	activity.Last_Login = html.EscapeString(activity.Last_Login)
}
