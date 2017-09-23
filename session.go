/*
	Digivance MVC Application Framework
	Browser Session Features
	Dan Mayor (dmayor@digivance.com)

	This file defines the basic Browser Session object. This object provides a server side memory map
	that can store values between requests. This session is identified by the session ID provided in
	the cookies of the incoming request.

	This package is released under as open source under the LGPL-3.0 which can be found:
	https://opensource.org/licenses/LGPL-3.0
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
