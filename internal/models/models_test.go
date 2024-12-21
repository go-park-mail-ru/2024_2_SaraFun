package models

import (
	"html"
	"testing"
	"time"
)

func TestUser_Sanitize(t *testing.T) {
	user := &User{
		ID:       1,
		Username: "<script>alert('XSS')</script>",
		Email:    "<b>email@example.com</b>",
		Password: "\" onmouseover=alert('XSS')",
		Profile:  10,
	}

	user.Sanitize()

	if user.Username != html.EscapeString("<script>alert('XSS')</script>") {
		t.Errorf("Username not sanitized properly: got %v", user.Username)
	}
	if user.Email != html.EscapeString("<b>email@example.com</b>") {
		t.Errorf("Email not sanitized properly: got %v", user.Email)
	}
	if user.Password != html.EscapeString("\" onmouseover=alert('XSS')") {
		t.Errorf("Password not sanitized properly: got %v", user.Password)
	}
}

func TestProfile_Sanitize(t *testing.T) {
	profile := &Profile{
		ID:           2,
		FirstName:    "<first>name</first>",
		LastName:     "<last>name</last>",
		Age:          30,
		BirthdayDate: "<birthday>2000-01-01</birthday>",
		Gender:       "<gender>male</gender>",
		Target:       "<target>some target</target>",
		About:        "<about>something</about>",
	}

	profile.Sanitize()

	if profile.FirstName != html.EscapeString("<first>name</first>") {
		t.Errorf("FirstName not sanitized: %v", profile.FirstName)
	}
	if profile.LastName != html.EscapeString("<last>name</last>") {
		t.Errorf("LastName not sanitized: %v", profile.LastName)
	}
	if profile.BirthdayDate != html.EscapeString("<birthday>2000-01-01</birthday>") {
		t.Errorf("BirthdayDate not sanitized: %v", profile.BirthdayDate)
	}
	if profile.Gender != html.EscapeString("<gender>male</gender>") {
		t.Errorf("Gender not sanitized: %v", profile.Gender)
	}
	if profile.Target != html.EscapeString("<target>some target</target>") {
		t.Errorf("Target not sanitized: %v", profile.Target)
	}
	if profile.About != html.EscapeString("<about>something</about>") {
		t.Errorf("About not sanitized: %v", profile.About)
	}
}

func TestImage_Sanitize(t *testing.T) {
	img := &Image{
		Id:     1,
		Link:   "<img src=x onerror=alert('XSS')>",
		Number: 5,
	}
	img.Sanitize()

	if img.Link != html.EscapeString("<img src=x onerror=alert('XSS')>") {
		t.Errorf("Link not sanitized: %v", img.Link)
	}
}

func TestPersonCard_Sanitize(t *testing.T) {
	card := &PersonCard{
		UserId:   10,
		Username: "<u>user</u>",
		Profile:  Profile{ID: 1, FirstName: "John", LastName: "Doe"},
		Images:   []Image{{Id: 1, Link: "<img>", Number: 1}},
	}
	card.Sanitize()

	if card.Username != html.EscapeString("<u>user</u>") {
		t.Errorf("Username not sanitized: %v", card.Username)
	}
}

func TestSession_Sanitize(t *testing.T) {
	session := &Session{
		SessionID: "<session>",
		UserID:    123,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour),
	}
	session.Sanitize()

	if session.SessionID != html.EscapeString("<session>") {
		t.Errorf("SessionID not sanitized: %v", session.SessionID)
	}
}

func TestReport_Sanitize(t *testing.T) {
	report := &Report{
		ID:       1,
		Author:   2,
		Receiver: 3,
		Reason:   "<reason>hack</reason>",
		Body:     "<body>text</body>",
	}
	report.Sanitize()

	if report.Reason != html.EscapeString("<reason>hack</reason>") {
		t.Errorf("Reason not sanitized: %v", report.Reason)
	}
	if report.Body != html.EscapeString("<body>text</body>") {
		t.Errorf("Body not sanitized: %v", report.Body)
	}
}

func TestMessage_Sanitize(t *testing.T) {
	msg := &Message{
		ID:       1,
		Author:   10,
		Receiver: 20,
		Body:     "<script>alert('msg')</script>",
		Time:     "2024-12-17T20:29:10Z",
	}
	msg.Sanitize()

	if msg.Body != html.EscapeString("<script>alert('msg')</script>") {
		t.Errorf("Body not sanitized: %v", msg.Body)
	}
}

func TestProduct_Sanitize(t *testing.T) {
	product := &Product{
		Title:       "<title>prod</title>",
		Description: "<desc>description</desc>",
		ImageLink:   "<img>image</img>",
		Price:       100,
	}
	product.Sanitize()

	if product.Title != html.EscapeString("<title>prod</title>") {
		t.Errorf("Title not sanitized: %v", product.Title)
	}
	if product.Description != html.EscapeString("<desc>description</desc>") {
		t.Errorf("Description not sanitized: %v", product.Description)
	}
	if product.ImageLink != html.EscapeString("<img>image</img>") {
		t.Errorf("ImageLink not sanitized: %v", product.ImageLink)
	}
}

func TestSurvey_Initialize(t *testing.T) {
	survey := Survey{
		ID:       1,
		Author:   2,
		Question: "What is your name?",
		Comment:  "Nice survey",
		Rating:   5,
		Grade:    2,
	}
	// Здесь нет Sanitize(), просто проверим, что поля корректно присвоены
	if survey.ID != 1 || survey.Author != 2 || survey.Question != "What is your name?" ||
		survey.Comment != "Nice survey" || survey.Rating != 5 || survey.Grade != 2 {
		t.Errorf("Survey fields not set correctly: %+v", survey)
	}
}

func TestSurveyStat_Initialize(t *testing.T) {
	stat := SurveyStat{
		Question: "Q?",
		Grade:    3,
		Rating:   4.5,
		Sum:      9,
		Count:    2,
	}
	// Проверим базовую инициализацию
	if stat.Question != "Q?" || stat.Grade != 3 || stat.Rating != 4.5 ||
		stat.Sum != 9 || stat.Count != 2 {
		t.Errorf("SurveyStat fields not set correctly: %+v", stat)
	}
}

func TestAdminQuestion_Initialize(t *testing.T) {
	aq := AdminQuestion{
		Content: "Q content",
		Grade:   1,
	}
	if aq.Content != "Q content" || aq.Grade != 1 {
		t.Errorf("AdminQuestion fields not set correctly: %+v", aq)
	}
}

func TestReaction_Initialize(t *testing.T) {
	reaction := Reaction{
		Id:       10,
		Author:   20,
		Receiver: 30,
		Type:     true,
	}
	if reaction.Id != 10 || reaction.Author != 20 || reaction.Receiver != 30 || reaction.Type != true {
		t.Errorf("Reaction fields not set correctly: %+v", reaction)
	}
}
