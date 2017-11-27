/*
	Digivance MVC Application Framework - Unit Tests
	Helper Functions Tests
	Dan Mayor (dmayor@digivance.com)

	This file defines the version 0.1.0 compatibility of helpers.go functions. These functions are written
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

	if !mvcapp.TemplateExists("test", "_test_template.htm") {
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

	templateList := mvcapp.MakeTemplateList("test", []string{"_start.tpl", "_shared.tpl", "_sharedwidget.tpl", "_page.tpl"})
	if len(templateList) != 4 {
		t.Error("Failed to parse all standard paths for templates")
	}
}

func TestGetLogFilename(t *testing.T) {
	apppath := mvcapp.GetApplicationPath()
	filename := fmt.Sprintf("%s/%s", apppath, "mvcapp.log")
	mvcapp.SetLogFilename(filename)
	if mvcapp.GetLogFilename() != filename {
		t.Error("Failed to set and validate log file")
	}
}

func TestGetLogLevel(t *testing.T) {
	level := 3
	mvcapp.SetLogLevel(level)
	if mvcapp.GetLogLevel() != level {
		t.Error("Failed to set and validate log level")
	}
}

func TestLogMessage(t *testing.T) {
	apppath := mvcapp.GetApplicationPath()
	filename := fmt.Sprintf("%s/%s", apppath, "mvcapp_test.log")
	os.RemoveAll(filename)

	mvcapp.SetLogFilename(filename)
	mvcapp.SetLogLevel(3)

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
}

func TestLogWarning(t *testing.T) {
	apppath := mvcapp.GetApplicationPath()
	filename := fmt.Sprintf("%s/%s", apppath, "mvcapp_test.log")
	os.RemoveAll(filename)

	mvcapp.SetLogFilename(filename)
	mvcapp.SetLogLevel(3)

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
}

func TestLogError(t *testing.T) {
	apppath := mvcapp.GetApplicationPath()
	filename := fmt.Sprintf("%s/%s", apppath, "mvcapp_test.log")
	os.RemoveAll(filename)

	mvcapp.SetLogFilename(filename)
	mvcapp.SetLogLevel(3)

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
}
