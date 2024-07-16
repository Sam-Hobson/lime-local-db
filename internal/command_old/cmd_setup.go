package command

import (
	"github.com/go-errors/errors"
	"log/slog"

	. "github.com/sam-hobson/internal/command/executors"
)

const (
	setupCmdUsage   = "limedb setup [optional BASEDIR]"
	setupCmdNumArgs = 1
)

type setupCmd struct {
	argNum  int
	baseDir string
	err     error
}

func (ps *setupCmd) error() error {
	return ps.err
}

func (ps *setupCmd) onFinish() Executor {
	if ps.argNum == 0 {
		return NewSetupExecutor(ps.baseDir)
	}
	return nil
}

func (ps *setupCmd) process(key string) (Executor, argProcessor) {
	ps.argNum++

	if ps.argNum == 1 {
		ps.baseDir = key
	} else {
		slog.Error("The setup keyword takes 0 or 1 argument.", "Usage", setupCmdUsage)
		ps.err = errors.Errorf("The new-db keyword takes 0 or 1 argument.\n Usage: %s", setupCmdUsage)
	}

	return NewSetupExecutor(ps.baseDir), ps
}

func newSetupCmd() *setupCmd {
	return &setupCmd{
		argNum:  0,
		baseDir: "",
		err:     nil,
	}
}
