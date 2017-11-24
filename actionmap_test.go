/*
	Digivance MVC Application Framework - Unit Tests
	Action Map Feature Tests
	Dan Mayor (dmayor@digivance.com)

	This file defines the version 0.1.0 compatibility of actionmap.go functions. These functions are written
	to demonstrate and test the intended use cases of the functions in actionmap.go
*/

package mvcapp_test

import (
	"strings"
	"testing"

	"github.com/digivance/mvcapp"
)

// actionHandler is an internal helper used to test the execution of an Action Map method
func actionHandler(params []string) *mvcapp.ActionResult {
	return mvcapp.NewActionResult([]byte("Version 0.1.0 Compliant"))
}

// validateResultData is an internal helper used to test that the execution of the actionHandler
// returned the expected payload data
func validateResultData(actionResult *mvcapp.ActionResult) bool {
	if string(actionResult.Data) == "Version 0.1.0 Compliant" {
		return true
	}

	return false
}

func TestNewActionMap(t *testing.T) {
	// Create a new action map
	actionMap := mvcapp.NewActionMap("", "Test", actionHandler)
	if actionMap == nil {
		t.Fatal("Failed to create action map")
	}

	// Execute the actionHandler we registered above
	actionResult := actionMap.Method([]string{})
	if actionResult == nil {
		t.Fatal("Failed to execute actionMap method")
	}

	// Ensure that the results returned are as expected
	if !validateResultData(actionResult) {
		t.Error("Failed to validate result data")
	}
}

func TestNewGetActionMap(t *testing.T) {
	// Create a new GET ActionMap
	actionMap := mvcapp.NewGetActionMap("Test", actionHandler)
	if actionMap == nil {
		t.Fatal("Failed to create GET action map")
	}

	// Ensure it is a GET ActionMap
	if !strings.EqualFold(actionMap.Verb, "GET") {
		t.Error("Failed to set GET HTTP Verb")
	}

	// Execute the actionHandler we registered above
	actionResult := actionMap.Method([]string{})
	if actionResult == nil {
		t.Fatal("Failed to execute actionMap method")
	}

	// Ensure that the results returned are as expected
	if !validateResultData(actionResult) {
		t.Error("Failed to validate result data")
	}
}

func TestNewPostActionMap(t *testing.T) {
	// Create a new POST ActionMap
	actionMap := mvcapp.NewPostActionMap("Test", actionHandler)
	if actionMap == nil {
		t.Fatal("Failed to create POST action map")
	}

	// Ensure it is a POST ActionMap
	if !strings.EqualFold(actionMap.Verb, "POST") {
		t.Error("Failed to set POST Verb")
	}

	// Execute the actionHandler we registered above
	actionResult := actionMap.Method([]string{})
	if actionResult == nil {
		t.Fatal("Failed to execute actionMap method")
	}

	// Ensure that the results returned are as expected
	if !validateResultData(actionResult) {
		t.Error("Failed to validate result data")
	}
}

func TestNewPutActionMap(t *testing.T) {
	// Create a new PUT ActionMap
	actionMap := mvcapp.NewPutActionMap("Test", actionHandler)
	if actionMap == nil {
		t.Fatal("Failed to create PUT action map")
	}

	// Ensure if is a PUT ActionMap
	if !strings.EqualFold(actionMap.Verb, "PUT") {
		t.Error("Failed set PUT Verb")
	}

	// Execute the actionHandler we registered above
	actionResult := actionMap.Method([]string{})
	if actionResult == nil {
		t.Fatal("Failed to execute actionMap method")
	}

	// Ensure that the results returned are as expected
	if !validateResultData(actionResult) {
		t.Error("Failed to validate result data")
	}
}

func TestNewDeleteActionMap(t *testing.T) {
	// Create a new DELETE ActionMap
	actionMap := mvcapp.NewDeleteActionMap("Test", actionHandler)
	if actionMap == nil {
		t.Fatal("Failed to create DELETE action map")
	}

	// Ensure if is a DELETE ActionMap
	if !strings.EqualFold(actionMap.Verb, "DELETE") {
		t.Error("Failed set DELETE Verb")
	}

	// Execute the actionHandler we registered above
	actionResult := actionMap.Method([]string{})
	if actionResult == nil {
		t.Fatal("Failed to execute actionMap method")
	}

	// Ensure that the results returned are as expected
	if !validateResultData(actionResult) {
		t.Error("Failed to validate result data")
	}
}
