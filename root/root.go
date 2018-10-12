package root

import "github.com/spf13/cobra"

// Register `cmd`.
func Register(cmd *cobra.Command) {
	Command.AddCommand(cmd)
}

// Command represents the base command when called without any subcommands
var Command = &cobra.Command{
	Use:               "casher",
	PersistentPreRunE: preRun,
}

func init() {
}

// PreRunNoop noop.
func PreRunNoop(c *cobra.Command, args []string) {
}

// preRun sets up global tasks used for most commands, some use PreRunNoop
// to remove this default behaviour.
func preRun(c *cobra.Command, args []string) error {
	return nil
}
