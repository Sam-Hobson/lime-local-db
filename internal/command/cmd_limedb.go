package command

import (
	"log/slog"

	"github.com/go-errors/errors"
)

const limedbCmdUsage = "limedb [OPERATION] [OPTIONS]"

type limedbCmd struct {
	err error
}

func (ps *limedbCmd) error() error {
	return ps.err
}

func (ps *limedbCmd) onFinish() {
	slog.Error("The limedb command requires an operation.", "Usage", limedbCmdUsage)
	ps.err = errors.Errorf("The limedb command requires an operation.\n Usage: %s", limedbCmdUsage)
}

func (ps *limedbCmd) process(key string) (Executor, argProcessor) {
	switch key {
	case "new-db":
		return nil, newNewDbCmd()
	}

	ps.err = errors.Errorf("Could not identify keyword.")
	return nil, ps
}

func newLimedbCmd() *limedbCmd {
	return &limedbCmd{
        err: nil,
    }
}