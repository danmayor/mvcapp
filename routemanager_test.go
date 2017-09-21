package mvcapp_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Digivance/mvcapp"
)

type TestModel struct {
	Title   string
	Welcome string
}

type TestController struct {
	*mvcapp.Controller
}

func NewTestController(request *http.Request) mvcapp.IController {
	rtn := &TestController{
		Controller: mvcapp.NewBaseController(request),
	}

	rtn.RegisterAction("GET", "Index", rtn.Index)
	return rtn
}

func (controller TestController) Index(params []string) mvcapp.IActionResult {
	sid := controller.Session.ID
	fmt.Println(sid)
	var saidHello bool
	if controller.Session.Values["SaidHello"] != nil {
		saidHello = controller.Session.Values["SaidHello"].(bool)
	} else {
		saidHello = false
	}

	if !saidHello {
		saidHello = true
		controller.Session.Values["SaidHello"] = saidHello
	}

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

	result := mvcapp.NewViewResult(templates, model)
	result.AddHeader("Framework", "Digivance MVC Application Framework")

	return result
}

func TestRouteManager(t *testing.T) {
	mgr := mvcapp.NewRouteManager()
	mgr.RegisterController("Test", NewTestController)

	request := httptest.NewRequest("GET", "/Test/Index", nil)
	response := httptest.NewRecorder()

	mgr.HandleRequest(response, request)
	mgr.HandleRequest(response, request)

	fmt.Println(response.Body.String())

	res := http.Response{Header: response.Header()}
	cookies := res.Cookies()
	name := cookies[0].Name
	value := cookies[0].Value

	fmt.Println(fmt.Sprintf("[%s] = \"%s\"\n", name, value))
}
