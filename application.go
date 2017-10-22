/*
	Digivance MVC Application Framework
	Application object
	Dan Mayor (dmayor@digivance.com)

	This file defines functionality for mapping an action method to an http request optionally
	boud to an http verb.

	This package is released under as open source under the LGPL-3.0 which can be found:
	https://opensource.org/licenses/LGPL-3.0
*/

package mvcapp

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// Application is our global scope object (E.g. application wide configuration
// and manager systems)
type Application struct {
	// SessionManager is our HTTP Session Manager system
	SessionManager *SessionManager

	// RouteManager is our Route Manager system
	RouteManager *RouteManager

	// DomainName is used by methods that generate full links, such as redirect secure methods
	DomainName string

	// BindAddress is the ip address that this application will listen on
	BindAddress string

	// HTTPPort is the port used to stream plain text http protocol
	HTTPPort int

	// HTTPSPort is the port used to stream TLS secured http protocol
	HTTPSPort int
}

// NewApplication returns a new default MVC Application object
func NewApplication() *Application {
	rtn := &Application{
		SessionManager: NewSessionManager(),
		RouteManager:   NewRouteManager(),
		BindAddress:    "",
		HTTPPort:       80,
		HTTPSPort:      443,
	}

	// I know, it's weird... Just roll with it
	rtn.RouteManager.SessionManager = rtn.SessionManager
	return rtn
}

// Run is used to execute this MVC Application (direct http socket server)
func (app *Application) Run() error {
	addr := fmt.Sprintf("%s:%d", app.BindAddress, app.HTTPPort)
	return http.ListenAndServe(addr, http.HandlerFunc(app.RouteManager.HandleRequest))
}

// RunSecure is used to execute this MVC Application over HTTPS/TLS (direct https socket server)
func (app *Application) RunSecure(certFile string, keyFile string) error {
	addr := fmt.Sprintf("%s:%d", app.BindAddress, app.HTTPSPort)
	return http.ListenAndServeTLS(addr, certFile, keyFile, http.HandlerFunc(app.RouteManager.HandleRequest))
}

// RunForcedSecure is used to execute this MVC Application in both HTTP and
// HTTPS/TLS modes, the HTTP mode will force redirection to HTTPS only. (direct http and https
// socket servers)
func (app *Application) RunForcedSecure(certFile string, keyFile string) error {
	addr := fmt.Sprintf("%s:%d", app.BindAddress, app.HTTPPort)
	go http.ListenAndServe(addr, http.HandlerFunc(app.RedirectSecure))

	addr = fmt.Sprintf("%s:%d", app.BindAddress, app.HTTPSPort)
	return http.ListenAndServeTLS(addr, certFile, keyFile, http.HandlerFunc(app.RouteManager.HandleRequest))
}

// RunForcedSecureJS is used to run a web application in forced TLS secure mode by returning
// a simple page with javascript that redirects the browser to https://DomainName/path
func (app *Application) RunForcedSecureJS(certFile string, keyFile string) error {
	addr := fmt.Sprintf("%s:%d", app.BindAddress, app.HTTPPort)
	go http.ListenAndServe(addr, http.HandlerFunc(app.RedirectSecureJS))

	addr = fmt.Sprintf("%s:%d", app.BindAddress, app.HTTPSPort)
	return http.ListenAndServeTLS(addr, certFile, keyFile, http.HandlerFunc(app.RouteManager.HandleRequest))
}

// RedirectSecure is used to submit an http redirect from http to https when forcing secure
func (app *Application) RedirectSecure(w http.ResponseWriter, req *http.Request) {
	target := fmt.Sprintf("https://%s%s", app.DomainName, req.URL.Path)

	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}

	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

// RedirectSecureJS is used to submit a javascript redirect from http to https when forcing secure
func (app *Application) RedirectSecureJS(w http.ResponseWriter, req *http.Request) {
	data := fmt.Sprintf("<html><head><title>Redirecting to secure site mode</title></head><body>window.location.href='https://%s%s';</body>", app.DomainName, req.URL.Path)
	w.Write([]byte(data))
}

var appPath = ""

// GetApplicationPath should return the full path to the executable.
// This is the root of the site and where the assembly file is
func GetApplicationPath() string {
	if appPath == "" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			appPath = "."
		}

		appPath = dir
	}

	return appPath
}
