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
	cmd.Flags().StringP("from-directory", "d", "", "Add all the triggers within files in a given directory")
	cmd.Flags().StringP("name", "n", "", "Specify a name for the provided trigger")
	cmd.Flags().StringP("message", "m", "", "Add a message/note associated with the trigger")

	return cmd
}

func runNewTriggerCommand(cmd *cobra.Command, args []string) error {
	databaseName := state.ApplicationState().GetSelectedDb()
	if databaseName == "" {
		util.Log("5a4cfcc4").Error("Cannot add trigger if database is not selected.")
		return errors.Errorf("Cannot add trigger if database is not selected")
	}

	fileName := util.PanicIfErr(cmd.Flags().GetString("from-file"))
	directory := util.PanicIfErr(cmd.Flags().GetString("from-directory"))

	if fileName != "" && directory != "" {
		util.Log("56d9d120").Error("Triggers cannot be added from a file and directory at the same time.", "File name", fileName, "Directory", directory)
		return errors.Errorf("Triggers cannot be added from a file and directory at the same time")
	}

	if fileName != "" {
		name := util.PanicIfErr(cmd.Flags().GetString("name"))
		if name == "" {
			util.Log("01fb5ce7").Error("Cannot add a trigger from a file if a name is not provided for the trigger.", "File name", fileName)
			return errors.Errorf("Cannot add a trigger from a file if a name is not provided using --name, -n")
		}

		relFs := util.NewRelativeFsManager()
		if contents, err := relFs.ReadFileIntoMemry(fileName); err != nil {
			return err
		} else {
			message := util.PanicIfErr(cmd.Flags().GetString("message"))
			err := database.CreateTriggerRaw(databaseName, name, contents, message)
			return err
		}
	} else if directory != "" {
	}

	return nil
}
