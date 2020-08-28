package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	logger "github.com/dm1trypon/easy-logger"
	"github.com/dm1trypon/rd-eye/models"
	"github.com/xeipuuv/gojsonschema"
)

// LC - logging category
const LC = "CONFIG"

// protocol - using for config's paths
const protocol = "file://"

// Cfg - contains config data
var Cfg = &models.Config{}

/*
Run - initializes the module for working with the service config.
Returns boolean result. True - config is OK, False - wrong config.
	- path <string> - path to config file json.

	- pathSchema <string> - path to config schema file json.
*/
func Run(path, pathSchema string) bool {
	logger.Info(LC, "Starting module Config")
	return make(path, pathSchema)
}

/*
make - checks the config file using the JSON scheme, then writes it to the Cfg structure.
Returns boolean result. True - config is OK, False - wrong config.

	- path <string> - path to config file json.

	- pathSchema <string> - path to config schema file json.
*/
func make(path, pathSchema string) bool {
	schemaLoader := gojsonschema.NewReferenceLoader(fmt.Sprint(protocol, pathSchema))
	documentLoader := gojsonschema.NewReferenceLoader(fmt.Sprint(protocol, path))

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		logger.Critical(LC, fmt.Sprint("Validation failed:", err.Error()))
		return false
	}

	if !result.Valid() {
		logger.Critical(LC, fmt.Sprint("Config is not valid"))
		for _, desc := range result.Errors() {
			logger.Error(LC, fmt.Sprint(desc))
		}

		return false
	}

	configData, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Critical(LC, fmt.Sprint("Error reading config file:", err.Error()))
		return false
	}

	if err := json.Unmarshal(configData, Cfg); err != nil {
		logger.Critical(LC, fmt.Sprint("Error parsing config file:", err.Error()))
		return false
	}

	return true
}
