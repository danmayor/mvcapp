package mvcapp_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/digivance/mvcapp"
)

func TestNewConfigurationManagerFromFile(t *testing.T) {
	configFile := mvcapp.GetApplicationPath() + "/testconfig.json"
	configData := []byte("{\"HTTPPort\": 80,\"HTTPSPort\": 443,\"LogFilename\": \"./app.log\",\"LogLevel\": 4,\"TLSCertFile\": \"./mycert.crt\",\"TLSKeyFile\": \"./mycert.key\"}")
	err := ioutil.WriteFile(configFile, configData, 0644)
	defer os.RemoveAll(configFile)

	if err != nil {
		t.Errorf("Failed to create new configuration manager json file: %s", err)
	}

	config, err := mvcapp.NewConfigurationManagerFromFile(configFile)
	if err != nil {
		t.Errorf("Failed to create new configuration manager from file: %s", err)
	}

	if config.HTTPPort != 80 || config.HTTPSPort != 443 || config.LogFilename != "./app.log" ||
		config.LogLevel != 4 || config.TLSCertFile != "./mycert.crt" || config.TLSKeyFile != "./mycert.key" {
		t.Errorf("Failed to load the configuration manager object from file")
	}
}

func TestConfigurationManager_SaveFile(t *testing.T) {
	config := &mvcapp.ConfigurationManager{
		HTTPPort:    80,
		HTTPSPort:   443,
		LogFilename: "./app.log",
		LogLevel:    4,
		TLSCertFile: "./mycert.crt",
		TLSKeyFile:  "./mycert.key",
	}

	err := config.SaveFile("./testconfig.json")
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
