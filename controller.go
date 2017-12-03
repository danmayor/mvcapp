/*
	Digivance MVC Application Framework
	Base Controller Features
	Dan Mayor (dmayor@digivance.com)

	This file defines the base controller functionality that the caller will use to derrive
	their custom controller objects.
*/

package mvcapp

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// ControllerCallback is a simple declaration to provide a callback method
// members (e.g. variables that point to methods to be executed)
type ControllerCallback func()

// ErrorResultCallback is a simple declaration to provide a callback method
// used when there is an internal server error (such as custom error page)
type ErrorResultCallback func(err error) *ActionResult

// NotFoundResultCallback is a simple declaration to provide a callback method
// used when the requested content can not be found (custom 404)
type NotFoundResultCallback func() *ActionResult

// IController defines the RegisterAction and Execute methods that
// need to be implemented by all controllers
type IController interface {
	// RegisterAction should add the Verb, Controller and Action method names to the
	// base controllers action map
	RegisterAction(string, string, ActionMethod)

	// Execute should query and execute the mapped action method
	Execute() *ActionResult

	// WriteResponse should write the provided action result to the response stream
	WriteResponse(*ActionResult)

	// RedirectJS is a method that should write an html page with javascript redirect
	// function directly to the response stream (this can be used to lock pages from
	// users based of custom conditions, such as logged in or not)
	RedirectJS(string)

	// ToController should return the actual base controller object (so that the system
	// can interact with controller variable members)
	ToController() *Controller
}

// Controller contains the basic members shared by custom controllers,
// also defines the required RegisterAction and Execute methods (below)
type Controller struct {
	// IController is the interface that the base controller implements
	IController

	// Request is the http.Request object that this controller is responding to
	Request *http.Request

	// Response is the response writer stream to write data to the client
	Response http.ResponseWriter

	// Session is the User browser session data collection for the user who made this request
	Session *Session

	// Cookies are populated from the collection submitted from the client. Server can alter
	// or add cookies to this collection to have them delivered back to the client. (Call the
	// controllers DeleteCookie method to signal the client to forget a cookie)
	Cookies []*http.Cookie

	// ContinuePipeline is used from methods that mean to prevent the execution of a controller
	// action method. E.g. if a controller signals to RedirectJS in the BeforeExecute callback,
	// this will be set to false and will prevent the action method from assembling and executing
	// the action result.
	ContinuePipeline bool

	// ControllerName is used to help the template path mapping of the View method of this controller.
	// Should be set in the controller creator method and should represent the folder name where the
	// views for this controller can be found
	ControllerName string

	// RequestedPath is a quick reference member that contains the url path (after the domain) that
	// was requested.
	RequestedPath string

	// QueryString is a quick reference member that contains the key value pair query string parameters
	// that were submitted with this request
	QueryString map[string]string

	// Fragment represents the #Value portion of the requested URL (Note some sources indicate that URL
	// fragments are or should be depreciated. Use URL fragments and this member at your discretion)
	Fragment string

	// DefaultAction is the name of the default action method name (in the route map) to try and execute
	// when making a request directly to the controller.  (E.g. site.com/controller) (This should be your
	// Index, Home or Default page name)
	DefaultAction string

	// ActionRoutes stores the routes which have been registered for this controller. This collection is
	// used in the Execute method to find the appropriate action method function to call
	ActionRoutes []*ActionMap

	// LastErrorMessage allows you to pass a simple string indicating an error or page failure. Your templates
	// should be aware of this member, and if present should display it as an error message to your user.
	LastErrorMessage string

	// LastErrorDetails allows for extended error message details (Example: if the LastErrorMessage was caused
	// due to failing to validate an email address, these details could contain the site rules for valid emails
	// and the template that see's these details set could then display said rules to the user)
	LastErrorDetails string

	// BeforeExecute is a callback method that a controller can set to provide a global method called before
	// the action method is executed. (Controller global prep function)
	BeforeExecute ControllerCallback

	// AfterExecute is a callback method that a controller can set to provide a global method called after
	// the action method is executed. (Controller global clean up function)
	AfterExecute ControllerCallback

	// ErrorResult is a callback method that a controller can set to respond to error conditions with a custom
	// error page
	ErrorResult ErrorResultCallback

	// NotFoundResult is a callback method that a controller can set to respond to a content not found condition
	// with a custom 404 page
	NotFoundResult NotFoundResultCallback
}

