package mvcapp

import (
	"fmt"
	"math/rand"
	"os"
	"time"
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

// Some constant configuration values for random string generation methods
const (
	// letterBytes : Available characters for random string
	letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// letterIDBits : Used in reduced byte masking
	letterIDBits = 6

	// letterIDMask : Used in reduced byte masking
	letterIDMask = 1<<letterIDBits - 1

	// letterIDMax : Used in reduced byte masking
	letterIDMax = 63 / letterIDBits
)

// randomizer : Internal rand.Source
var randomizer = rand.NewSource(time.Now().UnixNano())

// RandomString returns a randomly generated string of the given length.
func RandomString(length int) string {
	data := make([]byte, length)
	for i, cache, remain := length-1, randomizer.Int63(), letterIDMax; i >= 0; {
		if remain == 0 {
			cache, remain = randomizer.Int63(), letterIDMax
		}

		if id := int(cache & letterIDMask); id < len(letterBytes) {
			data[i] = letterBytes[id]
			i--
		}

		cache >>= letterIDBits
		remain--
	}

	return string(data)
}
