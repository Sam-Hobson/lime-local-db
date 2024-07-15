package command

import (
	"log/slog"
	"slices"

	cp "github.com/bigkevmcd/go-configparser"
)

type argProcessor interface {
	error() error
	process(key string) (Executor, argProcessor)
	onFinish()
}

type Executor interface {
	Priority() int
	Execute(state *ExecutionState) (*ExecutionState, error)
}

type ExecutionState struct {
	Config *cp.ConfigParser
}

func ProcessArgs(args []string) ([]Executor, error) {
	slog.Info("Arguments received.\n", "Args", args)

	var processor argProcessor = newLimedbCmd()

	executors := make([]Executor, 0)

	for _, arg := range args {
		executor, new_processor := processor.process(arg)

		if processor.error() != nil {
			slog.Error("Could not process arguments.\n", "Args", args, "Failed_arg", arg)
			return nil, processor.error()
		}

		if executor != nil {
			executors = append(executors, executor)
		}

		processor = new_processor
	}

	processor.onFinish()

	if processor.error() != nil {
		slog.Error("Something went wrong processing args.\n", "Args", args)
		return nil, processor.error()
	}

	// Sort by priority
	slices.SortFunc[[]Executor, Executor](executors, func(a Executor, b Executor) int {
		return a.Priority() - b.Priority()
	})

	slog.Info("Successfully processed arguments.\n", "Args", args)
	return executors, nil
}
