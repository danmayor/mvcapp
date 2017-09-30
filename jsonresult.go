package mvcapp

import (
	"encoding/json"
	"net/http"
)

// JsonResult is used to encode a provided payload into a json result string
type JsonResult struct {
	IActionResult
	*ActionResult
}

// NewJsonResult returns a new JsonResult with the payload json encoded to Data
func NewJsonResult(payload interface{}) *JsonResult {
	data, _ := json.Marshal(payload)

	return &JsonResult{
		ActionResult: NewActionResult(data),
	}
}

// Execute is called to render the data to the client (browser)
func (result *JsonResult) Execute(response http.ResponseWriter) (int, error) {
	if _, err := response.Write(result.Data); err != nil {
		return 500, err
	}

	return 200, nil
}
