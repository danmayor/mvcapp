package mvcapp

import (
	"time"

	"github.com/Digivance/str"
)

type Session struct {
	ID           string
	CreatedDate  time.Time
	ActivityDate time.Time
}

func NewSession() Session {
	return Session{
		ID:           str.Random(32),
		CreatedDate:  time.Now(),
		ActivityDate: time.Now(),
	}
}
