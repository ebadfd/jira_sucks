package bootstrap

import (
	console "github.com/ebadfd/jira_sucks/cmd/jira_sucks"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:              "jira_sucks",
	Short:            "A lightweight Jira web client",
	Long:             "A minimalist and efficient web client for interacting with Jira. Designed for users who need a faster, simpler way to manage tasks and projects without the clutter of the official Jira interface.",
	TraverseChildren: true,
}

// App root of the application
type App struct {
	*cobra.Command
}

// NewApp creates new root command
func NewApp() App {
	cmd := App{
		Command: rootCmd,
	}
	cmd.AddCommand(console.GetSubCommands(CommonModules)...)
	return cmd
}

var RootApp = NewApp()
