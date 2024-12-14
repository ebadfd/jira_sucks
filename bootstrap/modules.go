package bootstrap

import (
	"github.com/ebadfd/jira_sucks/lib"
	"github.com/ebadfd/jira_sucks/pkg/oauth"
	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	lib.Module,
	oauth.Module,
)
