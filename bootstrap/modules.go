package bootstrap

import (
	"github.com/ebadfd/jira_sucks/lib"
	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	lib.Module,
)
