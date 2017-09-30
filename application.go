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
	"net/http"
)

// Application is our global scope object (E.g. application wide configuration
// and manager systems)
type Application struct {
	// SessionManager is our HTTP Session Manager system
	SessionManager *SessionManager

	// RouteManager is our Route Manager system
	RouteManager *RouteManager
}

// NewApplication returns a new default MVC Application object
func NewApplication() *Application {
	rtn := &Application{
		SessionManager: NewSessionManager(),
		RouteManager:   NewRouteManager(),
	}

	// I know, it's weird... Just roll with it
	rtn.RouteManager.SessionManager = rtn.SessionManager
	return rtn
}

// Run is used to execute this MVC Application
func (app *Application) Run() error {
	return http.ListenAndServe(":80", http.HandlerFunc(app.RouteManager.HandleRequest))
}

// RunSecure is used to execute this MVC Application over HTTPS/TLS
func (app *Application) RunSecure(certFile string, keyFile string) error {
	return http.ListenAndServeTLS(":443", certFile, keyFile, http.HandlerFunc(app.RouteManager.HandleRequest))
}

// RunForcedSecure is used to execute this MVC Application in both HTTP and
// HTTPS/TLS modes, the HTTP mode will force redirection to HTTPS only.
func (app *Application) RunForcedSecure(certFile string, keyFile string) error {
	go http.ListenAndServe(":80", http.HandlerFunc(redirect))
	return http.ListenAndServeTLS(":443", certFile, keyFile, http.HandlerFunc(app.RouteManager.HandleRequest))
}

// redirect is used internally to submit an http redirect from http to https
func redirect(w http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.URL.Path

	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}

	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}
