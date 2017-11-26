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
}

// NewRouteManager returns a new route manager object with default
// controller and action tokens set to "Home" and "Index".
func NewRouteManager() *RouteManager {
	return &RouteManager{
		SessionIDKey: "MvcApp.SessionID",

		DefaultController: "Home",
		DefaultAction:     "Index",

		Routes:         make([]*RouteMap, 0),
		SessionManager: NewSessionManager(),
	}
}

// RegisterController is used to map a custom controller object to the
// controller section of the requested url (E.g. "site.com/CONTROLLER/action")
func (manager *RouteManager) RegisterController(name string, creator ControllerCreator) {
	manager.Routes = append(manager.Routes, NewRouteMap(name, creator))
}

// HandleRequest is mapped to the http handler method and processes the
// HTTP request pipeline
func (manager *RouteManager) HandleRequest(response http.ResponseWriter, request *http.Request) {
	// Gets the controller objects responsible for this route (if they exist)
	icontroller, controller := manager.getController(response, request)

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
		} else {
			icontroller.WriteResponse(result)
		}
	}

	// Regardless of executing the controller or not, we call the after execute callback
	// if it exists
	if controller.AfterExecute != nil {
		controller.AfterExecute()
	}
}

// parseFragment will extract, and remove the fragment (or named anchor) section
// of the url. Returns strings representing the fragment and the url without it
func (manager *RouteManager) parseFragment(url string) (string, string) {
	fragment := ""

	if strings.Contains(url, "#") {
		fragment = url[strings.Index(url, "#"):]
		url = url[0 : strings.Index(url, "#")-1]
	}

	return fragment, url
}

// parseQueryString will extract the path and parse the query string key value pairs.
// Returns the path (relative to the app's domain) and a map of the query string pairs.
func (manager *RouteManager) parseQueryString(url string) (string, map[string]string) {
	path := ""
	queryString := map[string]string{}

	if strings.Contains(url, "?") {
		path = url[0 : strings.Index(url, "?")-1]
		qsLine := url[strings.Index(url, "?"):]

		for _, pair := range strings.Split(qsLine, "&") {
			kvp := strings.Split(pair, "=")
			if len(kvp) > 1 {
				queryString[kvp[0]] = kvp[1]
			}

			if len(kvp) == 1 {
				queryString[kvp[0]] = ""
			}
		}
	} else {
		path = url
	}

	return path, queryString
}

// parseControllerName returns the controller name requested, will fallback and return
// the default controller if this is a root request.
func (manager *RouteManager) parseControllerName(path string) string {
	rtn := manager.DefaultController
	parts := strings.Split(path, "/")

	if len(parts) > 0 && parts[0] != "" {
		rtn = parts[0]
	}

	return rtn
}

// getController takes the response and request from our http server and map it to the
// registered icontroller and controller objects (if they exist)
func (manager *RouteManager) getController(response http.ResponseWriter, request *http.Request) (IController, *Controller) {
	fragment, url := manager.parseFragment(request.URL.Path)
	path, queryString := manager.parseQueryString(url)
	controllerName := manager.parseControllerName(path)

	for _, route := range manager.Routes {
		if strings.HasPrefix(strings.ToLower(route.ControllerName), strings.ToLower(controllerName)) {
			// Construct the appropriate controller
			icontroller := route.CreateController(request)
			controller := icontroller.ToController()

			controller.Response = response
			controller.DefaultAction = manager.DefaultAction
			controller.RequestedPath = path
			controller.QueryString = queryString
			controller.Fragment = fragment

			controller.Cookies = request.Cookies()

			return icontroller, controller
		}
	}

	return nil, nil
}

// setControllerSessions is called if there is an active session manager. This method will
// read the browser cookies to find the browser session ID (as defined by the managers SessionIDKey)
// and if present, will load the browser session value collection for this user into the controllers
// Session member.
func (manager *RouteManager) setControllerSessions(controller *Controller) {
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
	controller.Session = browserSession
	controller.Session.ActivityDate = time.Now()
	controller.SetCookie(&http.Cookie{Name: manager.SessionIDKey, Value: browserSessionID, Path: "/"})
}

// handleFile is called if HandleRequest fails to load the controller or the result, if this fails
// we will fall back on MVC 404 functionality
func (manager *RouteManager) handleFile(response http.ResponseWriter, request *http.Request) bool {
	// TODO: Add a case insensitive file finding thingie here (because linux kind of sucks like that)
	_, url := manager.parseFragment(request.URL.Path)
	path, _ := manager.parseQueryString(url)

	if strings.HasPrefix(strings.ToLower(path), "/") {
		path = fmt.Sprintf("%s/%s", GetApplicationPath(), path[1:])
	}

	if path == "" {
		return false
	}

	f, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	// refuse to serve directory contents for security
	mode := f.Mode()
	if mode.IsDir() {
		return false
	}

	if validPath(path) {
		http.ServeFile(response, request, path)
		return true
	}

	return false
}

// validPath is used internally to ignore paths that are used by the mvcapp system
func validPath(path string) bool {
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

	if !strings.Contains(strings.ToLower(path), "/") {
		return false
	}

	return true
}
