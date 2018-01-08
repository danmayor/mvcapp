/*
	Digivance MVC Application Framework - Unit Tests
	Helper Functions Tests
	Dan Mayor (dmayor@digivance.com)

	This file defines the version 0.3.0 compatibility of helpers.go functions. These functions are written
	to demonstrate and test the intended use cases of the functions in helpers.go
*/

package mvcapp_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/digivance/mvcapp"
)

// TestTemplateExists ensures the mvcapp.TemplateExists method returns the expected value
func TestTemplateExists(t *testing.T) {
	var err error
	apppath := mvcapp.GetApplicationPath()

	os.MkdirAll(fmt.Sprintf("%s/%s", apppath, "views/test"), 0644)
	os.MkdirAll(fmt.Sprintf("%s/%s", apppath, "views/shared"), 0644)
	os.MkdirAll(fmt.Sprintf("%s/%s", apppath, "views/shared/test"), 0644)

	defer os.RemoveAll(fmt.Sprintf("%s/%s", apppath, "views/test"))
	defer os.RemoveAll(fmt.Sprintf("%s/%s", apppath, "views/shared/test"))
	defer os.RemoveAll(fmt.Sprintf("%s/%s", apppath, "views/shared"))
	defer os.RemoveAll(fmt.Sprintf("%s/%s", apppath, "views"))

	filename := fmt.Sprintf("%s/%s", apppath, "_test_teplate.htm")
	err = ioutil.WriteFile(filename, []byte(apppath), 0644)
	if err != nil {
		t.Fatal(err)
	}

	if !mvcapp.TemplateExists("test", filename) {
		t.Error("Failed to test direct path")
	}
	os.RemoveAll(filename)

	filename = fmt.Sprintf("%s/%s", apppath, "views/_test_template.htm")
	err = ioutil.WriteFile(filename, []byte(apppath), 0644)
	if err != nil {
		t.Fatal(err)
	}

	if !mvcapp.TemplateExists("test", "./_test_template.htm") {
		t.Error("Failed to test ./views path")
	}
	os.RemoveAll(filename)

	filename = fmt.Sprintf("%s/%s", apppath, "/views/test/_test_template.htm")
	err = ioutil.WriteFile(filename, []byte(apppath), 0644)
	if err != nil {
		t.Fatal(err)
	}

	if !mvcapp.TemplateExists("test", "_test_template.htm") {
		t.Error("Failed to test ./views/controller path")
	}
	os.RemoveAll(filename)

	filename = fmt.Sprintf("%s/%s", apppath, "/views/shared/_test_template.htm")
	err = ioutil.WriteFile(filename, []byte(apppath), 0644)
	if err != nil {
		t.Fatal(err)
	}

	if !mvcapp.TemplateExists("test", "_test_template.htm") {
		t.Error("Failed to test ./views/shared path")
	}
	os.RemoveAll(filename)

	filename = fmt.Sprintf("%s/%s", apppath, "/views/shared/test/_test_template.htm")
	err = ioutil.WriteFile(filename, []byte(apppath), 0644)
	if err != nil {
		t.Fatal(err)
	}

	if !mvcapp.TemplateExists("test", "_test_template.htm") {
		t.Error("Failed to test ./views/shared/controller path")
	}
	os.RemoveAll(filename)

	if mvcapp.TemplateExists("test", "_test_template.htm") {
		t.Error("Failed to detect that file is NOT present :(")
	}
}

