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

	SetupOff = uint64(1 << 0)
	NewdbOff = uint64(1 << 1)
	RmdbOff  = uint64(1 << 2)
)

type Flags struct {
	SetupDir bool
	Newdb    string
	Rmdb     string

	ProvidedFields uint64
}

func (f *Flags) FlagsSet(flagOff uint64) bool {
	return (f.ProvidedFields & flagOff) > 0
}

func (f *Flags) OtherFlagsSet(flagOff uint64) bool {
	return (f.ProvidedFields ^ flagOff) > 0
}

func (f *Flags) String() string {
	return fmt.Sprintf("%+v", *f)
}

func (f *Flags) setFlagIfProvided(name string, offset uint64) {
	if pflag.CommandLine.Changed(name) {
		f.ProvidedFields |= offset
	}
}

func GetFlags() *Flags {
	slog.Info("Parsing flags.", "Hash", "ae0a99b8")
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

	slog.Info("Successfully parsed flags.", "Hash", "fea120e1", "flags", flags)

	return flags
}
