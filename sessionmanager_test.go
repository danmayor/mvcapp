package mvcapp

import (
	"testing"
	"time"

	"github.com/Digivance/str"
)

func TestSessionManager(t *testing.T) {
	mgr := NewSessionManager()

	// Tests that we create a new browser session
	createdSession := mgr.CreateSession(str.Random(32))
	if createdSession == nil {
		t.Error("Failed to test session manager: No session created")
	}

	// Tests that we can receive the session from our manager
	returnedSession := mgr.GetSession(createdSession.ID)
	if returnedSession == nil {
		t.Error("Failed to test session manager: No session returned")
	}

	// We use this to test changing session members
	returnedSession.ActivityDate = time.Now().Add(-30 * time.Minute)

	// We use this to test that setting the returned session's member affected
	// the session in our manager
	comparedSession := mgr.GetSession(createdSession.ID)
	if !comparedSession.ActivityDate.Equal(returnedSession.ActivityDate) {
		t.Error("Failed to test session manager: Can not compare saved value")
	}

	// Testing that we can create another new session
	newSession := NewSession()
	newSession.ID = "TEST123456789"

	// Saving this session to the
	mgr.SetSession(newSession)

	newReturnedSession := mgr.GetSession(newSession.ID)
	if newReturnedSession == nil {
		t.Error("Failed to test session manager: Can not get set session value")
	}

	newReturnedSession.Values["TestKey"] = "TestVal"
	testVal := newReturnedSession.Values["TestKey"].(string)
	if !str.Compare(testVal, "TestVal") {
		t.Error("Failed to test setting session values")
	}

	// Testing drop method (and clean method indirectly), we add a duplicate browser
	// session, and attempt to drop by the ID. Expected behavior is to delete both
	// browser sessions.
	otherSession := &Session{
		ID: createdSession.ID,
	}

	mgr.Sessions = append(mgr.Sessions, otherSession)
	if len(mgr.Sessions) < 3 {
		t.Error("Failed to append duplicate session for testing :(")
	}

	mgr.DropSession(createdSession.ID)

	if len(mgr.Sessions) > 1 {
		t.Error("Failed to drop browser sessions :(")
	}
}
