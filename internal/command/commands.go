package command

import (
	"log/slog"
)

type argProcessor interface {
	error() error
	process(key string) argProcessor
	finished()
}

func ProcessArgs(args []string) error {
	slog.Info("Arguments received.\n", "Args", args)

	var processor argProcessor = newLimedbCmd()

	for _, arg := range args {
		new_processor := processor.process(arg)

		if processor.error() != nil {
			slog.Error("Could not process arguments.\n", "Args", args, "Failed_arg", arg)
			return processor.error()
		}

		processor = new_processor
	}

	processor.finished()

	if processor.error() != nil {
		slog.Error("Something went wrong processing args.\n", "Args", args)
		return processor.error()
	}

	slog.Info("Successfully processed arguments.\n", "Args", args)
	return nil
}
