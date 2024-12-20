package issues

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/andygrunwald/go-jira"
	"github.com/ebadfd/jira_sucks/lib"
	"github.com/ebadfd/jira_sucks/views"
	"github.com/ebadfd/jira_sucks/views/home"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

type IssueServiceImpl struct {
	log  lib.Logger
	conf *lib.Configuration
}

func NewIssueServiceImpl(log lib.Logger, conf *lib.Configuration) *IssueServiceImpl {
	return &IssueServiceImpl{
		log:  log,
		conf: conf,
	}
}

var Module = fx.Options(
	fx.Provide(NewIssueServiceImpl),
)

func (p *IssueServiceImpl) IssueDetails(w http.ResponseWriter, r *http.Request) {
	s := context.Get(r, lib.AuthResults).(lib.AuthSession)
	client, err := lib.JiraClient(s.CloudId, s.Token)

	if err != nil {
		lib.Render(w, http.StatusBadRequest, views.ErrorPage(err))
		return
	}

	projectKey := mux.Vars(r)["key"]
	issueKey := mux.Vars(r)["issueKey"]

	issue, _, err := client.Issue.Get(issueKey, &jira.GetQueryOptions{
		ProjectKeys: projectKey,
	})

	if err != nil {
		lib.Render(w, http.StatusBadRequest, views.ErrorPage(err))
		return
	}

	customFields, _, err := client.Issue.GetCustomFields(issueKey)

	if err != nil {
		lib.Render(w, http.StatusBadRequest, views.ErrorPage(err))
		return
	}

	transitions, _, err := client.Issue.GetTransitions(issueKey)

	if err != nil {
		lib.Render(w, http.StatusBadRequest, views.ErrorPage(err))
		return
	}

	lib.Render(w, http.StatusOK, home.Issue(issue, customFields, transitions))
}

func (p *IssueServiceImpl) IssueDetailsTransitions(w http.ResponseWriter, r *http.Request) {
	s := context.Get(r, lib.AuthResults).(lib.AuthSession)
	client, err := lib.JiraClient(s.CloudId, s.Token)

	if err != nil {
		lib.Render(w, http.StatusBadRequest, views.ErrorPage(err))
		return
	}

	projectKey := mux.Vars(r)["key"]
	issueKey := mux.Vars(r)["issueKey"]
	transitionId := r.FormValue("transitionId")

	_, err = client.Issue.DoTransition(issueKey, transitionId)

	if err != nil {
		lib.Render(w, http.StatusBadRequest, views.ErrorPage(err))
		return
	}

	issue, _, err := client.Issue.Get(issueKey, &jira.GetQueryOptions{
		ProjectKeys: projectKey,
	})

	if err != nil {
		lib.Render(w, http.StatusBadRequest, views.ErrorPage(err))
		return
	}

	transitions, _, err := client.Issue.GetTransitions(issueKey)

	if err != nil {
		lib.Render(w, http.StatusBadRequest, views.ErrorPage(err))
		return
	}

	lib.Render(w, http.StatusOK, home.Transitions(issue, transitions))

}

func (p *IssueServiceImpl) Issues(w http.ResponseWriter, r *http.Request) {
	s := context.Get(r, lib.AuthResults).(lib.AuthSession)
	client, err := lib.JiraClient(s.CloudId, s.Token)

	if err != nil {
		lib.Render(w, http.StatusBadRequest, views.ErrorPage(err))
		return
	}

	projectKey := mux.Vars(r)["key"]

	startAtStr := r.URL.Query().Get("startAt")
	maxResultsStr := r.URL.Query().Get("maxResults")
	jql := r.URL.Query().Get("jql")

	if jql == "" {
		jql = fmt.Sprintf("project = \"%s\" ORDER BY created DESC", projectKey)
	}

	startAt, err := strconv.Atoi(startAtStr)
	if err != nil {
		startAt = 0
	}

	maxResults, err := strconv.Atoi(maxResultsStr)
	if err != nil {
		maxResults = 20
	}

	issues, res, err := client.Issue.Search(
		jql,
		&jira.SearchOptions{
			MaxResults: maxResults,
			StartAt:    startAt,
		},
	)

	if err != nil {
		lib.Render(w, http.StatusBadRequest, views.ErrorPage(err))
		return
	}

	if lib.IsHTMXRequest(r) {
		lib.Render(w, http.StatusOK, home.IssuesList(issues, res.StartAt, res.MaxResults, res.Total, jql))
		return
	} else {
		lib.Render(w, http.StatusOK, home.Issues(issues, res.StartAt, res.MaxResults, res.Total, jql))
		return
	}
}
