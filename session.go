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
	"strings"
	"time"
)

// SessionValue is a simple Key Value Pair struct used in the values
// collection of a Session (such as a per user session)
type SessionValue struct {
	// Key represents the key that serves as the index of this session value in a
	// browser session value collection
	Key string

	// Value is the data stored for this session value
	Value interface{}
}

// Session represents an http browser session data model
type Session struct {
	// ID is the unique key string that represents this browser session
	ID string

	// CreatedDate is the date and time when this browser session was created
	CreatedDate time.Time

	// ActivityDate is the date and time when this browser session was last active
	ActivityDate time.Time

	// Values is the collection of key value pair data stored in this browser session
	Values []*SessionValue
}

// NewSession returns a new Session model
func NewSession() *Session {
	return &Session{
		ID:           RandomString(32),
		CreatedDate:  time.Now(),
		ActivityDate: time.Now(),
		Values:       make([]*SessionValue, 0),
	}
}

// Get returns the interface{} of raw data value of the requested session value
func (session *Session) Get(key string) interface{} {
	for _, v := range session.Values {
		if strings.EqualFold(v.Key, key) {
			return v.Value
		}
	}

	return nil
}

// Set will overwrite or create a new value with the provided interface{} of raw data
func (session *Session) Set(key string, value interface{}) {
	for k, v := range session.Values {
		if strings.EqualFold(v.Key, key) {
			session.Values[k].Value = value
			return
		}
	}

	session.Values = append(session.Values, &SessionValue{Key: key, Value: value})
}

// Remove will remove the session value, identified by the provided key from this users
// session value collection
func (session *Session) Remove(key string) {
	for k, v := range session.Values {
		if strings.EqualFold(v.Key, key) {
			if k > 1 {
				session.Values = append(session.Values[:k], session.Values[k+1:]...)
				return
			}

			if k == 1 {
				session.Values = append(session.Values[2:], session.Values[0])
				return
			}

			if k == 0 {
				session.Values = session.Values[1:]
			}
		}
	}
}
