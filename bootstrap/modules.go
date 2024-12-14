package bootstrap

import (
	"github.com/ebadfd/jira_sucks/lib"
	"github.com/ebadfd/jira_sucks/pkg/issues"
	"github.com/ebadfd/jira_sucks/pkg/oauth"
	"github.com/ebadfd/jira_sucks/pkg/projects"
	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	lib.Module,
	oauth.Module,
	projects.Module,
	issues.Module,
)
