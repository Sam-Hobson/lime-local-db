package cli

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/pflag"
)

const (
	setupDir       = "setup"
	newdb          = "new-db"
	newdbShorthand = "n"
	rmdb           = "rm-db"
	rmdbShorthand  = "D"
)

type Flags struct {
	setupDir         string
	setupDirProvided bool
	newdb            string
	newdbProvided    bool
	rmdb             string
	rmdbProvided     bool
}

func (f *Flags) String() string {
    return fmt.Sprintf("%+v", *f)
}

func GetFlags() *Flags {
    slog.Info("Parsing flags.")
	flags := &Flags{}

	pflag.StringVar(&flags.setupDir, setupDir, "", "Provide the location of a configuration file to be generated. $HOME will be used by default.")
	pflag.Lookup(setupDir).NoOptDefVal = os.Getenv("HOME")
	pflag.StringVarP(&flags.newdb, newdb, newdbShorthand, "", "Provide the name of a new database which will be created.")
	pflag.StringVarP(&flags.rmdb, rmdb, rmdbShorthand, "", "Provide the name of a new database which will be deleted PERMANENTLY.")

	pflag.Parse()

	flags.setupDirProvided = pflag.CommandLine.Changed(setupDir)
	flags.newdbProvided = pflag.CommandLine.Changed(newdb)
	flags.rmdbProvided = pflag.CommandLine.Changed(rmdb)

    slog.Info("Successfully parsed flags.", "flags", flags)

	return flags
}

func ProcessFlags() {
}

