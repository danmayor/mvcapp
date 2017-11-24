package mvcapp_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/digivance/mvcapp"
)

func TestNewActionResult(t *testing.T) {
	actionResult := mvcapp.NewActionResult([]byte("Version 0.1.0 Compliant"))
	if actionResult == nil {
		t.Fatal("Failed to create new action result")
	}

	if string(actionResult.Data) != "Version 0.1.0 Compliant" {
		t.Error("Failed to validate result data")
	}
}

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
	viewResult := mvcapp.NewViewResult([]string{filename}, nil)
	if viewResult == nil {
		t.Fatal("Failed to create view result")
	}

	// Validate the resulting result data
	if string(viewResult.Data) != expectedResultData {
		t.Error("Failed to validate view result data")
	}
}

func TestNewJSONResult(t *testing.T) {
	// Create a json encoded payload
	payload := "Version 0.1.0 Compliant"
	jsonResult := mvcapp.NewJSONResult(payload)
	if jsonResult == nil {
		t.Fatal("Failed to create JSON result")
	}

	// Deserialize the created json byte array
	var res string
	err := json.Unmarshal(jsonResult.Data, &res)
	if err != nil {
		t.Fatal(err)
	}

	// Test that the returned value is the intended payload
	if res != payload {
		t.Error("Failed to validate payload")
	}
}
