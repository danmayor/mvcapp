/*
	Digivance MVC Application Framework
	Action Result Features
	Dan Mayor (dmayor@digivance.com)

	This file defines the base action result functionality
*/

package mvcapp

import (
	"errors"
	"net/http"
)

// IActionResult is used internally to require the execute method
type IActionResult interface {
	Execute(http.ResponseWriter) (int, error)
}

// ActionResult is a base level struct that implements the Execute
// method and provides the Data []byte member
type ActionResult struct {
	IActionResult

	Headers map[string]string
	Data    []byte
}

// NewActionResult returns a new action result populated with the provided data
func NewActionResult(data []byte) *ActionResult {
	return &ActionResult{
		Data:    data,
		Headers: map[string]string{},
	}
}

// AddHeader adds an http header key value pair combination to the result
func (result *ActionResult) AddHeader(key string, val string) {
	result.Headers[key] = val
}

// Execute writes the raw data to the client
func (result ActionResult) Execute(response http.ResponseWriter) (int, error) {
	if len(result.Data) <= 0 {
		return 404, errors.New("No response from request")
	}

	response.Write(result.Data)
	return 200, nil
}
