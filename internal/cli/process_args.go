package cli

import (
	"log/slog"
	"math"

	"github.com/go-errors/errors"
	op "github.com/sam-hobson/internal/operations"
)

const SetupIncompatibleWithFlags = math.MaxUint64 ^ (1 << SetupOff)

func ProcessArgs(flags *Flags) error {

	if flags.FlagSet(SetupOff) {
		if flags.OneOfSet(SetupIncompatibleWithFlags) {
			slog.Error("--setup flag used when other flags are provided flags.", "flags", flags)
			return errors.Errorf("--setup flag used when other flags are provided: %s", flags)
		}

		err := op.Setup()
		return err
	}

	if flags.FlagSet(NewdbOff) {
	}

	return nil
}
