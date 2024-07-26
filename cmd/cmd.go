package cmd

import (
	"fmt"

	newdb "github.com/sam-hobson/cmd/new-db"
	rmdb "github.com/sam-hobson/cmd/rm-db"
	"github.com/spf13/cobra"
)

func NewCommand(version, commit string) *cobra.Command {
    cmd := &cobra.Command{
        Use: "limedb",
        Short: "Application for interacting with a simple database",
        Long:  "TODO: This",
        Version: buildVersion(version, commit),
    }

    cmd.AddCommand(newdb.NewCommand(), rmdb.NewCommand())

    return cmd
}

func buildVersion(version, commit string) string {
    return fmt.Sprintf("%s(%s)", version, commit)
}
