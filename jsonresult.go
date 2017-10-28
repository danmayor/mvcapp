package mvcapp

import (
	"encoding/json"
	"net/http"
)

// JSONResult is used to encode a provided payload into a json result string
type JSONResult struct {
	IActionResult
	*ActionResult
}

// NewJSONResult returns a new JSONResult with the payload json encoded to Data
func NewJSONResult(payload interface{}) *JSONResult {
	data, _ := json.Marshal(payload)

	return &JSONResult{
		ActionResult: NewActionResult(data),
	}
}

// Execute is called to render the data to the client (browser)
func (result *JSONResult) Execute(response http.ResponseWriter) error {
	if _, err := response.Write(result.Data); err != nil {
		return err
	}

	return nil
}

// ToActionResult returns the base action result struct
func (result *JSONResult) ToActionResult() *ActionResult {
	return result.ActionResult
}
