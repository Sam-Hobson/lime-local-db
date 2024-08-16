package trigger

import "github.com/spf13/cobra"

// TODO: Flesh out use/examples documentation.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "trigger [Subcommand]",
		Short:   "Operate on database triggers",
		Example: "limedb trigger",
	}

	cmd.AddCommand(
		newTriggerCommand(),
		templateTriggerCommand(),
		lsTriggerCommand(),
	)

	return cmd
}
