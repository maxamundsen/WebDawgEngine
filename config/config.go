package config

import (
	"encoding/json"
	"gohttp/build"
	"io"
	"log"
	"os"
)

// the config package allows the use of runtime configuration options from a file.
// the config options live in a global, non-exported struct called 'config'.
// the InitConfiguration() function reads data from a config file, and
// populates this struct. This function should be called once when the program
// starts.
//
// typically the configuration is read from `config.json`, however
// when the `devel` build tag is set, the options are read from `config.devel.json`
//
// since the configuration is initialized at runtime, embedding the configuration
// options would not make sense. the config file *must* be present on the host filesystem.
//
// the config struct is not exported, and must be retrieved via the GetConfiguration() function.
// because of this, you cannot modify the values inside configuration from anywhere in the program
// (except the `config` package directly)

// define all configuration options here
type configuration struct {
	Host             string `json:"Host"`
	CookieExpiryDays int    `json:"CookieExpiryDays"`
}

var config configuration

// specify default values if none are specified in the config file
func setDefaultValues() {
	if config.CookieExpiryDays == 0 {
		config.CookieExpiryDays = 7
	}

	if config.Host == "" {
		config.Host = "localhost:8080"
	}
}

func ReadConfiguration() {
	var configFile *os.File
	var err error

	if build.DEVEL {
		// if a config.devel.json exists
		configFile, err = os.Open("config.devel.json")

		// fallback to default config if devel config does not exist
		if err != nil {
			configFile, err = os.Open("config.json")
		}
	} else {
		configFile, err = os.Open("config.json")
	}

	if err != nil {
		log.Println(err)
	}

	configBytes, readErr := io.ReadAll(configFile)

	if readErr != nil {
		log.Fatal("Cannot read the configuration file.")
	}

	defer configFile.Close()

	json.Unmarshal(configBytes, &config)

	setDefaultValues()
	log.Println("Loaded configuration file (" + configFile.Name() + ")")
}

func GetConfiguration() configuration {
	return config
}
