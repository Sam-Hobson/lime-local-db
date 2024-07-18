package cli

import (
	"log/slog"
	"math"

	"github.com/go-errors/errors"
	op "github.com/sam-hobson/internal/operations"
)

const (
	SetupIncompatibleFlags = math.MaxUint64 ^ SetupOff
    NewdbIncompatibleFlags = math.MaxUint64 ^ NewdbOff
)

func ProcessArgs(flags *Flags) error {

    // Handle --config
	if flags.FlagsSet(SetupOff) {
		if flags.FlagsSet(SetupIncompatibleFlags) {
			slog.Error("--setup flag used when other flags are provided flags.", "flags", flags)
			return errors.Errorf("--setup flag used when other flags are provided: %s", flags)
		}

		err := op.Setup()
		return err
	}

    // Handle --new-db/-n
	if flags.FlagsSet(NewdbOff) {
        if flags.FlagsSet(NewdbIncompatibleFlags) {
			slog.Error("--new-db or -n flag used when incompatible flags are provided.", "flags", flags)
			return errors.Errorf("--new-db or -n flag used when incompatible flags are provided: %s", flags)
        }
	}

	return nil
}
