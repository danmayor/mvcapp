/*
	Digivance MVC Application Framework
	Configuration Manager Object
	Dan Mayor (dmayor@digivance.com)

	This file defines the general application settings and configurations system. This system
	allows the package caller to easily interact with all the configurable settings as well as
	read / write them to human readable json configuration files
*/

package mvcapp

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

// ConfigurationManager is a simple json serializer wrapper that allows an
// mvcapp to easily save and load human friendly json files for application
// configuration values
type ConfigurationManager struct {
	// AppName is the display name for this application wherever applicable
	AppName string

	// AppVersion is the display version of this application wherever applicable
	AppVersion string

	// DomainName is the FQDN for this application
	DomainName string

	// BindAddress is the TCP/IP address to bind the listener to (normally leave blank for all)
	BindAddress string

	// HTTPPort is the TCP/IP Port number to listen on for plain http
	HTTPPort int

	// HTTPSPort is the TCP/IP Port number to listen on for TLS http
	HTTPSPort int

	// LogFilename is the full path and filename to write logging messages to (based on log level)
	LogFilename string

	// LogLevel represents the level of messages that are written to the log file. Allows for caller
	// to switch at runtime based on conditions or needs
	LogLevel int

	// TLSCertFile is the full path and filename of the TLS Certificate file to use for HTTPS
	TLSCertFile string

	// TLSKeyFile is the full path and filename of the TLS Key file to use for HTTPS
	TLSKeyFile string

	// AllowGoogleAuthFiles will allow the app to serve google site authentication files over plain
	// http even if the app is forcing all traffic to https (normally irrelevent)
	AllowGoogleAuthFiles bool

	// HTTPSessionIDKey is the name of the cookie that will store our user's session id if using
	// the built in http session manager system
	HTTPSessionIDKey string

	// HTTPSessionTimeout is the number of minutes that a users http session data will be stored in
	// memory between requests
	HTTPSessionTimeout int64

	// TaskDuration is the number of seconds to idle between evaluating internal tasks (such as cleaning
	// user http sessions in memory)
	TaskDuration int64

	// DefaultController is used to define where requests to the root of this DomainName are routed
	// Should be Home in most cases
	DefaultController string

	// DefaultAction is used to define the default action method to be executed if unable to map the
	// the requested action. Should be Index in most cases
	DefaultAction string
}

// NewConfigurationManager returns an empty new Configuration Manager struct
func NewConfigurationManager() *ConfigurationManager {
	return &ConfigurationManager{
		AppName:    "MyApp",
		AppVersion: "0.0.0",

		DomainName:  "domain.tld",
		BindAddress: "",
		HTTPPort:    80,
		HTTPSPort:   443,

		LogFilename: "./app.log",
		LogLevel:    2,

		TLSCertFile: "",
		TLSKeyFile:  "",

		AllowGoogleAuthFiles: true,

		HTTPSessionIDKey:   "mvcapp.sessionid",
		HTTPSessionTimeout: 30,
		TaskDuration:       60,

		DefaultController: "Home",
		DefaultAction:     "Index",
	}
}

// NewConfigurationManagerFromFile returns a new Configuration Manager struct
// constructed from the values in the provided json file
func NewConfigurationManagerFromFile(filename string) (*ConfigurationManager, error) {
	if strings.HasPrefix(filename, "~/") || strings.HasPrefix(filename, "./") {
		filename = GetApplicationPath() + filename[1:]
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := NewConfigurationManager()
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// SaveFile will save the current values to a human readable json configuration file
func (config *ConfigurationManager) SaveFile(filename string) error {
	if strings.HasPrefix(filename, "~/") || strings.HasPrefix(filename, "./") {
		filename = GetApplicationPath() + filename[1:]
	}

	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
