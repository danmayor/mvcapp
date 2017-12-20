/*
	Digivance MVC Application Framework
	Route Manager Features
	Dan Mayor (dmayor@digivance.com)

	This file defines the generic route manager functionality of the MVC application. This manager
	allows the caller to register route maps and bind the handler method. This system drives the
	request pipeline of an MVC application made with this framework.
*/

package mvcapp

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// ControllerCreator is a delegate to the creation method of a controller
// that is mapped as the primary route. (E.g. site.com/CONTROLLER is mapped
// to a NewXController method that implements this signature)
type ControllerCreator func(*http.Request) IController

// RouteManager provides the basic http request pipeline of the
// mvcapp framework
type RouteManager struct {
	// SessionIDKey is the name of the Cookie / Session key to use when identifying
	// the browser session ID. (E.g. name of the cookie that contains this users
	// browser session ID)
	SessionIDKey string

	// DefaultController is a string defining the name of the controller to execute
	// when a request comes in to the root of the site (Should be your home /
	// site index controller)
	DefaultController string

	// DefaultAction is a string defining the name of the action method to be called
	// when a request is made to the root of a controller. (This should be your home
	// / default or index page name)
	DefaultAction string

	// Routes is the collection of RouteMaps that define the controllers which are
	// registered in this manager
	Routes []*RouteMap

	// SessionManager is a pointer to the SessionManager object to use for this app
	SessionManager *SessionManager

	// BundleManager is a pointer to the BundleManager that can be used by controllers
	// that derrive from the BundleController type (Is set during execution pipeline)
	BundleManager *BundleManager
}

// NewRouteManager returns a new route manager object with default
// controller and action tokens set to "Home" and "Index".
func NewRouteManager() *RouteManager {
	return &RouteManager{
		SessionIDKey: "SessionID",

		DefaultController: "Home",
		DefaultAction:     "Index",

		Routes:         make([]*RouteMap, 0),
		SessionManager: NewSessionManager(),
	}
}

// toQuerystringMap will parse the provided url encoded query string into a map of kvp's
func (manager *RouteManager) toQueryStringMap(queryString string) map[string]string {
	rtn := map[string]string{}

	pairs := strings.Split(queryString, "&")
	for _, pair := range pairs {
		kvp := strings.Split(pair, "=")
		if len(kvp) >= 2 {
			k := kvp[0]
			v := strings.Join(kvp[1:], "=")

			rtn[k] = strings.TrimRight(v, "=")
		}
	}

	return rtn
}

// parseControllerName returns the controller name requested, will fallback and return
// the default controller if this is a root request.
func (manager *RouteManager) parseControllerName(path string) string {
	rtn := manager.DefaultController
	parts := strings.Split(strings.TrimLeft(path, "/"), "/")

	if len(parts) > 0 && parts[0] != "" {
		rtn = parts[0]
	}

	return rtn
}

// getController takes the response and request from our http server and map it to the
// registered icontroller and controller objects (if they exist)
func (manager *RouteManager) getController(response http.ResponseWriter, request *http.Request) (IController, *Controller) {
	path := strings.TrimLeft(request.URL.Path, "/")
	controllerName := manager.parseControllerName(path)

	TraceLog(fmt.Sprintf("Getting controller request to controller: %s", controllerName))

	for _, route := range manager.Routes {
		if strings.HasPrefix(strings.ToLower(route.ControllerName), strings.ToLower(controllerName)) {
			// Construct the appropriate controller
			icontroller := route.CreateController(request)
			controller := icontroller.ToController()

			controller.ControllerName = controllerName
			controller.Response = response
			controller.DefaultAction = manager.DefaultAction
			controller.RequestedPath = path
			controller.QueryString = manager.toQueryStringMap(request.URL.RawQuery)
			controller.Fragment = request.URL.Fragment
			controller.Cookies = request.Cookies()

			TraceLog(fmt.Sprintf("Constructed controller: %s", controllerName))
			return icontroller, controller
		}
	}

	TraceLog(fmt.Sprintf("Failed to obtain controller for request to: %s", controllerName))
	return nil, nil
}

// setControllerSessions is called if there is an active session manager. This method will
// read the browser cookies to find the browser session ID (as defined by the managers SessionIDKey)
// and if present, will load the browser session value collection for this user into the controllers
// Session member.
func (manager *RouteManager) setControllerSessions(controller *Controller) {
	if controller == nil {
		LogError("Can not set controller sessions, no controller registered")
		return
	}

	if controller.Request == nil {
		LogError("Can not set controller sessions, no request received?")
		return
	}

	// Get the browser session ID from the request cookies
	browserSessionCookie, err := controller.Request.Cookie(manager.SessionIDKey)
	browserSessionID := ""
	if err != nil || browserSessionCookie == nil || len(browserSessionCookie.Value) < 32 || !manager.SessionManager.Contains(browserSessionCookie.Value) {
		browserSessionID = RandomString(32)
	} else {
		browserSessionID = browserSessionCookie.Value
	}

	// Get the browserSession from the SessionManager and set
	// the controllers reference to it
	browserSession := manager.SessionManager.GetSession(browserSessionID)
	if browserSession == nil {
		browserSession = manager.SessionManager.CreateSession(browserSessionID)
	}

	controller.Session = browserSession
	controller.Session.ActivityDate = time.Now()
	controller.SetCookie(&http.Cookie{Name: manager.SessionIDKey, Value: browserSessionID, Path: "/"})
}

