package mvcapp

import (
	"time"

	"github.com/Digivance/str"
)

type SessionManager struct {
	Sessions []Session
}

func NewSessionManager() SessionManager {
	return SessionManager{
		Sessions: make([]Session, 0),
	}
}

func (manager *SessionManager) GetSession(id string) *Session {
	for key, val := range manager.Sessions {
		if str.Equals(val.ID, id) {
			return &manager.Sessions[key]
		}
	}

	return nil
}

func (manager *SessionManager) CreateSession() *Session {
	i := len(manager.Sessions)
	session := NewSession()
	manager.Sessions = append(manager.Sessions, session)
	return &manager.Sessions[i]
}

func (manager *SessionManager) SetSession(session Session) {
	id := session.ID
	res := manager.GetSession(id)

	if res != nil {
		res = &session
	} else {
		manager.Sessions = append(manager.Sessions, session)
	}
}

func (manager *SessionManager) DropSession(id string) {
	for key, val := range manager.Sessions {
		if str.Equals(val.ID, id) {
			manager.Sessions = append(manager.Sessions[:key-1], manager.Sessions[key:]...)
		}
	}
}

func (manager *SessionManager) CleanSessions() {
	now := time.Now()
	expires := now.Add(15 * time.Minute)

	for key, val := range manager.Sessions {
		if expires.After(val.ActivityDate) {
			manager.Sessions = append(manager.Sessions[:key-1], manager.Sessions[key:]...)
		}
	}
}
