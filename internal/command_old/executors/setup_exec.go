package executors

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/go-errors/errors"
	"github.com/sam-hobson/internal/config"
)


type SetupExecutor struct {
	baseDir string
}

func (e *SetupExecutor) Priority() int {
	return 0
}

func (e *SetupExecutor) Execute(state *ExecutionState) (*ExecutionState, error) {
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

	return state, err
}

func NewSetupExecutor(baseDir string) *SetupExecutor {
    return &SetupExecutor{
        baseDir: baseDir,
    }
}
