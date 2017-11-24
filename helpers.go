package mvcapp

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
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

// appPath is used internally so that we don't have to query the os args
// every time we ask to GetApplicationPath
var appPath = ""

// GetApplicationPath should return the full path to the executable.
// This is the root of the site and where the assembly file is
func GetApplicationPath() string {
	if appPath == "" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			appPath = "."
		}

		appPath = dir
	}

	return appPath
}

// logFilename is used internally to hold the name of the file that holds our
// application logs
var logFilename = ""

// GetLogFilename returns the current, or default log file that we will write to
func GetLogFilename() string {
	return logFilename
}

// SetLogFilename will set the filename that log messages will be written to
func SetLogFilename(filename string) {
	logFilename = filename
}

// logLevel is the internal value representing what levels of log messages are written
// to our log file. Where 0 = Off 1 = Errors Only, 2 = Warnings (Such as 404),
// 3 = Verbose (It'll say a lot)
var logLevel = 1

// GetLogLevel returns the level of log messages that are written to our log file
func GetLogLevel() int {
	return logLevel
}

// SetLogLevel sets the internal log level of messages that are written to our log file
// Where 0 = Off 1 = Errors Only, 2 = Warnings (Such as 404), 3 = Verbose (It'll say a lot)
func SetLogLevel(level int) {
	logLevel = level
}

// LogMessage writes an information message to the log file if our internal log level is 3
func LogMessage(message string) {
	if logLevel >= 3 && logFilename != "" {
		f, err := os.OpenFile(logFilename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
		if err != nil {
			return
		}

		defer f.Close()
		f.WriteString(fmt.Sprintf("[%s] Information: %s\r\n", time.Now().String(), message))
	}
}

// LogWarning writes a warning message to the log file if our internal log level is >= 2
func LogWarning(message string) {
	if logLevel >= 2 && logFilename != "" {
		f, err := os.OpenFile(logFilename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
		if err != nil {
			return
		}

		defer f.Close()
		f.WriteString(fmt.Sprintf("[%s] Warning: %s\r\n", time.Now().String(), message))
	}
}

// LogError writes an error message to the log file if our internal log level is >= 1
func LogError(message string) {
	if logLevel >= 1 && logFilename != "" {
		f, err := os.OpenFile(logFilename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
		if err != nil {
			return
		}

		defer f.Close()
		f.WriteString(fmt.Sprintf("[%s] Critical: %s\r\n\r\n", time.Now().String(), message))
	}
}
