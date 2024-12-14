package issues

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/andygrunwald/go-jira"
	"github.com/ebadfd/jira_sucks/lib"
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
	client := lib.JiraClient(s.CloudId, s.Token)

	projectKey := mux.Vars(r)["key"]
	issueKey := mux.Vars(r)["issueKey"]

	issue, _, err := client.Issue.Get(issueKey, &jira.GetQueryOptions{
		ProjectKeys: projectKey,
	})

	customFields, _, err := client.Issue.GetCustomFields(issueKey)

	if err != nil {
		panic(err)
	}

	d, _ := issue.Fields.MarshalJSON()
	fmt.Println(string(d))

	lib.Render(w, http.StatusOK, home.Issue(issue, customFields))

}

func (p *IssueServiceImpl) Issues(w http.ResponseWriter, r *http.Request) {
	s := context.Get(r, lib.AuthResults).(lib.AuthSession)
	client := lib.JiraClient(s.CloudId, s.Token)
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
		panic(err)
	}

	lib.Render(w, http.StatusOK, home.Issues(issues, res.StartAt, res.MaxResults, res.Total, jql))
}
