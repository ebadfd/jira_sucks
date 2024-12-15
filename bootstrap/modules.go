package bootstrap

import (
	"github.com/ebadfd/jira_sucks/lib"
	"github.com/ebadfd/jira_sucks/pkg/issues"
	"github.com/ebadfd/jira_sucks/pkg/oauth"
	"github.com/ebadfd/jira_sucks/pkg/projects"
	"github.com/ebadfd/jira_sucks/pkg/releases"
	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	lib.Module,
	oauth.Module,
	projects.Module,
	issues.Module,
	releases.Module,
)
