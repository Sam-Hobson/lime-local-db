package command

import (
	"log/slog"

	"github.com/go-errors/errors"
)

const newDbCmdUsage = "limedb new-db [DB NAME]"
const numArgs = 1

type newDbCmd struct {
	argNum int
	dbName string
	err    error
}

func (ps *newDbCmd) error() error {
	return ps.err
}

func (ps *newDbCmd) finished() {
    if ps.argNum != numArgs {
        slog.Error("The new-db keyword takes 1 argument.", "Usage", newDbCmdUsage)
        ps.err = errors.Errorf("The new-db keyword takes 1 argument.\n Usage: %s", newDbCmdUsage)
    }
}

func (ps *newDbCmd) process(key string) argProcessor {
	ps.argNum++

	if ps.argNum == 1 {
		ps.dbName = key
	} else {
        slog.Error("The new-db keyword only take 1 argument.", "Usage", newDbCmdUsage)
        ps.err = errors.Errorf("The new-db keyword only take 1 argument.\n Usage: %s", newDbCmdUsage)
	}

	return ps
}

func newNewDbCmd() *newDbCmd {
	return &newDbCmd{
		argNum: 0,
		dbName: "",
		err:    nil,
	}
}
