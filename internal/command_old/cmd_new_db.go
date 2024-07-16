package command

import (
	"log/slog"

	"github.com/go-errors/errors"
	. "github.com/sam-hobson/internal/command/executors"
)

const newDbCmdUsage = "limedb new-db [DB NAME]"
const newDbCmdnumArgs = 1

type newDbCmd struct {
	argNum int
	dbName string
	err    error
}

func (ps *newDbCmd) error() error {
	return ps.err
}

func (ps *newDbCmd) onFinish() Executor {
	if ps.argNum != newDbCmdnumArgs {
		slog.Error("The new-db keyword takes 1 argument.", "Usage", newDbCmdUsage)
		ps.err = errors.Errorf("The new-db keyword takes 1 argument.\n Usage: %s", newDbCmdUsage)
	}

	return nil
}

func (ps *newDbCmd) process(key string) (Executor, argProcessor) {
	ps.argNum++

	if ps.argNum == 1 {
		ps.dbName = key
	} else {
		slog.Error("The new-db keyword only takes 1 argument.", "Usage", newDbCmdUsage)
		ps.err = errors.Errorf("The new-db keyword only takes 1 argument.\n Usage: %s", newDbCmdUsage)
	}

	return nil, ps
}

func newNewDbCmd() *newDbCmd {
	return &newDbCmd{
		argNum: 0,
		dbName: "",
		err:    nil,
	}
}
