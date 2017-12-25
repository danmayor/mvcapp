/*
	Digivance MVC Application Framework - Unit Tests
	Session Manager Tests
	Dan Mayor (dmayor@digivance.com)

	This file defines the version 0.2.0 compatibility of sessionmanager.go functions. These functions are written
	to demonstrate and test the intended use cases of the functions in sessionmanager.go
*/

package mvcapp_test

import (
	"testing"
	"time"

	"github.com/digivance/mvcapp"
)

// TestNewSessionManager ensures that mvcapp.NewSessionManager returns the expected value
func TestNewSessionManager(t *testing.T) {
	manager := mvcapp.NewSessionManager()
	if manager == nil {
		t.Fatal("Failed to create session manager")
	}
}

// TestSessionManager_GetSession ensures that SessionManager.GetSession returns the expected value
func TestSessionManager_GetSession(t *testing.T) {
	session := mvcapp.NewSession()
	session.Set("Test", "Val")

	manager := mvcapp.NewSessionManager()
	manager.SetSession(session)

	testSession := manager.GetSession(session.ID)
	if testSession.ID != session.ID {
		t.Fatal("Failed to get browser session")
	}

	if testSession.Get("Test") != session.Get("Test") {
		t.Error("Failed to validate that the session key value pairs pass through")
		t.Log(testSession.Get("Test"))
	}
}

// TestSessionManager_Contains ensures that SessionMnager.Contains returns the expected value
func TestSessionManager_Contains(t *testing.T) {
	session := mvcapp.NewSession()
	manager := mvcapp.NewSessionManager()
	manager.SetSession(session)

	if !manager.Contains(session.ID) {
		t.Error("Failed to verify that session manager contains the test browser session")
	}

	if manager.Contains("fubar") {
		t.Error("Failed to identify missing browser session")
	}
}

// TestSessionManager.CreateSession ensures that SessionManager.CreateSession operates as expected
func TestSessionManager_CreateSession(t *testing.T) {
	manager := mvcapp.NewSessionManager()
	session := manager.CreateSession("TestID")
	if !manager.Contains(session.ID) {
		t.Error("Failed to create new browser session from session manager")
	}
}

// TestSessionManager.SetSession ensures that SessionManager.SetSession operates as expected
func TestSessionManager_SetSession(t *testing.T) {
	manager := mvcapp.NewSessionManager()
	session := manager.CreateSession("TestID")

	newSession := mvcapp.NewSession()
	newSession.ID = session.ID
	newSession.Set("Test", "Value")

	manager.SetSession(newSession)
	if len(newSession.Values) <= 0 {
		t.Error("Failed to set newSession values")
	}
}

// TestSessionManager_DropSession ensures that SessionManager.DropSession operates as expected
func TestSessionManager_DropSession(t *testing.T) {
	manager := mvcapp.NewSessionManager()
	manager.CreateSession("Deletable")
	manager.DropSession("Deletable")

	if manager.Contains("Deletable") {
		t.Error("Failed to drop zero browser session")
	}

	manager.CreateSession("Deletable")
	manager.CreateSession("First")
	manager.DropSession("Deletable")

	if manager.Contains("Deletable") {
		t.Error("Failed to drop first browser session")
	}

	manager.CreateSession("Deletable")
	manager.CreateSession("Second")
	manager.DropSession("Deletable")

	if manager.Contains("Deletable") {
		t.Error("Failed to drop second browser session")
	}

	manager.CreateSession("Deletable")
	manager.CreateSession("Third")
	manager.DropSession("Deletable")

	if manager.Contains("Deletable") {
		t.Error("Failed to drop third browser session")
	}

	if !manager.Contains("First") || !manager.Contains("Second") || !manager.Contains("Third") {
		t.Error("Failed to retain the proper browser sessions")
	}
}

// TestSessionManager_CleanSessions ensures that SessionManager.CleanSessions operates as expected
func TestSessionManager_CleanSessions(t *testing.T) {
	manager := mvcapp.NewSessionManager()
	session := manager.CreateSession("Deletable")
	session.ActivityDate = time.Now().Add(-30 * time.Minute)
	manager.CreateSession("A")
	manager.CleanSessions()
	if manager.Contains("Deletable") {
		t.Error("Failed to delete zero expired session")
	}

	session = manager.CreateSession("Deletable")
	session.ActivityDate = time.Now().Add(-30 * time.Minute)
	manager.CreateSession("B")
	manager.CleanSessions()
	if manager.Contains("Deletable") {
		t.Error("Failed to delete first expired session")
	}

	session = manager.CreateSession("Deletable")
	session.ActivityDate = time.Now().Add(-30 * time.Minute)
	manager.CreateSession("C")
	manager.CleanSessions()
	if manager.Contains("Deletable") {
		t.Error("Failed to delete second expired session")
	}

	session = manager.CreateSession("Deletable")
	session.ActivityDate = time.Now().Add(-30 * time.Minute)
	manager.CreateSession("D")
	manager.CleanSessions()
	if manager.Contains("Deletable") {
		t.Error("Failed to delete third expired session")
	}
}
