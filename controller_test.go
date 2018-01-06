/*
	Digivance MVC Application Framework - Unit Tests
	Controller Feature Tests
	Dan Mayor (dmayor@digivance.com)

	This file defines the version 0.2.0 compatibility of controller.go functions. These functions are written
	to demonstrate and test the intended use cases of the functions in controller.go
*/

package mvcapp_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/digivance/mvcapp"
)

type testController struct {
	*mvcapp.Controller
}

// newTestController is used internally to create a mock controller object that we will
// use during various unit tests below
func newTestController(request *http.Request) mvcapp.IController {
	rtn := &testController{
		mvcapp.NewBaseController(request),
	}

	rtn.RegisterAction("", "Index", rtn.Index)

	rtn.NotFoundResult = rtn.NotFoundPage
	rtn.ErrorResult = rtn.ErrorPage

	return rtn
}

// Index is an emulated action result method that will respond to some of the requests
// built into the unit tests below
func (controller *testController) Index(params []string) *mvcapp.ActionResult {
	return mvcapp.NewActionResult([]byte("test"))
}

// NotFoundPage is an emulated custom 404 handler
func (controller *testController) NotFoundPage() *mvcapp.ActionResult {
	return mvcapp.NewActionResult([]byte("Not Found"))
}

// ErrorPage is an emulated custom error page
func (controller *testController) ErrorPage(err error) *mvcapp.ActionResult {
	return mvcapp.NewActionResult([]byte("Error Page"))
}

