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
	HTTPPort    int
	HTTPSPort   int
	LogFilename string
	LogLevel    int
	TLSCertFile string
	TLSKeyFile  string
}

// NewConfigurationManager returns an empty new Configuration Manager struct
func NewConfigurationManager() *ConfigurationManager {
	return &ConfigurationManager{}
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
