package jira_sucks

import (
	"fmt"
	"net/http"

	"github.com/ebadfd/jira_sucks/lib"
	"github.com/ebadfd/jira_sucks/pkg/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

type ServerCommand struct{}

func (s *ServerCommand) Short() string {
	return "Run the Jira Client application"
}

func (s *ServerCommand) Setup(cmd *cobra.Command) {}

func (s *ServerCommand) Run() lib.CommandRunner {
	return func(
		conf *lib.Configuration,
		logger lib.Logger,
	) {
		serverHost := fmt.Sprintf(":%s", conf.Port)
		r := mux.NewRouter()
		r.HandleFunc("/", handlers.Test)

		logger.Info(fmt.Sprintf("starting the web server on %s", serverHost))
		http.ListenAndServe(serverHost, r)
	}
}

func NewServeCommand() *ServerCommand {
	return &ServerCommand{}
}
