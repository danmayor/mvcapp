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
	"strings"
	"time"
)

// Application is our global scope object (E.g. application wide configuration
// and manager systems)
type Application struct {
	// RouteManager is our Route Manager system
	RouteManager *RouteManager

	Config *ConfigurationManager

	// HTTPServer is the http.Server object being used to host http transport
	HTTPServer *http.Server

	// HTTPSServer is the http.Server object being used to host https transport
	HTTPSServer *http.Server
}

// NewApplication returns a new default MVC Application object
func NewApplication() *Application {
	rtn := &Application{
		RouteManager: NewRouteManager(),
		Config:       NewConfigurationManager(),
		HTTPServer:   nil,
		HTTPSServer:  nil,
	}

	if LogFilename == "" {
		SetLogFilename("./mvcapp.log")
	}

	LogTrace("Application initialized")
	return rtn
}

// NewApplicationFromConfig returns a new MVC Application object with the provided
// configuration manager object
func NewApplicationFromConfig(config *ConfigurationManager) *Application {
	rtn := &Application{
		RouteManager: NewRouteManager(),
		Config:       config,
		HTTPServer:   nil,
		HTTPSServer:  nil,
	}

	if LogFilename == "" {
		SetLogFilename("./mvcapp.log")
	}

	rtn.RouteManager.SessionManager.SessionTimeout = time.Duration(config.HTTPSessionTimeout) * time.Minute

	LogTrace("Application initialized")
	return rtn
}

// NewApplicationFromConfigFile returns a new MVC Application object  constructed
// from the provided json config file
func NewApplicationFromConfigFile(filename string) (*Application, error) {
	config, err := NewConfigurationManagerFromFile(filename)
	if err != nil {
		return nil, err
	}

	return NewApplicationFromConfig(config), nil
}

// Stop is used to stop hosting this MVC Application. You can call one of the Run methods to restart
func (app *Application) Stop() error {
	if app.HTTPServer != nil {
		app.HTTPServer.Shutdown(nil)
	}

	if app.HTTPSServer != nil {
		if err := app.HTTPSServer.Shutdown(nil); err != nil {
			return err
		}
	}

	return nil
}

// Run is used to execute this MVC Application (direct http socket server)
func (app *Application) Run() error {
	if app.HTTPServer != nil {
		return errors.New("Can not run application, HTTPServer already in use")
	}

	config := app.Config
	addr := fmt.Sprintf("%s:%d", config.BindAddress, config.HTTPPort)
	app.HTTPServer = &http.Server{Addr: addr}
	app.HTTPServer.Handler = http.HandlerFunc(app.RouteManager.HandleRequest)

	return app.HTTPServer.ListenAndServe()
}

// RunSecure is used to execute this MVC Application over HTTPS/TLS (direct https socket server)
func (app *Application) RunSecure(certFile string, keyFile string) error {
	if app.HTTPSServer != nil {
		return errors.New("Can not RunSecure, HTTPSServer already in use")
	}

	config := app.Config
	addr := fmt.Sprintf("%s:%d", config.BindAddress, config.HTTPSPort)
	app.HTTPSServer = &http.Server{Addr: addr}
	app.HTTPSServer.Handler = http.HandlerFunc(app.RouteManager.HandleRequest)

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

	var err error

	config := app.Config
	addr := fmt.Sprintf("%s:%d", config.BindAddress, config.HTTPPort)
	app.HTTPServer = &http.Server{Addr: addr}
	app.HTTPServer.Handler = http.HandlerFunc(app.RedirectSecure)
	go func() {
		err = app.HTTPServer.ListenAndServe()
	}()

	addr = fmt.Sprintf("%s:%d", config.BindAddress, config.HTTPSPort)
	app.HTTPSServer = &http.Server{Addr: addr}
	app.HTTPSServer.Handler = http.HandlerFunc(app.RouteManager.HandleRequest)
	go func() {
		if certFile == "" {
			certFile = config.TLSCertFile
		}

		if keyFile == "" {
			keyFile = config.TLSKeyFile
		}

		err = app.HTTPSServer.ListenAndServeTLS(certFile, keyFile)
	}()

	// Here is an internal management thread if needed, it waits for app.Config.TaskDuration per tick
	for err == nil {
		time.Sleep(time.Duration(config.TaskDuration) * time.Second)
	}

	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("Failed to launch application in forced secure mode: %s", err)
		}

		return err
	}

	return err
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

	var err error

	config := app.Config
	addr := fmt.Sprintf("%s:%d", config.DomainName, config.HTTPPort)
	app.HTTPServer = &http.Server{Addr: addr}
	app.HTTPServer.Handler = http.HandlerFunc(app.RedirectSecureJS)
	go func() {
		err = app.HTTPServer.ListenAndServe()
	}()

	addr = fmt.Sprintf("%s:%d", config.BindAddress, config.HTTPSPort)
	app.HTTPSServer = &http.Server{Addr: addr}
	app.HTTPSServer.Handler = http.HandlerFunc(app.RouteManager.HandleRequest)
	go func() {
		err = app.HTTPSServer.ListenAndServeTLS(certFile, keyFile)
	}()

	for err == nil {
		time.Sleep(time.Duration(config.TaskDuration) * time.Second)
	}

	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("Failed to launch application in forced secure via javascript mode: %s", err)
		}

		return err
	}

	return err
}

// RedirectSecure is used to submit an http redirect from http to https when forcing secure
func (app *Application) RedirectSecure(w http.ResponseWriter, req *http.Request) {
	// To allow for google site ownership verification
	if app.Config.AllowGoogleAuthFiles && strings.HasPrefix(req.URL.Path, "google") && strings.HasSuffix(req.URL.Path, ".html") {
		path := fmt.Sprintf("%s/%s", GetApplicationPath(), strings.TrimLeft(req.URL.Path, "/"))
		http.ServeFile(w, req, path)
		return
	}

	target := fmt.Sprintf("https://%s%s", req.Host, req.URL.Path)

	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}

	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

// RedirectSecureJS is used to submit a javascript redirect from http to https when forcing secure
func (app *Application) RedirectSecureJS(w http.ResponseWriter, req *http.Request) {
	// To allow for google site ownership verification
	if app.Config.AllowGoogleAuthFiles && strings.HasPrefix(req.URL.Path, "/google") && strings.HasSuffix(req.URL.Path, ".html") {
		path := fmt.Sprintf("%s/%s", GetApplicationPath(), strings.TrimLeft(req.URL.Path, "/"))
		http.ServeFile(w, req, path)
		return
	}

	data := fmt.Sprintf("<html><head><title>Redirecting to secure site mode</title></head><body><script type=\"text/javascript\">window.location.href='//%s';</script></body>", req.URL.Path)
	w.Write([]byte(data))
}
