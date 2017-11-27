package mvcapp_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/digivance/mvcapp"
)

type rmTestController struct {
	*mvcapp.Controller
}

func newRMTestController(request *http.Request) mvcapp.IController {
	rtn := &rmTestController{
		Controller: mvcapp.NewBaseController(request),
	}

	rtn.BeforeExecute = rtn.OnBeforeExecute
	rtn.AfterExecute = rtn.OnAfterExecute
	rtn.ErrorResult = rtn.OnError
	rtn.NotFoundResult = rtn.OnNotFound

	rtn.RegisterAction("", "Index", rtn.Index)
	rtn.RegisterAction("", "NotFound", rtn.NotFound)
	rtn.RegisterAction("", "DefaultNotFound", rtn.DefaultNotFound)

	return rtn
}

func (controller *rmTestController) OnBeforeExecute() {
	controller.ContinuePipeline = true
}

func (controller *rmTestController) Index(params []string) *mvcapp.ActionResult {
	return controller.Result([]byte("Test Data"))
}

func (controller *rmTestController) NotFound(params []string) *mvcapp.ActionResult {
	return nil
}

func (controller *rmTestController) DefaultNotFound(params []string) *mvcapp.ActionResult {
	controller.NotFoundResult = nil
	return nil
}

func (controller *rmTestController) OnAfterExecute() {
	controller.ContinuePipeline = true
}

func (controller *rmTestController) OnNotFound() *mvcapp.ActionResult {
	return controller.Result([]byte("Not Found"))
}

func (controller *rmTestController) OnError(err error) *mvcapp.ActionResult {
	return controller.Result([]byte("Error"))
}

func TestNewRouteManager(t *testing.T) {
	manager := mvcapp.NewRouteManager()
	if manager == nil {
		t.Error("Failed to create a new route manager")
	}
}

func TestRouteManager_HandleRequest(t *testing.T) {
	// Create a route manager
	manager := mvcapp.NewRouteManager()
	manager.RegisterController("test", newRMTestController)
	manager.DefaultController = "test"
	recorder := httptest.NewRecorder()

	// Ensure that the routes collection has our test controller registered
	if len(manager.Routes) != 1 {
		t.Fatal("Failed to register test controller")
	}

	// test a mapped request for /test/index
	req, err := http.NewRequest("GET", "http://localhost/test/index/param1/param2?qs=value&ext=more#MyFragment", nil)
	if err != nil {
		t.Fatal(err)
	}

	manager.HandleRequest(recorder, req)
	data, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "Test Data" {
		t.Error("Failed to validate known route response data")
		t.Log(string(data))
	}

	// test default mapping (/ should go to /test/index)
	req, err = http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	manager.HandleRequest(recorder, req)
	data, err = ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "Test Data" {
		t.Error("Failed to validate known route response data")
		t.Log(string(data))
	}

	// test an unmapped action
	req, err = http.NewRequest("", "http://localhost/test/404", nil)
	if err != nil {
		t.Fatal(err)
	}

	manager.HandleRequest(recorder, req)
	data, err = ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "Not Found" {
		t.Error("Failed to test custom not found page")
		t.Log(string(data))
	}

	// test an action that fails to return a result with custom not found
	req, err = http.NewRequest("", "http://localhost/test/NotFound", nil)
	if err != nil {
		t.Fatal(err)
	}

	manager.HandleRequest(recorder, req)
	data, err = ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "Not Found" {
		t.Error("Failed to test custom not found page")
		t.Log(string(data))
	}

	// test an action that fails to return a result with the default not found
	req, err = http.NewRequest("", "http://localhost/test/DefaultNotFound", nil)
	if err != nil {
		t.Fatal(err)
	}

	manager.HandleRequest(recorder, req)
	data, err = ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Fatal(err)
	}

	if len(data) < 50 {
		t.Error("Failed to test custom not found page")
		t.Log(string(data))
	}

	// try to request an invalid path / file (controllers)
	req, err = http.NewRequest("", "http://localhost/controllers/testcontroller.go", nil)
	if err != nil {
		t.Fatal(err)
	}

	manager.HandleRequest(recorder, req)
	data, err = ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "Error" {
		t.Error("Failed to test custom error when requesting invalid path page")
		t.Log(string(data))
	}

	// try to request an invalid path / file (emails)
	req, err = http.NewRequest("", "http://localhost/emails/index", nil)
	if err != nil {
		t.Fatal(err)
	}

	manager.HandleRequest(recorder, req)
	data, err = ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "Error" {
		t.Error("Failed to test invalid path")
		t.Log(string(data))
	}

	// try to request an invalid path / file (models)
	req, err = http.NewRequest("", "http://localhost/models/index", nil)
	if err != nil {
		t.Fatal(err)
	}

	manager.HandleRequest(recorder, req)
	data, err = ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "Error" {
		t.Error("Failed to test invalid path")
		t.Log(string(data))
	}

	// try to request an invalid path / file (views)
	req, err = http.NewRequest("", "http://localhost/views/index.htm", nil)
	if err != nil {
		t.Fatal(err)
	}

	manager.HandleRequest(recorder, req)
	data, err = ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "Error" {
		t.Error("Failed to test invalid path")
		t.Log(string(data))
	}

	// test getting a raw file
	apppath := mvcapp.GetApplicationPath()
	os.MkdirAll(fmt.Sprintf("%s/downloads/apps", apppath), 0644)
	defer os.RemoveAll(fmt.Sprintf("%s/downloads", apppath))
	filename := fmt.Sprintf("%s/downloads/apps/pretend.app", apppath)
	payload := []byte("Super cool application thingie here")
	err = ioutil.WriteFile(filename, payload, 0644)
	if err != nil {
		t.Fatal(err)
	}

	req, err = http.NewRequest("GET", "http://localhost/downloads/apps/pretend.app", nil)
	if err != nil {
		t.Fatal(err)
	}

	manager.HandleRequest(recorder, req)
	data, err = ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != string(payload) {
		t.Error("Failed to validate raw file download")
		t.Log(string(data))
	}
}
