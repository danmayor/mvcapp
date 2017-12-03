/*
	Digivance MVC Application Framework
	Application object
	Dan Mayor (dmayor@digivance.com)

	This file defines functionality for mapping an action method to an http request optionally
	boud to an http verb.
*/

package mvcapp

import (
	"errors"
	"fmt"
	"net/http"
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

	// HTTPServer is the http.Server object being used to host http transport
	HTTPServer *http.Server

	// HTTPSPort is the port used to stream TLS secured http protocol
	HTTPSPort int

	// HTTPSServer is the http.Server object being used to host https transport
	HTTPSServer *http.Server
}

// NewApplication returns a new default MVC Application object
func NewApplication() *Application {
	rtn := &Application{
		SessionManager: NewSessionManager(),
		RouteManager:   NewRouteManager(),
		BindAddress:    "",
		HTTPPort:       80,
		HTTPServer:     nil,
		HTTPSPort:      443,
		HTTPSServer:    nil,
	}

	rtn.RouteManager.SessionManager = rtn.SessionManager

	if LogFilename == "" {
		SetLogFilename(fmt.Sprintf("%s/%s", GetApplicationPath(), "mvcapp.log"))
	}

	LogMessage("Application initialized")
	return rtn
}

// Stop is used to stop hosting this MVC Application. You can call one of the Run methods to restart
func (app *Application) Stop() {
	if app.HTTPServer != nil {
		app.HTTPServer.Shutdown(nil)
	}

	if app.HTTPSServer != nil {
		app.HTTPSServer.Shutdown(nil)
	}
}

// Run is used to execute this MVC Application (direct http socket server)
func (app *Application) Run() error {
	if app.HTTPServer != nil {
		return errors.New("Can not run application, HTTPServer already in use")
	}

	addr := fmt.Sprintf("%s:%d", app.BindAddress, app.HTTPPort)
	app.HTTPServer = &http.Server{Addr: addr}
	app.HTTPServer.Handler = http.HandlerFunc(app.RouteManager.HandleRequest)

	return app.HTTPServer.ListenAndServe()
}

// RunSecure is used to execute this MVC Application over HTTPS/TLS (direct https socket server)
func (app *Application) RunSecure(certFile string, keyFile string) error {
	if app.HTTPSServer != nil {
		return errors.New("Can not RunSecure, HTTPSServer already in use")
	}

	addr := fmt.Sprintf("%s:%d", app.BindAddress, app.HTTPSPort)
	app.HTTPSServer = &http.Server{Addr: addr}
	app.HTTPServer.Handler = http.HandlerFunc(app.RouteManager.HandleRequest)

	return app.HTTPSServer.ListenAndServeTLS(certFile, keyFile)
}

// RunForcedSecure is used to execute this MVC Application in both HTTP and
// HTTPS/TLS modes, the HTTP mode will force redirection to HTTPS only. (direct http and https
// socket servers)
func (app *Application) RunForcedSecure(certFile string, keyFile string) error {
	if app.HTTPServer != nil {
		return errors.New("Can not RunForcedSecure, HTTPServer already in use")
	}

	if app.HTTPSServer != nil {
		return errors.New("Can not RunForcedSecure, HTTPSServer already in use")
	}

	addr := fmt.Sprintf("%s:%d", app.BindAddress, app.HTTPPort)
	app.HTTPServer = &http.Server{Addr: addr}
	app.HTTPServer.Handler = http.HandlerFunc(app.RedirectSecure)
	go app.HTTPServer.ListenAndServe()

	addr = fmt.Sprintf("%s:%d", app.BindAddress, app.HTTPSPort)
	app.HTTPSServer = &http.Server{Addr: addr}
	app.HTTPSServer.Handler = http.HandlerFunc(app.RouteManager.HandleRequest)
	return app.HTTPSServer.ListenAndServeTLS(certFile, keyFile)
}

// RunForcedSecureJS is used to run a web application in forced TLS secure mode by returning
// a simple page with javascript that redirects the browser to https://DomainName/path
func (app *Application) RunForcedSecureJS(certFile string, keyFile string) error {
	if app.HTTPServer != nil {
		return errors.New("Can not RunForcedSecureJS, HTTPServer already in use")
	}

	if app.HTTPSServer != nil {
		return errors.New("Can not RunForcedSecure, HTTPSServer already in use")
	}

	addr := fmt.Sprintf("%s:%d", app.BindAddress, app.HTTPPort)
	app.HTTPServer = &http.Server{Addr: addr}
	app.HTTPServer.Handler = http.HandlerFunc(app.RedirectSecureJS)
	go app.HTTPServer.ListenAndServe()

	addr = fmt.Sprintf("%s:%d", app.BindAddress, app.HTTPSPort)
	app.HTTPSServer = &http.Server{Addr: addr}
	app.HTTPSServer.Handler = http.HandlerFunc(app.RouteManager.HandleRequest)
	return app.HTTPSServer.ListenAndServeTLS(certFile, keyFile)
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
	data := fmt.Sprintf("<html><head><title>Redirecting to secure site mode</title></head><body><script type=\"text/javascript\">window.location.href='https://%s%s';</script></body>", app.DomainName, req.URL.Path)
	w.Write([]byte(data))
}
