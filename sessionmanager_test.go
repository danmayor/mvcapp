package mvcapp

import (
	"testing"
	"time"

	"github.com/Digivance/str"
)

func TestSessionManager(t *testing.T) {
	mgr := NewSessionManager()

	createdSession := mgr.CreateSession(str.Random(32))
	if createdSession == nil {
		t.Error("Failed to test session manager: No session created")
	}

	returnedSession := mgr.GetSession(createdSession.ID)
	if returnedSession == nil {
		t.Error("Failed to test session manager: No session returned")
	}

	returnedSession.ActivityDate = time.Now().Add(-30 * time.Minute)

	comparedSession := mgr.GetSession(createdSession.ID)
	if !comparedSession.ActivityDate.Equal(returnedSession.ActivityDate) {
		t.Error("Failed to test session manager: Can not compare saved value")
	}

	newSession := NewSession()
	newSession.ID = "TEST123456789"

	mgr.SetSession(newSession)

	newReturnedSession := mgr.GetSession(newSession.ID)
	if newReturnedSession == nil {
		t.Error("Failed to test session manager: Can not get set session value")
	}
}