// NewBaseController returns a reference to a new Base Controller
func NewBaseController(request *http.Request) *Controller {
	controllerName := strings.Split(request.URL.Path, "/")[0]
	if controllerName == "" {
		controllerName = "Home"
	}

	rtn := &Controller{
		Request:          request,
		Session:          NewSession(),
		Cookies:          make([]*http.Cookie, 0),
		ContinuePipeline: true,

		ControllerName: controllerName,
		RequestedPath:  request.URL.Path,
		QueryString:    map[string]string{},
		Fragment:       "",

		DefaultAction: "",
		ActionRoutes:  make([]*ActionMap, 0),

		LastErrorMessage: "",
		LastErrorDetails: "",
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

// GetCookie returns the requested cookie from this controllers collection
func (controller *Controller) GetCookie(name string) *http.Cookie {
	for _, v := range controller.Cookies {
		if strings.EqualFold(v.Name, name) {
			return v
		}
	}

	return nil
}

// SetCookie will overwrite or create a cookie in this controllers collection
func (controller *Controller) SetCookie(cookie *http.Cookie) {
	for k, v := range controller.Cookies {
		if strings.EqualFold(v.Name, cookie.Name) {
			controller.Cookies[k] = cookie
			return
		}
	}

	controller.Cookies = append(controller.Cookies, cookie)
}

// DeleteCookie will set the cookie (identified by provided ccokieName) to expire in
// the past, thus making the browser remove it and stop sending it back.
func (controller *Controller) DeleteCookie(cookieName string) {
	for _, v := range controller.Cookies {
		if strings.EqualFold(v.Name, cookieName) {
			v.Expires = time.Now().Add(-1 * time.Hour)
		}
	}
}

// Execute is called by the route manager instructing this controller to respond
func (controller *Controller) Execute() *ActionResult {
	verb := controller.Request.Method
	actionName := controller.DefaultAction
	params := []string{}

	if strings.Contains(strings.ToLower(controller.RequestedPath), "/") && controller.RequestedPath != "/" {
		// Strips the leading / so we prevent the empty first parts element below
		url := controller.RequestedPath
		if strings.HasPrefix(url, "/") {
			url = url[1:]
		}

		parts := strings.Split(url, "/")

		if len(parts) > 1 {
			actionName = parts[1]

			if len(parts) > 2 {
				params = parts[2:]
			}
		}
	}

	for _, actionMethod := range controller.ActionRoutes {
		if strings.EqualFold(actionMethod.Name, actionName) && (len(actionMethod.Verb) <= 0 || strings.EqualFold(actionMethod.Verb, verb)) {
			if strings.EqualFold(verb, "POST") {
				controller.Request.ParseForm()
			}

			res := actionMethod.Method(params)
			return res
		}
	}

	if controller.NotFoundResult != nil {
		return controller.NotFoundResult()
	}

	return controller.DefaultNotFoundPage()
}

// WriteResponse is called from the route manager to execute the result that was constructed
// from this controllers Execute method (E.g. the result returned from the action if mapped)
func (controller *Controller) WriteResponse(result *ActionResult) {
	if controller.ContinuePipeline {
		if result == nil || len(result.Data) <= 0 {
			if controller.NotFoundResult != nil {
				result = controller.NotFoundResult()
			} else {
				result = controller.DefaultNotFoundPage()
			}
		}

		result.Execute(controller.Response)
	}
}

// RedirectJS is a helper method that will write a very simple html page using the
// window.location.href='url' method to redirect the borwser to the provided url
// Note this will also set the controller.ContinuePipeline to false, meaning that
// the ActionMethod for this request will NOT be called. This allows us to use this
// method from BeginExecute callbacks to lock down an entire controller to given
// conditions, such as if the user is logged in. Can be called anytime before
// AfterExecute.
func (controller *Controller) RedirectJS(url string) {
	data := fmt.Sprintf("<html><head><title>Redirecting...</title><body><script type=\"text/javascript\">window.location.href='%s';</script></body></html>", url)
	TraceLog(fmt.Sprintf("Redirecting user with javascript, payload to follow:\n%s", data))

	// We manually write the cookies to the browser here because we'll be breaking the
	// standard pipelint (eg ContinuePipeline = false)
	res := NewActionResult([]byte(data))
	res.StatusCode = 200
	res.Cookies = controller.Cookies
	res.Headers["Cache-Control"] = "no-cache, no-store, must-revalidate"
	res.Headers["Pragma"] = "no-cache"
	res.Headers["Expires"] = "0"

	TraceLog("Payload and headers set for redirection via javascript, submitting response.")
	res.Execute(controller.Response)
	controller.ContinuePipeline = false
}

// Result returns a new ActionResult and automatically assigns the controllers cookies
func (controller *Controller) Result(data []byte) *ActionResult {
	res := NewActionResult(data)
	res.Cookies = controller.Cookies
	return res
}

// View will take the provided array of template names and try to make an mvcapp Template
// List (see func mvcapp.MakeTemplateList) using the type name of this controller (from
// reflection). Then returns the ViewResult that is created.
func (controller *Controller) View(templates []string, model interface{}) *ActionResult {
	templateList := MakeTemplateList(strings.ToLower(controller.ControllerName), templates)
	res := NewViewResult(templateList, model)
	if res == nil {
		if controller.ErrorResult != nil {
			return controller.ErrorResult(errors.New("Internal server error, failed to render page"))
		}

		return controller.DefaultErrorPage(errors.New("Internal server error, failed to render page"))
	}

	res.Cookies = controller.Cookies
	return res
}

// JSON returns a new JSONResult object of the provided payload
func (controller *Controller) JSON(payload interface{}) *ActionResult {
	res := NewJSONResult(payload)
	if res == nil {
		return NewJSONResult(false)
	}

	res.Cookies = controller.Cookies
	return res
}

// ToController is a method defined by the controller object (which implements IController) that
// returns a reference to the Controller object it is called on. We use this in the route manager
// to gain access to the session and cookie collections of the base controller from a custom controller
func (controller *Controller) ToController() *Controller {
	return controller
}

// DefaultErrorPage will attempt to render the built in error page
func (controller *Controller) DefaultErrorPage(err error) *ActionResult {
	html := fmt.Sprintf("<html><head><title>Server Error</title></head><body><h1>Server Error :(</h1>%s</body></html>", err.Error())
	data := []byte(html)
	res := NewActionResult(data)
	res.Cookies = controller.Cookies
	res.StatusCode = 500

	return res
}

// DefaultNotFoundPage will attempt to render the built in 404 page
func (controller *Controller) DefaultNotFoundPage() *ActionResult {
	html := fmt.Sprintf("<html><head><title>Content Not Found</title></head><body><h1>Content Missing</h1>We're sorry, we could not find '%s' from this app :(</body></html>", controller.RequestedPath)
	data := []byte(html)
	res := NewActionResult(data)
	res.Cookies = controller.Cookies
	res.StatusCode = 404

	return res
}