// handleFile is called if HandleRequest fails to load the controller or the result, if this fails
// we will fall back on MVC 404 functionality
func (manager *RouteManager) handleFile(response http.ResponseWriter, request *http.Request) bool {
	path := request.URL.Path
	if strings.HasPrefix(strings.ToLower(path), "/") {
		path = fmt.Sprintf("%s/%s", GetApplicationPath(), path[1:])
	}

	if path == "" {
		return false
	}

	f, err := os.Stat(path)
	if os.IsNotExist(err) {
		LogWarning(fmt.Sprintf("404 Trying to serve raw file: %s", path))
		return false
	}

	// refuse to serve directory contents for security
	mode := f.Mode()
	if mode.IsDir() {
		LogWarning(fmt.Sprintf("User tried to request raw directory contents and was blocked: %s", path))
		return false
	}

	if manager.validPath(path) {
		LogMessage(fmt.Sprintf("Serving raw file: %s", path))
		http.ServeFile(response, request, path)
		return true
	}

	LogError(fmt.Sprintf("Unknown error serving file [%s], permissions problem?", path))
	return false
}

// validPath is used internally to ignore paths that are used by the mvcapp system
func (manager *RouteManager) validPath(path string) bool {
	if strings.HasPrefix(strings.ToLower(path), "controllers/") {
		return false
	}

	if strings.HasPrefix(strings.ToLower(path), "emails/") {
		return false
	}

	if strings.HasPrefix(strings.ToLower(path), "models/") {
		return false
	}

	if strings.HasPrefix(strings.ToLower(path), "views/") {
		return false
	}

	return true
}

// RegisterController is used to map a custom controller object to the
// controller section of the requested url (E.g. "site.com/CONTROLLER/action")
func (manager *RouteManager) RegisterController(name string, creator ControllerCreator) {
	LogMessage(fmt.Sprintf("Registering controller for: %s", name))
	manager.Routes = append(manager.Routes, NewRouteMap(name, creator))
}

// HandleRequest is mapped to the http handler method and processes the
// HTTP request pipeline
func (manager *RouteManager) HandleRequest(response http.ResponseWriter, request *http.Request) {
	TraceLog(fmt.Sprintf("Handling request: %s", request.URL.String()))

	// Gets the controller objects responsible for this route (if they exist)
	icontroller, controller := manager.getController(response, request)

	path := strings.TrimLeft(request.URL.Path, "/")
	if !manager.validPath(path) {
		// If the path is invalid, we use the default controller to render an
		// error page telling the user so
		request, _ = http.NewRequest("GET", manager.DefaultController, nil)
		icontroller, controller = manager.getController(response, request)

		if icontroller == nil || controller == nil {
			LogError("Failed to load default controller to serve invalid path error page")
			return
		}

		if controller.ErrorResult != nil {
			controller.ErrorResult(errors.New("Invalid path requested")).Execute(response)
		} else {
			controller.DefaultErrorPage(errors.New("Invalid path requested")).Execute(response)
		}

		LogWarning(fmt.Sprintf("Request to invalid path: %s", request.URL.String()))
		return
	}

	// If the controller is nil lets try to serve a raw file
	if controller == nil {
		if manager.handleFile(response, request) {
			return
		}

		request, _ = http.NewRequest("GET", manager.DefaultController, nil)
		icontroller, controller = manager.getController(response, request)
	}

	if controller == nil {
		LogError("Critical failure, could not load controller by request or the default controller!")
		response.WriteHeader(500)
		response.Write([]byte("Failed to handle request, please try again"))
		return
	}

	// If the route manager has a session manager, we'll fire that bad boy up
	// and try to get the browser session id from the submitted cookies
	// which is then loaded into the controller session value collection
	if manager.SessionManager != nil {
		manager.setControllerSessions(controller)
	}

	// Call our before execute callback if one is registered
	if controller.BeforeExecute != nil {
		controller.BeforeExecute()
	}

	// If our before execute needs to fail, it can do so and set continue pipeline
	// to false, which means we should not attempt to execute the controller.
	if controller.ContinuePipeline {
		result := icontroller.Execute()
		if result == nil {
			if !manager.handleFile(response, request) {
				if controller.NotFoundResult != nil {
					result = controller.NotFoundResult()
				} else {
					result = controller.DefaultNotFoundPage()
				}
			}
		}

		icontroller.WriteResponse(result)
	}

	// Regardless of executing the controller or not, we call the after execute callback
	// if it exists
	if controller.AfterExecute != nil {
		controller.AfterExecute()
	}
}
