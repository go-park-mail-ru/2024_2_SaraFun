package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_Sanitize(t *testing.T) {
	user := User{
		Username: "<script>alert('xss')</script>",
		Email:    "<img src='x' onerror='alert(1)'>",
		Password: "<b>password</b>",
	}
	user.Sanitize()

	assert.Equal(t, "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;", user.Username)
	assert.Equal(t, "&lt;img src=&#39;x&#39; onerror=&#39;alert(1)&#39;&gt;", user.Email)
	assert.Equal(t, "&lt;b&gt;password&lt;/b&gt;", user.Password)
}

func TestProfile_Sanitize(t *testing.T) {
	profile := Profile{
		FirstName: "<script>alert('xss')</script>",
		LastName:  "<b>Doe</b>",
		Gender:    "<img src='x' onerror='alert(1)'>",
		Target:    "<i>target</i>",
		About:     "<u>about me</u>",
	}
	profile.Sanitize()

	assert.Equal(t, "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;", profile.FirstName)
	assert.Equal(t, "&lt;b&gt;Doe&lt;/b&gt;", profile.LastName)
	assert.Equal(t, "&lt;img src=&#39;x&#39; onerror=&#39;alert(1)&#39;&gt;", profile.Gender)
	assert.Equal(t, "&lt;i&gt;target&lt;/i&gt;", profile.Target)
	assert.Equal(t, "&lt;u&gt;about me&lt;/u&gt;", profile.About)
}

func TestImage_Sanitize(t *testing.T) {
	image := Image{
		Id:   1,
		Link: "<img src='x' onerror='alert(1)'>",
	}
	image.Sanitize()

	assert.Equal(t, "&lt;img src=&#39;x&#39; onerror=&#39;alert(1)&#39;&gt;", image.Link)
}

func TestPersonCard_Sanitize(t *testing.T) {
	personCard := PersonCard{
		UserId:   1,
		Username: "<b>John</b>",
		Profile: Profile{
			FirstName: "<script>alert('xss')</script>",
			LastName:  "<b>Doe</b>",
		},
		Images: []Image{
			{Id: 1, Link: "<img src='x' onerror='alert(1)'>"},
		},
	}
	personCard.Sanitize()

	assert.Equal(t, "&lt;b&gt;John&lt;/b&gt;", personCard.Username)
}

func TestSession_Sanitize(t *testing.T) {
	session := Session{
		SessionID: "<script>alert('xss')</script>",
	}
	session.Sanitize()

	assert.Equal(t, "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;", session.SessionID)
}

func TestReport_Sanitize(t *testing.T) {
	report := Report{
		Body: "<script>alert('xss')</script>",
	}
	report.Sanitize()

	assert.Equal(t, "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;", report.Body)
}

func TestMessage_NoSanitize(t *testing.T) {
	message := Message{
		Author:   1,
		Receiver: 2,
		Body:     "No sanitization needed",
	}

	assert.Equal(t, "No sanitization needed", message.Body)
}
