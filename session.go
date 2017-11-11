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

func (session *Session) Get(key string) interface{} {
	for _, v := range session.Values {
		if str.Compare(v.Key, key) {
			return v.Value
		}
	}

	return nil
}

func (session *Session) Set(key string, value interface{}) {
	for k, v := range session.Values {
		if str.Compare(v.Key, key) {
			session.Values[k].Value = value
			return
		}
	}

	session.Values = append(session.Values, &SessionValue{Key: key, Value: value})
}
