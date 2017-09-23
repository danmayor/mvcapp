/*
	Digivance MVC Application Framework
	Route Manager Unit Tests
	Dan Mayor (dmayor@digivance.com)

	This file defines the unit tests for the Route Manager features.

	This package is released under as open source under the LGPL-3.0 which can be found:
	https://opensource.org/licenses/LGPL-3.0
*/

package mvcapp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Digivance/mvcapp"
)

// TestModel is used to test passing a data model to our view result / templates
type TestModel struct {
	Title   string
	Welcome string
}

// TestController is used to test basic custom controller coding
type TestController struct {
	*mvcapp.Controller
}

// NewTestController is used to test our routing controller creator method
func NewTestController(request *http.Request) mvcapp.IController {
	rtn := &TestController{
		Controller: mvcapp.NewBaseController(request),
	}

	rtn.RegisterAction("GET", "Index", rtn.Index)
	return rtn
}

// Index is used to test our basic action methods
func (controller TestController) Index(params []string) mvcapp.IActionResult {
	if controller.Session.Values != nil {
		// Here we test using the controllers' browser session
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
	}

	// Here we test setting cookies to be passed to the browser
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

// TestRouteManager is our unit test method, it makes 2 requests to test passing cookies
// and sessions, as well as testing the model rendering in view results.
func TestRouteManager(t *testing.T) {
	mgr := mvcapp.NewRouteManager()
	mgr.SessionManager = mvcapp.NewSessionManager()
	mgr.RegisterController("Test", NewTestController)

	request, err := http.NewRequest("GET", "/Test/Index", nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	request.AddCookie(&http.Cookie{
		Name:    mgr.SessionIDKey,
		Value:   "EIMA5VQOU4980S35AYPAEKYABL73GZBA",
		Expires: time.Now().Add(15 * time.Minute),
	})

	response := httptest.NewRecorder()

	mgr.HandleRequest(response, request)
	mgr.HandleRequest(response, request)
}
