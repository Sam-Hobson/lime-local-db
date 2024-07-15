package command

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/go-errors/errors"
	"github.com/sam-hobson/internal/config"
)

const setupCmdUsage = "limedb setup [optional BASEDIR]"
const setupCmdNumArgs = 1

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
		return &setupExecutor{
			baseDir: ps.baseDir,
		}
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

	return &setupExecutor{
		baseDir: ps.baseDir,
	}, ps
}

func newSetupCmd() *setupCmd {
	return &setupCmd{
		argNum:  0,
		baseDir: "",
		err:     nil,
	}
}

type setupExecutor struct {
	baseDir string
}

func (e *setupExecutor) Priority() int {
	return 0
}

func (e *setupExecutor) Execute(state *ExecutionState) (*ExecutionState, error) {
	if state.Config != nil {
		slog.Error("limedb setup called when config already exists.")
		return state, errors.Errorf("limedb setup called when config already exists.")
	}

	home := os.Getenv("HOME")

	var err error

	if e.baseDir == "" {
		err = config.CreateDefaultConfig(home, filepath.Join(home, config.LimeDir))
	} else {
		err = config.CreateDefaultConfig(home, e.baseDir)
	}

	if err != nil {
		return state, err
	}

	state.Config, err = config.ParseConfig(home)

	return state, err
}
