package cli

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/pflag"
)

const (
	setup          = "setup"
	newdb          = "new-db"
	newdbShorthand = "n"
	rmdb           = "rm-db"
	rmdbShorthand  = "D"

	SetupOff = uint8(0)
	NewdbOff = uint8(1)
	RmdbOff  = uint8(2)
)

type Flags struct {
	SetupDir bool
	Newdb    string
	Rmdb     string

	ProvidedFields uint64
}

func (f *Flags) FlagSet(flagOff uint8) bool {
	return (f.ProvidedFields & (1 << flagOff)) > 0
}

func (f *Flags) OtherFlagsSet(flagOff uint8) bool {
	return (f.ProvidedFields ^ (1 << flagOff)) > 0
}

func (f *Flags) OneOfSet(flagOff uint64) bool {
    return (f.ProvidedFields & flagOff) > 0
}

func (f *Flags) String() string {
	return fmt.Sprintf("%+v", *f)
}

func (f *Flags) setFlagIfProvided(name string, offset uint8) {
	if pflag.CommandLine.Changed(name) {
		f.ProvidedFields |= 1 << offset
	}
}

func GetFlags() *Flags {
	slog.Info("Parsing flags.")
	flags := &Flags{}

	pflag.BoolVar(&flags.SetupDir, setup, false, "Create a new default config in the $HOME directory.")
	pflag.StringVarP(&flags.Newdb, newdb, newdbShorthand, "", "Provide the name of a new database which will be created.")
	pflag.StringVarP(&flags.Rmdb, rmdb, rmdbShorthand, "", "Provide the name of a new database which will be deleted PERMANENTLY.")

	pflag.Parse()

	// If no flags were provided
	if pflag.NFlag() == 0 {
		pflag.Usage()
		os.Exit(2)
	}

	flags.setFlagIfProvided(setup, SetupOff)
	flags.setFlagIfProvided(newdb, NewdbOff)
	flags.setFlagIfProvided(rmdb, RmdbOff)

	slog.Info("Successfully parsed flags.", "flags", flags)

	return flags
}