// TestController_RegisterAction ensures that the Controller.RegisterAction method operates as expected
func TestController_RegisterAction(t *testing.T) {
	// Make a mock request to test routing and executing
	req, err := http.NewRequest("GET", "http://localhost/test/index", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Invoke an instance of our Test controller (as the route manager would)
	// based on the route maps
	controller := newTestController(req)
	res, err := controller.Execute()
	if err != nil {
		t.Fatalf("404 from expected result: %s", err)
	}

	// compare the resulting payload against the expected vaule "test"
	data := string(res.Data)
	if data != "test" {
		t.Error("Failed to validate paylad data")
	}
}

// TestController_GetCookie ensures that the Controller.GetCookie method returns the expected value
func TestController_GetCookie(t *testing.T) {
	req, err := http.NewRequest("", "http://localhost/test/index", nil)
	if err != nil {
		t.Fatal(err)
	}

	cookie := &http.Cookie{
		Name:  "TestCookie",
		Value: "TestValue",
	}

	req.AddCookie(cookie)

	// Invoke an instance of our Test controller (as the route manager would)
	// based on the route maps
	icontroller := newTestController(req)
	_, err = icontroller.Execute()
	if err != nil {
		t.Fatalf("404 from expected result: %s", err)
	}

	controller := icontroller.ToController()
	testCookie := controller.GetCookie(cookie.Name)
	if testCookie.Value != cookie.Value {
		t.Error("Test cookie not found in controller collection")
	}

	testCookie = controller.GetCookie("FailMe!")
	if testCookie != nil {
		t.Error("Returned unknown cookie, expected nil")
	}
}

// TestController_SetCookie ensures that the Controller.SetCookie method operates as expected
func TestController_SetCookie(t *testing.T) {
	req, err := http.NewRequest("", "http://localhost/test/index", nil)
	if err != nil {
		t.Fatal(err)
	}

	icontroller := newTestController(req)
	controller := icontroller.ToController()
	cookie := &http.Cookie{
		Name:  "TestCookie",
		Value: "TestValue",
	}

	// We try to set twice to test both patterns, creating and overwriting
	// controller cookie's
	controller.SetCookie(cookie)
	controller.SetCookie(cookie)

	testCookie := controller.GetCookie(cookie.Name)
	if testCookie.Value != cookie.Value {
		t.Error("Failed to retrieve test cookie from controller collection")
	}
}

// TestController_DeleteCookie ensures that the Controller.DeleteCookie method operates as expected
func TestController_DeleteCookie(t *testing.T) {
	req, err := http.NewRequest("", "http://localhost/test/index", nil)
	if err != nil {
		t.Fatal(err)
	}

	icontroller := newTestController(req)
	controller := icontroller.ToController()
	cookie := &http.Cookie{
		Name:  "TestCookie",
		Value: "TestValue",
	}

	beforeTime := cookie.Expires

	controller.SetCookie(cookie)
	controller.DeleteCookie(cookie.Name)
	testCookie := controller.GetCookie(cookie.Name)

	if testCookie.Expires == beforeTime {
		t.Error("Failed to reset cookie expiration date")
	}

	if !testCookie.Expires.Before(time.Now()) {
		t.Error("Failed to delete cookie, expiration not set in the past")
	}
}

// TestController_Execute ensures that the Controller.Execute method operates as expected
func TestController_Execute(t *testing.T) {
	req, err := http.NewRequest("POST", "http://localhost/test/index/with/parameters", nil)
	if err != nil {
		t.Fatal(err)
	}

	icontroller := newTestController(req)
	res, err := icontroller.Execute()
	if err != nil {
		t.Errorf("Failed to execute controller action: %s", err)
	}

	data := string(res.Data)
	if data != "test" {
		t.Error("Failed to validate returned value")
	}

	req, err = http.NewRequest("GET", "http://localhost/test/notfound", nil)
	if err != nil {
		t.Fatal(err)
	}

	icontroller = newTestController(req)
	res, err = icontroller.Execute()
	if err != nil {
		t.Errorf("Failed to execute controller action: %s", err)
	}

	data = string(res.Data)
	if data != "Not Found" {
		t.Error("Failed to validate custom content not found page")
	}

	controller := icontroller.ToController()
	controller.NotFoundResult = nil
	res, err = icontroller.Execute()
	if err != nil {
		t.Error("Failed to execute controller action")
	}

	data = string(res.Data)
	if len(data) <= 0 {
		t.Error("Failed to validate default content not found page")
	}
}

// TestController_WriteResponse ensures that the Controller.WriteResponse method operates as expected
func TestController_WriteResponse(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost/test/index", nil)
	if err != nil {
		t.Fatal(err)
	}

	icontroller := newTestController(req)
	recorder := httptest.NewRecorder()

	controller := icontroller.ToController()
	controller.Response = recorder

	res, err := icontroller.Execute()
	if err != nil {
		t.Fatal(err)
	}

	controller.WriteResponse(res)
	data, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Fatal(err)
	}

	body := string(data)
	if body != "test" {
		t.Error("Failed to retrieve expected payload")
	}

	req, err = http.NewRequest("GET", "http://localhost/test/notfound", nil)
	icontroller = newTestController(req)
	recorder = httptest.NewRecorder()

	controller = icontroller.ToController()
	controller.Response = recorder
	controller.WriteResponse(nil)

	data, err = ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Fatal(err)
	}

	body = string(data)
	if body != "Not Found" {
		t.Error("Failed to validate custom content not found result")
	}

	req, err = http.NewRequest("GET", "http://localhost/test/notfound", nil)
	icontroller = newTestController(req)
	recorder = httptest.NewRecorder()

	controller = icontroller.ToController()
	controller.Response = recorder
	controller.NotFoundResult = nil
	controller.WriteResponse(nil)

	data, err = ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Fatal(err)
	}

	body = string(data)
	if len(body) <= 0 {
		t.Error("Failed to validate default content not found result")
	}

	req, err = http.NewRequest("GET", "http://localhost/test/notfound", nil)
	icontroller = newTestController(req)
	recorder = httptest.NewRecorder()

	controller = icontroller.ToController()
	controller.Response = recorder
	controller.WriteResponse(mvcapp.NewActionResult([]byte{}))

	data, err = ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Fatal(err)
	}

	body = string(data)
	if body != "Not Found" {
		t.Error("Failed to validate custom content not found result")
	}

	req, err = http.NewRequest("GET", "http://localhost/test/notfound", nil)
	icontroller = newTestController(req)
	recorder = httptest.NewRecorder()

	controller = icontroller.ToController()
	controller.Response = recorder
	controller.NotFoundResult = nil
	controller.WriteResponse(mvcapp.NewActionResult([]byte{}))

	data, err = ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Fatal(err)
	}

	body = string(data)
	if len(body) <= 0 {
		t.Error("Failed to validate default content not found result")
	}
}

// TestController_RedirectJS ensures that the Controller.RedirectJS method operates as expected
func TestController_RedirectJS(t *testing.T) {
	expectedResult := "<html><head><title>Redirecting...</title><body><script type=\"text/javascript\">window.location.href='https://localhost/test/index';</script></body></html>"
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost/test/index", nil)
	if err != nil {
		t.Fatal(err)
	}

	icontroller := newTestController(req)
	controller := icontroller.ToController()
	controller.Response = recorder

	controller.RedirectJS("https://localhost/test/index")
	data, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Fatal(err)
	}

	body := string(data)
	if body != expectedResult {
		t.Error("Failed to validate expected result data")
	}
}

