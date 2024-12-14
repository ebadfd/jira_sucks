package oauth

import (
	"go.uber.org/fx"
)

var JiraPermissions = []string{"read:me", "read:jira-work", "manage:jira-project", "read:jira-user", "write:jira-work"}

var Module = fx.Options(
	fx.Provide(NewJiraOAuthServiceImpl),
)
