package mvcapp

import (
	"fmt"
	"net/http/httptest"
	"testing"
)

type TestModel struct {
	Title   string
	Welcome string
}

type TestController struct {
	Controller
}

func NewTestController() *TestController {
	rtn := &TestController{}
	rtn.RegisterAction("GET", "Index", rtn.Index)
	return rtn
}

func (controller TestController) Index(params []string) IActionResult {
	templates := []string{"testindex.htm"}
	model := TestModel{
		Title:   "Test Controller",
		Welcome: "This message is from the data model!",
	}

	result := NewViewResult(templates, model)
	result.AddHeader("Framework", "Digivance MVC Application Framework")

	return result
}

func TestRouteManager(t *testing.T) {
	controller := NewTestController()
	mgr := NewRouteManager()
	mgr.RegisterController("Test", controller)

	request := httptest.NewRequest("GET", "/Test/Index", nil)
	response := httptest.NewRecorder()

	mgr.HandleRequest(response, request)
	fmt.Println(response.Body.String())
}
