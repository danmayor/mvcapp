package mvcapp

import (
	"fmt"
	"os"
)

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