// TestMakeTemplateList ensures the mvcapp.MakeTemplateList method returns the expected values
func TestMakeTemplateList(t *testing.T) {
	var err error
	apppath := mvcapp.GetApplicationPath()

	os.MkdirAll(fmt.Sprintf("%s/%s", apppath, "views/test"), 0644)
	os.MkdirAll(fmt.Sprintf("%s/%s", apppath, "views/shared"), 0644)
	os.MkdirAll(fmt.Sprintf("%s/%s", apppath, "views/shared/test"), 0644)

	defer os.RemoveAll(fmt.Sprintf("%s/%s", apppath, "views/test"))
	defer os.RemoveAll(fmt.Sprintf("%s/%s", apppath, "views/shared/test"))
	defer os.RemoveAll(fmt.Sprintf("%s/%s", apppath, "views/shared"))
	defer os.RemoveAll(fmt.Sprintf("%s/%s", apppath, "views"))

	filename := fmt.Sprintf("%s/%s", apppath, "/views/_start.tpl")
	defer os.RemoveAll(filename)
	err = ioutil.WriteFile(filename, []byte(apppath), 0644)
	if err != nil {
		t.Fatal(err)
	}

	filename = fmt.Sprintf("%s/%s", apppath, "/views/shared/_shared.tpl")
	defer os.RemoveAll(filename)
	err = ioutil.WriteFile(filename, []byte(apppath), 0644)
	if err != nil {
		t.Fatal(err)
	}

	filename = fmt.Sprintf("%s/%s", apppath, "/views/shared/test/_sharedwidget.tpl")
	defer os.RemoveAll(filename)
	err = ioutil.WriteFile(filename, []byte(apppath), 0644)
	if err != nil {
		t.Fatal(err)
	}

	filename = fmt.Sprintf("%s/%s", apppath, "/views/test/_page.tpl")
	defer os.RemoveAll(filename)
	err = ioutil.WriteFile(filename, []byte(apppath), 0644)
	if err != nil {
		t.Fatal(err)
	}

	filename = fmt.Sprintf("%s/%s", apppath, "_root_template.tpl")
	defer os.RemoveAll(filename)
	err = ioutil.WriteFile(filename, []byte(apppath), 0644)
	if err != nil {
		t.Fatal(err)
	}

	templateList := mvcapp.MakeTemplateList("test", []string{"_start.tpl", "_shared.tpl", "_sharedwidget.tpl", "_page.tpl", "./_root_template.tpl"})
	if len(templateList) != 5 {
		t.Error("Failed to parse all standard paths for templates")
	}
}

// TestGetLogFilename ensures the mvcapp.GetLogFilename method returns the expected value
func TestGetLogFilename(t *testing.T) {
	apppath := mvcapp.GetApplicationPath()
	filename := fmt.Sprintf("%s/%s", apppath, "mvcapp.log")
	mvcapp.SetLogFilename(filename)
	if mvcapp.GetLogFilename() != filename {
		t.Error("Failed to set and validate log file")
	}
}

// TestGetLogLevel ensures the mvcapp.GetLogLevel method returns the expected value
func TestGetLogLevel(t *testing.T) {
	level := 3
	mvcapp.SetLogLevel(level)
	if mvcapp.GetLogLevel() != level {
		t.Error("Failed to set and validate log level")
	}
}

// TestLogMessage ensures the mvcapp.LogMessage method operates as expected
func TestLogMessage(t *testing.T) {
	apppath := mvcapp.GetApplicationPath()
	filename := fmt.Sprintf("%s/%s", apppath, "mvcapp_test.log")
	os.RemoveAll(filename)

	mvcapp.SetLogFilename(filename)
	mvcapp.SetLogLevel(mvcapp.LogLevelInfo)

	mvcapp.LogMessage("Hello logs!")
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()
	defer os.RemoveAll(filename)

	fi, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}

	if fi.Size() < 12 {
		t.Error("Failed to log message (womp womp)")
	}

	mvcapp.SetLogFilename("")
	if err := mvcapp.LogMessage("Should fail"); err == nil {
		t.Error("Failed to prevent writing to missing filename")
	}
}

// TestLogMessagef ensures the mvcapp.LogMessagef method operates as expected
func TestLogMessagef(t *testing.T) {
	apppath := mvcapp.GetApplicationPath()
	filename := fmt.Sprintf("%s/%s", apppath, "mvcapp_test.log")
	os.RemoveAll(filename)

	mvcapp.SetLogFilename(filename)
	mvcapp.SetLogLevel(mvcapp.LogLevelInfo)

	mvcapp.LogMessagef("Hello %s", "logs!")
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()
	defer os.RemoveAll(filename)

	fi, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}

	if fi.Size() < 12 {
		t.Error("Failed to log message (womp womp)")
	}

	mvcapp.SetLogFilename("")
	if err := mvcapp.LogMessage("Should fail"); err == nil {
		t.Error("Failed to prevent writing to missing filename")
	}
}

// TestLogWarning ensures the mvcapp.LogWarning method operates as expected
func TestLogWarning(t *testing.T) {
	apppath := mvcapp.GetApplicationPath()
	filename := fmt.Sprintf("%s/%s", apppath, "mvcapp_test.log")
	os.RemoveAll(filename)

	mvcapp.SetLogFilename(filename)
	mvcapp.SetLogLevel(mvcapp.LogLevelWarning)

	mvcapp.LogWarning("Hello logs!")
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()
	defer os.RemoveAll(filename)

	fi, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}

	if fi.Size() < 12 {
		t.Error("Failed to log warning (womp womp)")
	}

	mvcapp.SetLogFilename("")
	if err := mvcapp.LogWarning("Should fail"); err == nil {
		t.Error("Failed to prevent writing to missing filename")
	}
}