// TestController_Result ensures that the Controller.Result method returns the expected result
func TestController_Result(t *testing.T) {
	req, err := http.NewRequest("", "http://localhost/test/index", nil)
	if err != nil {
		t.Fatal(err)
	}

	icontroller := newTestController(req)
	controller := icontroller.ToController()

	res := controller.Result([]byte("TestData"))
	data := string(res.Data)

	if data != "TestData" {
		t.Error("Failed to construct generic result data")
	}
}

// TestController_View ensures that the Controller.View method returns the expected result
func TestController_View(t *testing.T) {
	req, err := http.NewRequest("", "http://localhost/test/index", nil)
	if err != nil {
		t.Fatal(err)
	}

	icontroller := newTestController(req)
	controller := icontroller.ToController()

	pathname := fmt.Sprintf("%s/%s", mvcapp.GetApplicationPath(), "views/shared")
	filename := fmt.Sprintf("%s/%s", pathname, "_test_template.htm")
	templateData := "{{ define \"mvcapp\" }}<html><head><title>test template</title></head><body>Hello {{ . }}</body></html>{{ end }}"
	expectedData := "<html><head><title>test template</title></head><body>Hello User</body></html>"
	defer os.RemoveAll(filename)
	defer os.RemoveAll(pathname)
	defer os.RemoveAll(fmt.Sprintf("%s/%s", mvcapp.GetApplicationPath(), "views"))

	os.MkdirAll(pathname, 0644)
	err = ioutil.WriteFile(filename, []byte(templateData), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Construct view result from temporary template file
	viewResult := controller.View([]string{filename}, "User")
	if viewResult == nil {
		t.Fatal("Failed to create view result")
	}

	// Validate the resulting result data
	if string(viewResult.Data) != expectedData {
		t.Error("Failed to validate view result data")
	}

	viewResult = controller.View([]string{""}, "User")
	if viewResult == nil {
		t.Fatal("Failed to create custom error view result")
	}

	controller.ErrorResult = nil
	viewResult = controller.View([]string{""}, "User")
	if viewResult == nil {
		t.Fatal("Failed to construct default error view result")
	}
}

// TestController_SimpleView ensures that the Controller.SimpleView method returns the expected result
func TestController_SimpleView(t *testing.T) {
	req, err := http.NewRequest("", "http://localhost/test/index", nil)
	if err != nil {
		t.Fatal(err)
	}

	icontroller := newTestController(req)
	controller := icontroller.ToController()

	pathname := fmt.Sprintf("%s/%s", mvcapp.GetApplicationPath(), "views/shared")
	filename := fmt.Sprintf("%s/%s", pathname, "_test_template.htm")
	templateData := "{{ define \"mvcapp\" }}<html><head><title>test template</title></head><body>Hello Template!</body></html>{{ end }}"
	expectedData := "<html><head><title>test template</title></head><body>Hello Template!</body></html>"

	defer os.RemoveAll(filename)
	defer os.RemoveAll(pathname)
	defer os.RemoveAll(fmt.Sprintf("%s/%s", mvcapp.GetApplicationPath(), "views"))

	os.MkdirAll(pathname, 0644)
	err = ioutil.WriteFile(filename, []byte(templateData), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Construct view result from temporary template file
	viewResult := controller.SimpleView(filename)
	if viewResult == nil {
		t.Fatal("Failed to create view result")
	}

	// Validate the resulting result data
	if string(viewResult.Data) != expectedData {
		t.Error("Failed to validate view result data")
	}
}

// TestController_JSON ensures that the Controll.JSON method returns the expected result
func TestController_JSON(t *testing.T) {
	req, err := http.NewRequest("", "http://localhost/test/index", nil)
	if err != nil {
		t.Fatal(err)
	}

	icontroller := newTestController(req)
	controller := icontroller.ToController()

	jsonResult := controller.JSON("Test Data")
	if jsonResult == nil {
		t.Fatal("Failed to construct json result")
	}

	data := ""
	err = json.Unmarshal(jsonResult.Data, &data)
	if err != nil {
		t.Fatal(err)
	}

	if data != "Test Data" {
		t.Error("Failed to validate json encoded result data")
	}

	jsonResult = controller.JSON(nil)
	if jsonResult == nil {
		t.Fatal("Failed to construct API failure json result")
	}

	// We set success to true then overwrite it with the value returned
	// from the call to JSON. We expect success to equal false if this
	// part of the test is to pass
	expectedResult := "{\"Success\":false,\"Error\":\"Failed to create json payload\"}"
	if !strings.EqualFold(string(jsonResult.Data), expectedResult) {
		t.Fatalf("Error comparing json result:\n%s", jsonResult.Data)
	}
}
