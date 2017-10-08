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
)

// Application is our global scope object (E.g. application wide configuration
// and manager systems)
type Application struct {
	// SessionManager is our HTTP Session Manager system
	SessionManager *SessionManager

	// RouteManager is our Route Manager system
	RouteManager *RouteManager

	BindAddress string
	HTTPPort    int
	HTTPSPort   int
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
	go http.ListenAndServe(addr, http.HandlerFunc(redirect))

	addr = fmt.Sprintf("%s:%d", app.BindAddress, app.HTTPSPort)
	return http.ListenAndServeTLS(addr, certFile, keyFile, http.HandlerFunc(app.RouteManager.HandleRequest))
}

// redirect is used internally to submit an http redirect from http to https when
// forcing secure
func redirect(w http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.URL.Path

	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}

	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}
