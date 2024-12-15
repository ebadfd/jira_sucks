package projects

import (
	"net/http"

	"github.com/ebadfd/jira_sucks/lib"
	"github.com/ebadfd/jira_sucks/views"
	"github.com/ebadfd/jira_sucks/views/home"
	"github.com/gorilla/context"
	"go.uber.org/fx"
)

type ProjectServiceImpl struct {
	log  lib.Logger
	conf *lib.Configuration
}

func NewProjectServiceImpl(log lib.Logger, conf *lib.Configuration) *ProjectServiceImpl {
	return &ProjectServiceImpl{
		log:  log,
		conf: conf,
	}
}

var Module = fx.Options(
	fx.Provide(NewProjectServiceImpl),
)

func (p *ProjectServiceImpl) Projects(w http.ResponseWriter, r *http.Request) {
	s := context.Get(r, lib.AuthResults).(lib.AuthSession)
	client, err := lib.JiraClient(s.CloudId, s.Token)

	if err != nil {
		lib.Render(w, http.StatusBadRequest, views.ErrorPage(err))
		return
	}

	projects, _, err := client.Project.GetList()

	if err != nil {
		lib.Render(w, http.StatusBadRequest, views.ErrorPage(err))
		return
	}

	lib.Render(w, http.StatusOK, home.Index(projects))
}
