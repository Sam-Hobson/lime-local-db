package trigger

import (
	"github.com/go-errors/errors"
	"github.com/sam-hobson/internal/database"
	"github.com/sam-hobson/internal/state"
	"github.com/sam-hobson/internal/util"
	"github.com/spf13/cobra"
)

func newTriggerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "new",
		Short:   "Create/add a new trigger.",
		Example: "limedb trigger new -f mytrigger.sqlite",
		Args:    cobra.MaximumNArgs(0),

		RunE: runNewTriggerCommand,
	}

	cmd.Flags().StringP("from-file", "f", "", "Add a trigger from the specified file")
	cmd.Flags().StringP("name", "n", "", "Specify a name for the provided trigger")
	cmd.Flags().StringP("message", "m", "", "Add a message/note associated with the trigger")

	return cmd
}

func runNewTriggerCommand(cmd *cobra.Command, args []string) error {
	fileName := util.PanicIfErr(cmd.Flags().GetString("from-file"))
    name := util.PanicIfErr(cmd.Flags().GetString("name"))
    message := util.PanicIfErr(cmd.Flags().GetString("message"))
	databaseName := state.ApplicationState().GetSelectedDb()

	if databaseName == "" {
		util.Log("5a4cfcc4").Error("Cannot add trigger if database is not selected.")
		return errors.Errorf("Cannot add trigger if database is not selected")
	}

	if fileName != "" {
        if name == "" {
            util.Log("01fb5ce7").Error("Cannot add a trigger from a file if a name is not provided for the trigger.", "File name", fileName)
            return errors.Errorf("Cannot add a trigger from a file if a name is not provided using --name, -n")
        }

		relFs := util.NewRelativeFsManager()
		if contents, err := relFs.ReadFileIntoMemry(fileName); err != nil {
			return err
		} else {
            err := database.CreateTriggerRaw(databaseName, name, contents, message)
            return err
		}
	}

	return nil
}
