package cli

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"github.com/go-errors/errors"
	op "github.com/sam-hobson/internal/operations"
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

func parseNewDbData(rawInput string) (*op.NewdbData, error) {
	re := regexp.MustCompile(`^(.*?)\s+\[(.*?)\]$`)
	matches := re.FindStringSubmatch(rawInput)

	if len(matches) != 3 {
		slog.Error("--new-db has malformed input. Example: limedb --new-db \"my-db [attr1 attr2 attr3 attr4]\"", "log_code", "20886ce7", "Raw_input", rawInput)
		return nil, errors.Errorf("--new-db has malformed input. Example: limedb --new-db \"my-db [attr1 attr2 attr3 attr4]\"")
	}

	name := matches[1]
	attrs := strings.Fields(matches[2])

	return &op.NewdbData{
		Dbname:   name,
		ColNames: attrs,
	}, nil
}

type Flags struct {
	SetupDir bool
	Newdb    op.NewdbData
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
	slog.Info("Parsing flags.", "log_code", "ae0a99b8")
	flags := &Flags{}

	pflag.BoolVar(&flags.SetupDir, setup, false, "Create a new default config in the $HOME directory.")
	newdbFlag := pflag.StringP(newdb, newdbShorthand, "", "Provide the name of a new database which will be created.")
	pflag.StringVarP(&flags.Rmdb, rmdb, rmdbShorthand, "", "Provide the name of a new database which will be deleted PERMANENTLY.")

	pflag.Parse()

	ndb, err := parseNewDbData(*newdbFlag)

	flags.Newdb = *ndb

	// If no flags were provided
	if pflag.NFlag() == 0 || err != nil {
		pflag.Usage()
		os.Exit(2)
	}

	flags.setFlagIfProvided(setup, SetupOff)
	flags.setFlagIfProvided(newdb, NewdbOff)
	flags.setFlagIfProvided(rmdb, RmdbOff)

	slog.Info("Successfully parsed flags.", "log_code", "fea120e1", "flags", flags)

	return flags
}
