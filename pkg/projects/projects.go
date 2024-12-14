package projects

import (
	"net/http"

	"github.com/ebadfd/jira_sucks/lib"
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
	client := lib.JiraClient(s.CloudId, s.Token)

	projects, _, err := client.Project.GetList()

	if err != nil {
		panic(err)
	}

	lib.Render(w, http.StatusOK, home.Index(projects))
}

