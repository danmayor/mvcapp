/*
	Digivance MVC Application Framework
	Action Map Features
	Dan Mayor (dmayor@digivance.com)

	This file defines functionality for mapping an action method to an http request optionally
	boud to an http verb.

	This package is released under as open source under the LGPL-3.0 which can be found:
	https://opensource.org/licenses/LGPL-3.0
*/

package mvcapp

// ActionMethod defines the method signature for controller action methods
type ActionMethod func([]string) IActionResult

// ActionMap is used to define the HTTP Verb, Controller's Action Name
// and the corresponding action method
type ActionMap struct {
	Verb   string
	Name   string
	Method ActionMethod
}

// NewActionMap returns a new ActionMap struct populated with the given parameters
func NewActionMap(httpVerb string, actionName string, actionMethod ActionMethod) *ActionMap {
	return &ActionMap{
		Verb:   httpVerb,
		Name:   actionName,
		Method: actionMethod,
	}
}

// NewGetActionMap returns a new ActionMap struct populated with the given parameters
// and sets the HTTP Verb to get
func NewGetActionMap(actionName string, actionMethod ActionMethod) *ActionMap {
	return &ActionMap{
		Verb:   "GET",
		Name:   actionName,
		Method: actionMethod,
	}
}

// NewPostActionMap returns a new ActionMap struct populated with the given parameters
// and sets the HTTP Verb to post
func NewPostActionMap(actionName string, actionMethod ActionMethod) *ActionMap {
	return &ActionMap{
		Verb:   "POST",
		Name:   actionName,
		Method: actionMethod,
	}
}

// NewPutActionMap returns a new ActionMap struct populated with the given parameters
// and sets the HTTP Verb to put
func NewPutActionMap(actionName string, actionMethod ActionMethod) *ActionMap {
	return &ActionMap{
		Verb:   "PUT",
		Name:   actionName,
		Method: actionMethod,
	}
}

// NewDeleteActionMap returns a new ActionMap struct populated with the given parameters
// and sets the HTTP Verb to delete
func NewDeleteActionMap(actionName string, actionMethod ActionMethod) *ActionMap {
	return &ActionMap{
		Verb:   "DELETE",
		Name:   actionName,
		Method: actionMethod,
	}
}
