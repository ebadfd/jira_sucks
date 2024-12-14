package jira_sucks

import (
	"fmt"
	"net/http"

	"github.com/ebadfd/jira_sucks/lib"
	"github.com/ebadfd/jira_sucks/pkg/handlers"
	"github.com/ebadfd/jira_sucks/pkg/oauth"
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
		oauth *oauth.JiraOAuthServiceImpl,
	) {
		serverHost := fmt.Sprintf(":%s", conf.Port)
		r := mux.NewRouter()
		r.Use(handlers.TrailingSlashMiddleware)

		r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/app", http.StatusPermanentRedirect)
		})

		auth := r.PathPrefix("/auth").Subrouter()
		auth.HandleFunc("/jira", oauth.OAuthJiraLogin)
		auth.HandleFunc("/callback", oauth.OAuthJiraCallback)

		app := r.PathPrefix("/app").Subrouter()
		app.Use(handlers.AuthMiddleware)

		app.HandleFunc("", handlers.Test)

		logger.Info(fmt.Sprintf("starting the web server on %s", serverHost))
		http.ListenAndServe(serverHost, r)
	}
}

func NewServeCommand() *ServerCommand {
	return &ServerCommand{}
}
