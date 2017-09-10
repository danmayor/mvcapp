/*
	Digivance MVC Application Framework
	View Result Object
	Dan Mayor (dmayor@digivance.com)

	This file defines the View Result functionality
*/

package mvcapp

import (
	"html/template"
	"net/http"
)

// ViewResult is a derivitive of the ActionResult struct and
// is used to render a template to the client as html
type ViewResult struct {
	IActionResult

	Headers   map[string]string
	Templates []string
	Model     interface{}
}

// NewViewResult returns a new ViewResult struct with the Data
// member set to the compiled templates requested
func NewViewResult(templates []string, model interface{}) *ViewResult {
	return &ViewResult{
		Headers:   map[string]string{},
		Templates: templates,
		Model:     model,
	}
}

func (result *ViewResult) AddHeader(key string, val string) {
	result.Headers[key] = val
}

// Execute will compile and execute the templates requested with the provided model
func (result *ViewResult) Execute(response http.ResponseWriter) (int, error) {
	page, err := template.ParseFiles(result.Templates...)
	if err != nil {
		return 500, err
	}

	for k, v := range result.Headers {
		response.Header().Set(k, v)
	}

	if err = page.ExecuteTemplate(response, "mvcapp", result.Model); err != nil {
		return 500, err
	}

	return 200, nil
}
