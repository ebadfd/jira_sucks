package lib

import (
	"os"

	"github.com/mitchellh/mapstructure"
)

var globalConfiguration = Configuration{
	Port: "3000",
}

type Configuration struct {
	Port string `mapstructure:"PORT"`
}

func (c Configuration) Validate() error {
	return nil
}

func NewConfiguration(logger Logger) *Configuration {
	vars := make(map[string]string)

	for _, key := range []string{
		"PORT",
	} {
		switch key {
		default:
			vars[key] = os.Getenv(key)
		}
	}

	err := mapstructure.Decode(vars, &globalConfiguration)

	if err != nil {
		logger.Fatal("enable to map env variables", err)
	}
	return &globalConfiguration
}
