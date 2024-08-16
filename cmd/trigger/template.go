package trigger

import (
	"github.com/go-errors/errors"
	"github.com/sam-hobson/internal/database"
	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/cobra"
)

func templateTriggerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "template [Trigger name]",
		Short:   "Prints a template for a sqlite trigger.",
		Example: "limedb trigger template my_trigger",
		Args:    cobra.ExactArgs(1),

		RunE: runTemplateTriggerCommand,
	}

	return cmd
}

func runTemplateTriggerCommand(cmd *cobra.Command, args []string) error {
	databaseName := state.ApplicationState().GetSelectedDb()

	if databaseName == "" {
		util.Log("bd7490f4").Error("Cannot create trigger template if database is not selected.")
		return errors.Errorf("Cannot create trigger template if database is not selected")
	}

	template, err := database.TriggerTemplate(databaseName, args[0], "[TRIGGER TYPE]", "[BODY]")
	if err != nil {
		return err
	}

	cmd.Println(template)

	return nil
}
