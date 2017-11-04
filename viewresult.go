/*
	Digivance MVC Application Framework
	View Result Features
	Dan Mayor (dmayor@digivance.com)

	This file defines the basic functionality of a ViewResult. View results represent a raw
	content result that is rendered to the browser.

	This package is released under as open source under the LGPL-3.0 which can be found:
	https://opensource.org/licenses/LGPL-3.0
*/

package mvcapp

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/digivance/applog"
	"github.com/digivance/str"
)

// ViewResult is a derivitive of the ActionResult struct and
// is used to render a template to the client as html
type ViewResult struct {
	IActionResult
	*ActionResult

	Model     interface{}
	Templates []string
}

// NewViewResult returns a new ViewResult struct with the Data
// member set to the compiled templates requested
func NewViewResult(templates []string, model interface{}) *ViewResult {
	return &ViewResult{
		ActionResult: NewActionResult([]byte{}),
		Templates:    templates,
		Model:        model,
	}
}

// Execute will compile and execute the templates requested with the provided model
func (result *ViewResult) Execute(response http.ResponseWriter) error {
	funcMap := template.FuncMap{
		"ToUpper": str.ToUpper,
		"ToLower": str.ToLower,
	}

	page, err := template.New("ViewTemplate").Funcs(funcMap).ParseFiles(result.Templates...)

	if err != nil {
		applog.WriteString(fmt.Sprintf("Failed to execute view result: %s", err.Error()))
		return err
	}

	for k, v := range result.Headers {
		response.Header().Set(k, v)
	}

	if err = page.ExecuteTemplate(response, "mvcapp", result.Model); err != nil {
		applog.WriteString(fmt.Sprintf("Failed to execute view result: %s", err.Error()))
		return err
	}

	return nil
}

// ToActionResult returns a pointer to the base action result class
func (result *ViewResult) ToActionResult() *ActionResult {
	return result.ActionResult
}

// TemplateExists checks the standard folder paths based on the provided controllerName
// to see if the template file can be found. (See MakeTemplateList for path structure)
func TemplateExists(controllerName string, template string) bool {
	if _, err := os.Stat(template); !os.IsNotExist(err) {
		return true
	}

	// Try /views/template
	viewPath := fmt.Sprintf("%s/views/%s", GetApplicationPath(), template)
	if _, err := os.Stat(viewPath); !os.IsNotExist(err) {
		return true
	}

	// Try /Views/controllerName/template
	controllerPath := fmt.Sprintf("%s/views/%s/%s", GetApplicationPath(), controllerName, template)
	if _, err := os.Stat(controllerPath); !os.IsNotExist(err) {
		return true
	}

	// Try /views/shared/template
	sharedPath := fmt.Sprintf("%s/views/shared/%s", GetApplicationPath(), template)
	if _, err := os.Stat(sharedPath); !os.IsNotExist(err) {
		return true
	}

	// Try /views/shared/controllerName/template
	sharedControllerPath := fmt.Sprintf("%s/views/shared/%s/%s", GetApplicationPath(), controllerName, template)
	if _, err := os.Stat(sharedControllerPath); !os.IsNotExist(err) {
		return true
	}

	return false
}

// MakeTemplateList provides some common view template path fallbacks. Will test
// if each of the template file names exist as is, if not will try the following:
//
// 	./views/template
// 	./views/controllerName/template
// 	./views/shared/template
// 	./views/shared/controllerName/template
func MakeTemplateList(controllerName string, templates []string) []string {
	rtn := []string{}

	for _, template := range templates {
		if _, err := os.Stat(template); !os.IsNotExist(err) {
			rtn = append(rtn, template)
		} else {
			// Try /views/template
			viewPath := fmt.Sprintf("%s/views/%s", GetApplicationPath(), template)
			if _, err := os.Stat(viewPath); !os.IsNotExist(err) {
				rtn = append(rtn, viewPath)
			} else {
				// Try /Views/controllerName/template
				controllerPath := fmt.Sprintf("%s/views/%s/%s", GetApplicationPath(), controllerName, template)
				if _, err := os.Stat(controllerPath); !os.IsNotExist(err) {
					rtn = append(rtn, controllerPath)
				} else {
					// Try /views/shared/template
					sharedPath := fmt.Sprintf("%s/views/shared/%s", GetApplicationPath(), template)
					if _, err := os.Stat(sharedPath); !os.IsNotExist(err) {
						rtn = append(rtn, sharedPath)
					} else {
						// Try /views/shared/controllerName/template
						sharedControllerPath := fmt.Sprintf("%s/views/shared/%s/%s", GetApplicationPath(), controllerName, template)
						if _, err := os.Stat(sharedControllerPath); !os.IsNotExist(err) {
							rtn = append(rtn, sharedControllerPath)
						} else {
							// TODO: Add 404 page here
						}
					}
				}
			}
		}
	}

	return rtn
}
