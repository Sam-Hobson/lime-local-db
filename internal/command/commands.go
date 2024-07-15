package command

import (
	"log/slog"

	"github.com/go-errors/errors"
)

type argProcessor interface {
    error() error
    process(key string) argProcessor
}

type genericProcessor struct {
    err error
}

func (ps *genericProcessor) error() error {
    return ps.err
}

func (ps *genericProcessor) process(key string) argProcessor {
    switch key {
    case "new":
        return newNewDbCmd()
    }

    ps.err = errors.Errorf("Could not identify keyword.")
    return ps
}

func ProcessArgs(args []string) error {
    slog.Info("Arguments received.\n", "Args", args)

    var processor = genericProcessor{}

    for _, arg := range args {
        processor.process(arg)

        if processor.error() != nil {
            slog.Error("Could not process arguments.\n", "Args", args, "Failed_arg", arg)
            return processor.error()
        }
    }

    slog.Info("Successfully processed arguments.\n", "Args", args)
    return nil
}
