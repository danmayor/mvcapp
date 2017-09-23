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

// Application is our global scope object (E.g. application wide configuration
// and manager systems)
type Application struct {
	// SessionKey is the Cookie Name used to store a connections
	// session ID
	SessionKey string

	// Sessions is our HTTP Session Manager system
	Sessions *SessionManager

	// Routes is our Route Manager system
	Routes *RouteManager
}

// NewApplication returns a new default MVC Application object
func NewApplication() Application {
	return Application{
		SessionKey: "SessionID",
		Routes:     NewRouteManager(),
	}
}
