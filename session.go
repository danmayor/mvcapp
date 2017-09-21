/*
	Digivance MVC Application Framework
	HTTP Session Features
	Dan Mayor (dmayor@digivance.com)

	This file defines some HTTP Session functionality
*/

package mvcapp

import (
	"time"

	"github.com/Digivance/str"
)

// Session represents an http session data model
type Session struct {
	ID           string
	CreatedDate  time.Time
	ActivityDate time.Time
	Values       map[string]interface{}
}

// NewSession returns a new Session model
func NewSession() *Session {
	return &Session{
		ID:           str.Random(32),
		CreatedDate:  time.Now(),
		ActivityDate: time.Now(),
		Values:       make(map[string]interface{}, 0),
	}
}
