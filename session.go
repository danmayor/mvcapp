/*
	Digivance MVC Application Framework
	Browser Session Features
	Dan Mayor (dmayor@digivance.com)

	This file defines the basic Browser Session object. This object provides a server side memory map
	that can store values between requests. This session is identified by the session ID provided in
	the cookies of the incoming request.
*/

package mvcapp

import (
	"time"
)

// Session represents an http browser session data model
type Session struct {
	// ID is the unique key string that represents this browser session
	ID string

	// CreatedDate is the date and time when this browser session was created
	CreatedDate time.Time

	// ActivityDate is the date and time when this browser session was last active
	ActivityDate time.Time

	// Values is the collection of key value pair data stored in this browser session
	Values map[string]interface{}
}

// NewSession returns a new Session model
func NewSession() *Session {
	return &Session{
		ID:           RandomString(32),
		CreatedDate:  time.Now(),
		ActivityDate: time.Now(),
		Values:       make(map[string]interface{}, 0),
	}
}

// Get returns the interface{} of raw data value of the requested session value
func (session *Session) Get(key string) interface{} {
	return session.Values[key]
}

// Set will overwrite or create a new value with the provided interface{} of raw data
func (session *Session) Set(key string, value interface{}) {
	session.Values[key] = value
}

// Remove will remove the session value, identified by the provided key from this users
// session value collection
func (session *Session) Remove(key string) {
	delete(session.Values, key)
}
