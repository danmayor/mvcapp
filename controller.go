/*
	Digivance MVC Application Framework
	Base Controller Feature
	Dan Mayor (dmayor@digivance.com)

	This file defines the base controller system functionality.
*/

package mvcapp

import (
	"net/http"

	"github.com/Digivance/str"
)

// IController defines the RegisterAction and Execute methods that
// need to be implemented by all controllers
type IController interface {
	RegisterAction(string, string, ActionMethod)
	Execute(string, []string) IActionResult
	SetRequest(*http.Request)
}

// Controller contains the basic members shared by custom controllers,
// also defines the required RegisterAction and Execute methods (below)
type Controller struct {
	IController

	Session      *Session
	Cookies      []*http.Cookie
	Request      *http.Request
	ActionRoutes []*ActionMap
}

func NewBaseController(request *http.Request) *Controller {
	rtn := &Controller{
		Session:      &Session{},
		Cookies:      make([]*http.Cookie, 0),
		Request:      request,
		ActionRoutes: make([]*ActionMap, 0),
	}

	for i, cookie := range request.Cookies() {
		rtn.Cookies[i] = cookie
	}

	return rtn
}

// RegisterAction allows package caller to map a controller action method to
// a given Http Request verb and action name (E.g. site.com/Controller/ActionName)
func (controller *Controller) RegisterAction(verb string, name string, method ActionMethod) {
	controller.ActionRoutes = append(controller.ActionRoutes, NewActionMap(verb, name, method))
}

// SetRequest is used to set the http.Request reference
func (controller *Controller) SetRequest(request *http.Request) {
	controller.Request = request
}

func (controller *Controller) SetCookie(cookie *http.Cookie) {
	controller.Cookies = append(controller.Cookies, cookie)
}

// Execute the route manager to call the action method mapped to the provided
// actionName. params is the remainder of the url split by / represented as strings
func (controller *Controller) Execute(actionName string, params []string) IActionResult {
	verb := controller.Request.Method

	for _, actionMethod := range controller.ActionRoutes {
		if str.Compare(actionMethod.Name, actionName) && (len(actionMethod.Verb) <= 0 || str.Compare(actionMethod.Verb, verb)) {
			res := actionMethod.Method(params)
			return res
		}
	}

	return NewActionResult([]byte{})
}
