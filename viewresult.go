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

// MakeTemplateList provides some common view template path fallbacks. Will test
// if each of the template file names exist as is, if not will try the following
// 	./Views/template
// 	./Views/controllerName/template
// 	./Views/Shared/template
// 	./Views/Shared/controllerName/template
func MakeTemplateList(controllerName string, templates []string) []string {
	rtn := []string{}

	for _, template := range templates {
		if _, err := os.Stat(template); !os.IsNotExist(err) {
			rtn = append(rtn, template)
		} else {
			// Try /Views/template
			viewPath := fmt.Sprintf("./Views/%s", template)
			if _, err := os.Stat(viewPath); !os.IsNotExist(err) {
				rtn = append(rtn, viewPath)
			} else {
				// Try /Views/controllerName/template
				controllerPath := fmt.Sprintf("./Views/%s/%s", controllerName, template)
				if _, err := os.Stat(controllerPath); !os.IsNotExist(err) {
					rtn = append(rtn, controllerPath)
				} else {
					// Try /Views/Shared/template
					sharedPath := fmt.Sprintf("./Views/Shared/%s", template)
					if _, err := os.Stat(sharedPath); !os.IsNotExist(err) {
						rtn = append(rtn, sharedPath)
					} else {
						// Try /Views/Shared/controllerName/template
						sharedControllerPath := fmt.Sprintf("./Views/Shared/%s/%s", controllerName, template)
						if _, err := os.Stat(sharedControllerPath); !os.IsNotExist(err) {
							rtn = append(rtn, sharedControllerPath)
						}
					}
				}
			}
		}
	}

	return rtn
}
