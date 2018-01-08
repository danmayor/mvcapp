/*
	Digivance MVC Application Framework - Unit Tests
	Action Map Feature Tests
	Dan Mayor (dmayor@digivance.com)

	This file defines the version 0.2.0 compatibility of actionresult.go functions. These functions are written
	to demonstrate and test the intended use cases of the functions in actionresult.go
*/

package mvcapp_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/digivance/mvcapp"
)

// TestNewActionResult ensures that mvcapp.NewActionResult returns the expected value
func TestNewActionResult(t *testing.T) {
	// Create a new generic action result
	actionResult := mvcapp.NewActionResult([]byte("Version 0.2.0 Compliant"))
	if actionResult == nil {
		t.Fatal("Failed to create new action result")
	}

	// Confirm that the payload data was set correctly
	if string(actionResult.Data) != "Version 0.2.0 Compliant" {
		t.Error("Failed to validate result data")
	}
}

// TestNewViewResult ensures that mvcapp.NewViewResult returns the expected value
func TestNewViewResult(t *testing.T) {
	// Create a temporary template file and set the expected resulting value
	filename := fmt.Sprintf("%s/%s", mvcapp.GetApplicationPath(), "_test_template.htm")
	templateData := "{{ define \"mvcapp\" }}<html><head><title>Test</title></head><body>Testing</body></html>{{ end }}"
	expectedResultData := "<html><head><title>Test</title></head><body>Testing</body></html>"
	defer os.RemoveAll(filename)

	err := ioutil.WriteFile(filename, []byte(templateData), 0644)
	if err != nil {
		t.Error(err)
	}

	// Construct view result from temporary template file
	viewResult, err := mvcapp.NewViewResult([]string{filename}, nil)
	if err != nil {
		t.Fatalf("Failed to create view result: %s", err)
	}

	// Validate the resulting result data
	if string(viewResult.Data) != expectedResultData {
		t.Error("Failed to validate view result data")
	}
}

// TestNewJSONResult ensures that mvcapp.NewJSONResult returns the expected value
func TestNewJSONResult(t *testing.T) {
	// Create a json encoded payload
	payload := "Version 0.1.0 Compliant"
	jsonResult, err := mvcapp.NewJSONResult(payload)
	if err != nil {
		t.Fatal("Failed to create JSON result")
	}

	// Deserialize the created json byte array
	var res string
	err = json.Unmarshal(jsonResult.Data, &res)
	if err != nil {
		t.Fatal(err)
	}

	// Test that the returned value is the intended payload
	if res != payload {
		t.Error("Failed to validate payload")
	}
}

// TestActionResult_AddHeader ensures that ActionResult.AddHeader operates as expected
func TestActionResult_AddHeader(t *testing.T) {
	// Create a generic action result to add header to
	actionResult := mvcapp.NewActionResult([]byte("Needs a body"))
	if actionResult == nil {
		t.Fatal("Failed to create action result")
	}

	// Add the header to the action result
	actionResult.AddHeader("TestHeader", "TestValue")
	if actionResult.Headers["TestHeader"] != "TestValue" {
		t.Error("Failed to set header to action result")
	}

	// Execute the action result to a httptest.ResponseRecorder
	res := httptest.NewRecorder()
	actionResult.Execute(res)

	// Parse the httptest.ResponseRecorder and validates that the header was
	// properly written and received
	if res.Result().Header.Get("TestHeader") != "TestValue" {
		t.Error("Failed to deliver header value to client")
	}
}

// TestActionResult_AddCookie ensures that ActionResult.AddCookie operates as expected
func TestActionResult_AddCookie(t *testing.T) {
	// Create a generic action result to add cookie to
	actionResult := mvcapp.NewActionResult([]byte("Needs a body"))
	if actionResult == nil {
		t.Fatal("Failed to create action result")
	}

	// Create a cookie object to add to the action result
	cookie := &http.Cookie{
		Name:  "TestCookie",
		Value: "TestValue",
	}

	// Add the cookie to the action result
	actionResult.AddCookie(cookie)
	found := false
	for _, v := range actionResult.Cookies {
		if v == cookie {
			found = true
		}
	}

	// Ensure the cookie was added to the collection
	if !found {
		t.Fatal("Failed to set cookie to action result")
	}

	// Execute the action result to a httptest.ResponseRecorder
	res := httptest.NewRecorder()
	actionResult.Execute(res)

	// Ensure that the cookie was delivered
	found = false
	for _, v := range res.Result().Cookies() {
		if v.Name == cookie.Name && v.Value == cookie.Value {
			found = true
		}
	}

	if !found {
		t.Error("Failed to read cookie from http response")
	}
}

// TestActionResult_Execute ensures that ActionResult.Execute operates as expected
func TestActionResult_Execute(t *testing.T) {
	// Create a generic action result to serve
	actionResult := mvcapp.NewActionResult([]byte("Test Payload"))
	if actionResult == nil {
		t.Fatal("Failed to create action result")
	}

	// Executes the action result to a httptest.ResponseRecorder
	res := httptest.NewRecorder()
	actionResult.Execute(res)

	// Reads the body of the response recirder
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Validates that the delivered body is the expected payload
	if string(body) != string(actionResult.Data) {
		t.Error("Failed to retrieve expected payload")
	}
}

// TestActionResult_RawHTML ensures that mvcapp.RawHTML returns the expected raw html template string
func TestActionResult_RawHTML(t *testing.T) {
	data := "<html><head><title>Testing Raw HTML</title></head></html>"
	testData := fmt.Sprintf("%s", mvcapp.RawHTML(data))
	if testData != data {
		t.Error("Failed to process raw html to template string")
	}
}
