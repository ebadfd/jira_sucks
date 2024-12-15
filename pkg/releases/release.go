package releases

import (
	"net/http"
	"strconv"

	"github.com/ebadfd/jira_sucks/lib"
	"github.com/ebadfd/jira_sucks/views/home"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

type ReleaseServiceImpl struct {
	log  lib.Logger
	conf *lib.Configuration
}

func NewReleaseServiceImpl(log lib.Logger, conf *lib.Configuration) *ReleaseServiceImpl {
	return &ReleaseServiceImpl{
		log:  log,
		conf: conf,
	}
}

var Module = fx.Options(
	fx.Provide(NewReleaseServiceImpl),
)

func (p *ReleaseServiceImpl) ReleaseDetails(w http.ResponseWriter, r *http.Request) {
	s := context.Get(r, lib.AuthResults).(lib.AuthSession)
	client := lib.JiraClient(s.CloudId, s.Token)

	projectKey := mux.Vars(r)["key"]
	releaseId := mux.Vars(r)["releaseId"]

	i, err := strconv.Atoi(releaseId)

	if err != nil {
		panic(err)
	}

	release, _, err := client.Version.Get(i)

	if err != nil {
		panic(err)
	}

	lib.Render(w, http.StatusOK, home.Release(release, projectKey))
}
