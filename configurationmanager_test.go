/*
	Digivance MVC Application Framework
	Configuration Manage Unit Tests
	Dan Mayor (dmayor@digivance.com)

	This file defines the version 0.3.0 compatibility of configurationmanager.go functions. These functions are written
	to demonstrate and test the intended use cases of the functions in configurationmanager.go
*/

package mvcapp_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/digivance/mvcapp"
)

// TestNewConfigurationManagerFromFile ensures that mvcapp.NewConfigurationManagerFromFile returns the expected value
func TestNewConfigurationManagerFromFile(t *testing.T) {
	configFile := mvcapp.GetApplicationPath() + "/testconfig.json"
	errorConfigFile := mvcapp.GetApplicationPath() + "/failconfig.json"

	configData := []byte("{\"HTTPPort\": 80,\"HTTPSPort\": 443,\"LogFilename\": \"./app.log\",\"LogLevel\": 4,\"TLSCertFile\": \"./mycert.crt\",\"TLSKeyFile\": \"./mycert.key\"}")
	errorConfigData := []byte("this is not json...")

	config, err := mvcapp.NewConfigurationManagerFromFile(configFile)
	if err == nil {
		t.Error("Failed to fail on missing file")
	}

	err = ioutil.WriteFile(configFile, configData, 0644)
	defer os.RemoveAll(configFile)

	if err != nil {
		t.Errorf("Failed to create new configuration manager json file: %s", err)
	}

	err = ioutil.WriteFile(errorConfigFile, errorConfigData, 0644)
	defer os.RemoveAll(errorConfigFile)

	if err != nil {
		t.Errorf("Failed to create new configuration manager json file: %s", err)
	}

	config, err = mvcapp.NewConfigurationManagerFromFile(errorConfigFile)
	if err == nil {
		t.Error("Failed to fail creating new configuration file from invalid json file")
	}

	config, err = mvcapp.NewConfigurationManagerFromFile(configFile)
	if err != nil {
		t.Errorf("Failed load config file: %s", err)
	}

	if config.HTTPPort != 80 || config.HTTPSPort != 443 || config.LogFilename != "./app.log" ||
		config.LogLevel != 4 || config.TLSCertFile != "./mycert.crt" || config.TLSKeyFile != "./mycert.key" {
		t.Errorf("Failed to load the configuration manager object from file")
	}
}

// TestConfigurationManager_SaveFile ensures that ConfigurationManager.SaveFile operates as expected
func TestConfigurationManager_SaveFile(t *testing.T) {
	config := &mvcapp.ConfigurationManager{
		HTTPPort:    80,
		HTTPSPort:   443,
		LogFilename: "./app.log",
		LogLevel:    4,
		TLSCertFile: "./mycert.crt",
		TLSKeyFile:  "./mycert.key",
	}

	err := config.SaveFile("?:\\*$%«╝╗")
	if err == nil {
		t.Error("Failed to fail to save to invalid filename")
	}

	err = config.SaveFile("./testconfig.json")
	if err != nil {
		t.Errorf("Failed to save configuration file: %s", err)
	}

	defer os.RemoveAll(mvcapp.GetApplicationPath() + "/testconfig.json")
	testConfig, err := mvcapp.NewConfigurationManagerFromFile("./testconfig.json")
	if err != nil {
		t.Errorf("Failed to load saved configuration file: %s", err)
	}

	if testConfig.HTTPPort != config.HTTPPort {
		t.Error("Failed to properly load configuration file, values don't match :(")
	}
}
