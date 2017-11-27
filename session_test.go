/*
	Digivance MVC Application Framework - Unit Tests
	Browser Session Collection Tests
	Dan Mayor (dmayor@digivance.com)

	This file defines the version 0.1.0 compatibility of session.go functions. These functions are written
	to demonstrate and test the intended use cases of the functions in session.go
*/

package mvcapp_test

import (
	"testing"

	"github.com/digivance/mvcapp"
)

func TestNewSession(t *testing.T) {
	session := mvcapp.NewSession()
	if session == nil {
		t.Error("Failed to create browser session collection")
	}
}

func TestSession_Get(t *testing.T) {
	session := mvcapp.NewSession()
	session.Set("Hello", "World")
	val := session.Get("Hello")
	if val != "World" {
		t.Error("Failed to validate session data")
	}
}

func TestSession_Set(t *testing.T) {
	session := mvcapp.NewSession()
	session.Set("Hello", "World")
	session.Set("More", "Worlds")
	session.Set("Hello", "World")
	val := session.Get("Hello")
	if val != "World" {
		t.Error("Failed to validate session data")
	}
}

func TestSession_Remove(t *testing.T) {
	session := mvcapp.NewSession()
	session.Set("Hello", "World")
	session.Remove("Hello")
	val := session.Get("Hello")
	if val != nil {
		t.Error("Failed to remove the zero world")
	}

	session.Set("First", "World")
	session.Set("Hello", "World")
	session.Set("Second", "World")
	session.Remove("Hello")
	val = session.Get("Hello")
	if val != nil {
		t.Error("Failed to remove the first world")
	}

	session.Set("Third", "Earth")
	session.Set("Hello", "World")
	session.Remove("Hello")
	val = session.Get("Hello")
	if val != nil {
		t.Error("Failed to remove after first world")
		t.Log(val)
	}

	if len(session.Values) != 3 {
		t.Error("Unexpected number of worlds remaining...")
		t.Log(len(session.Values))
	}
}
