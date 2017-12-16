/*
	Digivance MVC Application Framework
	Action Result Features
	Dan Mayor (dmayor@digivance.com)

	This file defines the base action result functionality
*/

package mvcapp

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
)

// ActionResult is a base level struct that implements the Execute
// method and provides the Data []byte member
type ActionResult struct {
	// StatusCode is the HTTP status code to write with this response. default is 200 ok
	StatusCode int

	// Headers is a key value pairs map of the names and values of headers to write with this response
	Headers map[string]string

	// Cookies is a collection of http cookie values to write with this response
	Cookies []*http.Cookie

	// Data is the raw byte array representing the payload to deliver
	Data []byte
}

// NewActionResult returns a new action result populated with the provided data
func NewActionResult(data []byte) *ActionResult {
	return &ActionResult{
		StatusCode: 200,
		Headers:    map[string]string{},
		Cookies:    []*http.Cookie{},
		Data:       data,
	}
}

func RawHTML(data string) template.HTML {
	return template.HTML(data)
}

// NewViewResult returns a new ViewResult struct with the Data
// member set to the compiled templates requested
func NewViewResult(templates []string, model interface{}) *ActionResult {
	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,
		"RawHTML": RawHTML,
	}

	page, err := template.New("ViewTemplate").Funcs(funcMap).ParseFiles(templates...)

	if err != nil {
		LogError(err.Error())
		return nil
	}

	buffer := new(bytes.Buffer)
	if err = page.ExecuteTemplate(buffer, "mvcapp", model); err != nil {
		LogError(err.Error())
		return nil
	}

	return NewActionResult(buffer.Bytes())
}

// NewJSONResult returns a new JSONResult with the payload json encoded to Data
func NewJSONResult(payload interface{}) *ActionResult {
	data, err := json.Marshal(payload)
	if len(data) <= 0 || payload == nil {
		if err != nil {
			LogError(err.Error())
		}

		return nil
	}

	return NewActionResult(data)
}

// AddHeader adds an http header key value pair combination to the result
func (result *ActionResult) AddHeader(key string, val string) {
	result.Headers[key] = val
}

// AddCookie adds the provided cookie to the result
func (result *ActionResult) AddCookie(cookie *http.Cookie) {
	result.Cookies = append(result.Cookies, cookie)
}

// Execute writes the header, cookies and data of this action result to the client.
func (result ActionResult) Execute(response http.ResponseWriter) {
	for k, v := range result.Headers {
		response.Header().Set(k, v)
	}

	for _, cookie := range result.Cookies {
		http.SetCookie(response, cookie)
	}

	response.WriteHeader(result.StatusCode)
	response.Write(result.Data)
}
