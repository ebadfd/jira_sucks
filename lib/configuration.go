package lib

import (
	"os"

	"github.com/mitchellh/mapstructure"
)

var globalConfiguration = Configuration{
	Port: "8080",
}

type Configuration struct {
	Port             string `mapstructure:"PORT"`
	JiraClientId     string `mapstructure:"JIRA_CLIENT_ID"`
	JiraClientSecret string `mapstructure:"JIRA_CLIENT_SECRET"`
	JiraAuthUri      string `mapstructure:"JIRA_AUTH_URI"`
	JiraTokenUri     string `mapstructure:"JIRA_TOKEN_URI"`
	SessionSecret    string `mapstructure:"SESSION_KEY"`
	Host             string `mapstructure:"APP_HOST"`
}

func (c Configuration) Validate() error {
	return nil
}

func NewConfiguration(logger Logger) *Configuration {
	vars := make(map[string]string)

	for _, key := range []string{
		"PORT",
		"JIRA_AUTH_URL",
		"JIRA_TOKEN_URI",
		"JIRA_CLIENT_ID",
		"JIRA_CLIENT_SECRET",
		"SESSION_KEY",
		"APP_HOST",
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
