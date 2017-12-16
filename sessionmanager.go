/*
	Digivance MVC Application Framework
	Session Manager Features
	Dan Mayor (dmayor@digivance.com)

	This file defines functionality for an in process browser session manager system. (E.g. per user
	server side memory map)
*/

package mvcapp

import (
	"time"
)

// SessionManager is the base struct that manages the collection
// of current http session models.
type SessionManager struct {
	// SessionIDKey is the name of the cookie value that will store the unique ID of the browser
	// session
	SessionIDKey string

	// Sessions is the collection of browser session objects
	Sessions map[string]*Session

	// SessionTimeout is the duration of time that a browser session will stay in memory between
	// requests / activity from the user
	SessionTimeout time.Duration
}

// NewSessionManager returns a new Session Manager object
func NewSessionManager() *SessionManager {
	return &SessionManager{
		Sessions:       make(map[string]*Session, 0),
		SessionTimeout: (15 * time.Minute),
	}
}

// GetSession returns the current http session for the provided session id
func (manager *SessionManager) GetSession(id string) *Session {
	return manager.Sessions[id]
}

// Contains detects if the requested id (key) exists in this session collection
func (manager *SessionManager) Contains(id string) bool {
	if session := manager.GetSession(id); session != nil {
		if session.ID == id {
			return true
		}
	}

	return false
}

// CreateSession creates and returns a new http session model
func (manager *SessionManager) CreateSession(id string) *Session {
	session := NewSession()
	session.ID = id
	manager.Sessions[id] = session
	return manager.Sessions[id]
}

// SetSession will set (creating if necessary) the provided session to
// the session manager collection
func (manager *SessionManager) SetSession(session *Session) {
	manager.Sessions[session.ID] = session
}

// DropSession will remove a session from the session manager collection based
// on the provided session id
func (manager *SessionManager) DropSession(id string) {
	for key, val := range manager.Sessions {
		if val.ID == id {
			delete(manager.Sessions, key)
		}
	}
}

// CleanSessions will drop inactive sessions
func (manager *SessionManager) CleanSessions() {
	expired := time.Now().Add(-manager.SessionTimeout)

	for key, val := range manager.Sessions {
		if val.ActivityDate.Before(expired) {
			delete(manager.Sessions, key)
		}
	}
}