// TestLogWarningf ensures the mvcapp.LogWarningf method operates as expected
func TestLogWarningf(t *testing.T) {
	apppath := mvcapp.GetApplicationPath()
	filename := fmt.Sprintf("%s/%s", apppath, "mvcapp_test.log")
	os.RemoveAll(filename)

	mvcapp.SetLogFilename(filename)
	mvcapp.SetLogLevel(mvcapp.LogLevelWarning)

	mvcapp.LogWarningf("Hello %s", "logs!")
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()
	defer os.RemoveAll(filename)

	fi, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}

	if fi.Size() < 12 {
		t.Error("Failed to log warning (womp womp)")
	}

	mvcapp.SetLogFilename("")
	if err := mvcapp.LogWarning("Should fail"); err == nil {
		t.Error("Failed to prevent writing to missing filename")
	}
}

// TestLogError ensures the mvcapp.LogError method operates as expected
func TestLogError(t *testing.T) {
	apppath := mvcapp.GetApplicationPath()
	filename := fmt.Sprintf("%s/%s", apppath, "mvcapp_test.log")
	os.RemoveAll(filename)

	mvcapp.SetLogFilename(filename)
	mvcapp.SetLogLevel(mvcapp.LogLevelError)

	mvcapp.LogError("Hello logs!")
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()
	defer os.RemoveAll(filename)

	fi, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}

	if fi.Size() < 12 {
		t.Error("Failed to log error (womp womp)")
	}

	mvcapp.SetLogFilename("")
	if err := mvcapp.LogError("Should fail"); err == nil {
		t.Error("Failed to prevent writing to missing filename")
	}
}

// TestLogErrorf ensures the mvcapp.LogErrorf method operates as expected
func TestLogErrorf(t *testing.T) {
	apppath := mvcapp.GetApplicationPath()
	filename := fmt.Sprintf("%s/%s", apppath, "mvcapp_test.log")
	os.RemoveAll(filename)

	mvcapp.SetLogFilename(filename)
	mvcapp.SetLogLevel(mvcapp.LogLevelError)

	mvcapp.LogErrorf("Hello %s", "logs!")
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()
	defer os.RemoveAll(filename)

	fi, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}

	if fi.Size() < 12 {
		t.Error("Failed to log error (womp womp)")
	}

	mvcapp.SetLogFilename("")
	if err := mvcapp.LogError("Should fail"); err == nil {
		t.Error("Failed to prevent writing to missing filename")
	}
}

// TestLogTrace ensures the mvcapp.LogTrace method operates as expected
func TestLogTrace(t *testing.T) {
	apppath := mvcapp.GetApplicationPath()
	filename := fmt.Sprintf("%s/%s", apppath, "mvcapp_test.log")
	os.RemoveAll(filename)

	mvcapp.SetLogFilename(filename)
	mvcapp.SetLogLevel(mvcapp.LogLevelTrace)

	mvcapp.LogTrace("Hello logs!")
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()
	defer os.RemoveAll(filename)

	fi, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}

	if fi.Size() < 12 {
		t.Error("Failed to log error (womp womp)")
	}

	mvcapp.SetLogFilename("")
	if err := mvcapp.LogTrace("Should fail"); err == nil {
		t.Error("Failed to prevent writing to missing filename")
	}
}

// TestLogTracef ensures the mvcapp.LogTracef method operates as expected
func TestLogTracef(t *testing.T) {
	apppath := mvcapp.GetApplicationPath()
	filename := fmt.Sprintf("%s/%s", apppath, "mvcapp_test.log")
	os.RemoveAll(filename)

	mvcapp.SetLogFilename(filename)
	mvcapp.SetLogLevel(mvcapp.LogLevelTrace)

	mvcapp.LogTracef("Hello %s", "logs!")
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()
	defer os.RemoveAll(filename)

	fi, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}

	if fi.Size() < 12 {
		t.Error("Failed to log error (womp womp)")
	}

	mvcapp.SetLogFilename("")
	if err := mvcapp.LogTrace("Should fail"); err == nil {
		t.Error("Failed to prevent writing to missing filename")
	}
}

// TestGetLogDateFormat ensures that the mvcapp.GetLogDateFormat method returns the
// expected value
func TestGetLogDateFormat(t *testing.T) {
	data := mvcapp.GetLogDateFormat()
	if data == "" {
		t.Error("Failed to get date format")
	}
}

// TestSetLogDateFormat ensures that the mvcapp.SetLogDateFormat method operates as expected
func TestSetLogDateFormat(t *testing.T) {
	mvcapp.SetLogDateFormat("1/2/3")
	data := mvcapp.GetLogDateFormat()
	if data != "1/2/3" {
		t.Error("Failed to set log date format")
	}
}
