/*
	Digivance MVC Application Framework
	Base Controller Features
	Dan Mayor (dmayor@digivance.com)

	This file defines the base controller functionality that the caller will use to derrive
	their custom controller objects.

	This package is released under as open source under the LGPL-3.0 which can be found:
	https://opensource.org/licenses/LGPL-3.0
*/

package mvcapp

import (
	"fmt"
	"net/http"

	"github.com/Digivance/str"
)

// ControllerCallback is a simple declaration to provide a callback method
// members (e.g. variables that point to methods to be executed)
type ControllerCallback func()

type ErrorResultCallback func(err error) IActionResult
type NotFoundResultCallback func(url string) IActionResult

// IController defines the RegisterAction and Execute methods that
// need to be implemented by all controllers
type IController interface {
	RegisterAction(string, string, ActionMethod)
	Execute() IActionResult
	SetRequest(*http.Request)
	ToController() *Controller
}

// Controller contains the basic members shared by custom controllers,
// also defines the required RegisterAction and Execute methods (below)
type Controller struct {
	IController

	Request        *http.Request
	Session        *Session
	Cookies        []*http.Cookie
	HTTPStatusCode int

	RequestedPath string
	QueryString   map[string]string
	Fragment      string

	DefaultAction string
	ActionRoutes  []*ActionMap

	BeforeExecute ControllerCallback
	AfterExecute  ControllerCallback

	ErrorResult    ErrorResultCallback
	NotFoundResult NotFoundResultCallback
}

// NewBaseController returns a reference to a new Base Controller
func NewBaseController(request *http.Request) *Controller {
	rtn := &Controller{
		Request:        request,
		Session:        &Session{},
		Cookies:        make([]*http.Cookie, 0),
		HTTPStatusCode: 200,

		RequestedPath: request.URL.Path,
		QueryString:   map[string]string{},
		Fragment:      "",

		DefaultAction: "",
		ActionRoutes:  make([]*ActionMap, 0),
	}

	for _, cookie := range request.Cookies() {
		rtn.Cookies = append(rtn.Cookies, cookie)
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

// AddCookie is used to set append a new cookie to this controllers cookie collection
func (controller *Controller) AddCookie(cookie *http.Cookie) {
	controller.Cookies = append(controller.Cookies, cookie)
}

// Execute is called by the route manager instructing this controller to respond
func (controller *Controller) Execute() IActionResult {
	verb := controller.Request.Method
	actionName := controller.DefaultAction
	params := []string{}

	if str.Contains(controller.RequestedPath, "/") {
		parts := str.Split(controller.RequestedPath, '/')

		if len(parts) > 1 {
			actionName = parts[1]

			if len(parts) > 2 {
				params = parts[2:]
			}
		}
	}

	for _, actionMethod := range controller.ActionRoutes {
		if str.Compare(actionMethod.Name, actionName) && (len(actionMethod.Verb) <= 0 || str.Compare(actionMethod.Verb, verb)) {
			res := actionMethod.Method(params)
			return res
		}
	}

	return NewActionResult([]byte{})
}

// ToController is a method defined by the controller object (which implements IController) that
// returns a reference to the Controller object it is called on. We use this in the route manager
// to gain access to the session and cookie collections of the base controller from a custom controller
func (controller *Controller) ToController() *Controller {
	return controller
}

// DefaultErrorPage will attempt to render the built in error page
func (controller *Controller) DefaultErrorPage(err error) IActionResult {
	controller.HTTPStatusCode = 500
	html := fmt.Sprintf("<html><head><title>Server Error</title></head><body><h1>Server Error :(</h1>%s</body></html>", err.Error())
	data := []byte(html)
	return NewActionResult(data)
}

// DefaultNotFoundPage will attempt to render the built in 404 page
func (controller *Controller) DefaultNotFoundPage() IActionResult {
	controller.HTTPStatusCode = 404
	url := controller.RequestedPath
	html := fmt.Sprintf("<html><head><title>Content Not Found</title></head><body><h1>Content Missing</h1>We're sorry, we could not find '%s' from this app :(</body></html>", url)
	data := []byte(html)
	return NewActionResult(data)
}
