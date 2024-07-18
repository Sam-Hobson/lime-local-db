package cli

import (
	"log/slog"
	"math"

	"github.com/go-errors/errors"
	conf "github.com/sam-hobson/internal/config"
	op "github.com/sam-hobson/internal/operations"
)

const (
	SetupIncompatibleFlags = math.MaxUint64 ^ SetupOff
	NewdbIncompatibleFlags = math.MaxUint64 ^ NewdbOff
	RmDbIncompatibleFlags  = math.MaxUint64 ^ RmdbOff
)

func ProcessArgs(flags *Flags) error {

	// Handle --config
	if flags.FlagsSet(SetupOff) {
		if flags.FlagsSet(SetupIncompatibleFlags) {
			slog.Error("--setup flag used when other flags are provided flags.", "log_code", "8e7e78d3", "flags", flags)
			return errors.Errorf("--setup flag used when other flags are provided: %s", flags)
		}

		err := op.Setup()
		return err
	}

	config, err := conf.GetConfig()

	if err != nil {
		slog.Warn("Cannot proceed with processing arguments, as no config is present.", "log_code", "36a6342f")
		return err
	}

	// Handle --new-db/-n
	if flags.FlagsSet(NewdbOff) {
		if flags.FlagsSet(NewdbIncompatibleFlags) {
			slog.Error("--new-db or -n flag used when incompatible flags are provided.", "log_code", "d2e62f94", "flags", flags)
			return errors.Errorf("--new-db or -n flag used when incompatible flags are provided: %s", flags)
		}

		err := op.NewDb(&flags.Newdb)
		return err
	}

	if flags.FlagsSet(RmdbOff) {
		if flags.FlagsSet(RmDbIncompatibleFlags) {
			slog.Error("--rm-db or -D flag used when incompatible flags are provided.", "log_code", "5e0a183f", "flags", flags)
			return errors.Errorf("--rm-db or -D flag used when incompatible flags are provided: %s", flags)
		}

		err := op.RmDb(flags.Rmdb, config)
		return err
	}

	return nil
}
