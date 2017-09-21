package mvcapp

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestModel struct {
	Title   string
	Welcome string
}

type TestController struct {
	*Controller
}

func NewTestController(request *http.Request) IController {
	rtn := &TestController{
		Controller: NewBaseController(request),
	}
	rtn.RegisterAction("GET", "Index", rtn.Index)
	return rtn
}

func (controller TestController) Index(params []string) IActionResult {
	controller.Cookies = append(controller.Cookies, &http.Cookie{
		Name:   "Dan",
		Value:  "is awesome!",
		MaxAge: 900,
	})

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
	mgr := NewRouteManager()
	mgr.RegisterController("Test", NewTestController)

	request := httptest.NewRequest("GET", "/Test/Index", nil)
	response := httptest.NewRecorder()

	mgr.HandleRequest(response, request)

	/*
		fmt.Println(response.Body.String())

		res := http.Response{Header: response.Header()}
		cookies := res.Cookies()
		name := cookies[0].Name
		value := cookies[0].Value

		fmt.Println(fmt.Sprintf("[%s] = \"%s\"\n", name, value))
	*/
}
