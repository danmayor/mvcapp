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

	"github.com/digivance/str"
)

// SessionValue is a simple Key Value Pair struct used in the values
// collection of a Session (such as a per user session)
type SessionValue struct {
	Key   string
	Value interface{}
}

// Session represents an http session data model
type Session struct {
	ID           string
	CreatedDate  time.Time
	ActivityDate time.Time
	Values       []*SessionValue
}

// NewSession returns a new Session model
func NewSession() *Session {
	return &Session{
		ID:           str.Random(32),
		CreatedDate:  time.Now(),
		ActivityDate: time.Now(),
		Values:       make([]*SessionValue, 0),
	}
}

// Get returns the interface{} of raw data value of the requested session value
func (session *Session) Get(key string) interface{} {
	for _, v := range session.Values {
		if str.Compare(v.Key, key) {
			return v.Value
		}
	}

	return nil
}

// Set will overwrite or create a new value with the provided interface{} of raw data
func (session *Session) Set(key string, value interface{}) {
	for k, v := range session.Values {
		if str.Compare(v.Key, key) {
			session.Values[k].Value = value
			return
		}
	}

	session.Values = append(session.Values, &SessionValue{Key: key, Value: value})
}
