package cli

import (
	"log/slog"
	"math"
	"strings"

	"github.com/go-errors/errors"
	op "github.com/sam-hobson/internal/operations"
	"github.com/spf13/cobra"
)

const NewDbCommandUsage = `new-db [Database name] [Database columns]...

A Database column is formatted as follows:

[Key flags][Nullable]:[Column type]:[Column name]{[Default value]}


[Key flags] is either P for primary key, or F for foreign key.
[Nullable] is N if the column should be Nullable.

[Column type] is the data type of the column, one of:
- INT
- REAL
- STR          -- Text field of unlimited length.
- STR{length}  -- Text field of fixed length.
- BLOB

[Column type] can be left blank, in which case the column will have a dynamic type.

[Default value] is the value of the column if a value is not specified.
`

const NewDbCommandExample = `limedb new-db petdb P:TEXT:name{default} N:TEXT:gender{F} N::breed{Dog}

limedb new-db petdb P:INT:age ::gender{M} FN:INT:id

limedb new-db petdb "P:STR{50}:name{a puppy}" "N:STR:goodBoy{Yes hes a good boy}"
`

var NewDbCommand = &cobra.Command{
	Use:     NewDbCommandUsage,
	Example: NewDbCommandExample,
	Short:   "Create a new database",
	Args:    cobra.RangeArgs(2, math.MaxInt8),

	DisableFlagsInUseLine: true, // TODO: Remove this when flags are added

	RunE: processNewDbCmd,
}

func processNewDbCmd(cmd *cobra.Command, args []string) error {
	// TODO: Add support for various flags
	anyFlagSet := cmd.Flags().NFlag() != 0

	if anyFlagSet {
		slog.Error("new-db does not take any flags.", "log_code", "998073c6")
		return errors.Errorf("new-db does not take any flags.")
	}

	name := args[0]
	attrs := args[1:]

	columns := make([]*op.Column, len(attrs))

	// Parse columns
	for i, attr := range attrs {
		column, err := parseColEntry(attr)

		if err != nil {
			return err
		}

		columns[i] = column
	}

	err := op.NewDb(name, columns)

	if err != nil {
		slog.Error("new-db command failed.", "log_code", "0f36afa0")
		return err
	}

	return nil
}

func parseColEntry(colEntry string) (*op.Column, error) {
	parts := strings.Split(colEntry, ":")

	if len(parts) != 3 {
		slog.Error("Column entry malformed.", "log_code", "b7144d7b", "Column", colEntry)
		return nil, errors.Errorf("Found malformed column entry input: %s", colEntry)
	}

    column := op.NewColumn()

	keyFlagsOrNullable := parts[0]
	// colType := parts[1]
	// nameAndDefaultVal := parts[2]

    err := parseKeyFlags(keyFlagsOrNullable, column)

    if err != nil {
        return nil, err
    }

	return nil, nil
}

func parseKeyFlags(flags string, column *op.Column) error {
	var primaryKey = false
	var foreignKey = false
	var nullable = false

	for _, flag := range flags {
		switch flag {
		case 'P':
			if primaryKey {
				slog.Error("Column entry malformed.", "log_code", "759534a7", "Column", flags)
				return errors.Errorf("Found malformed column entry input (P used more than once): %s", flags)
			}
            primaryKey = true
		case 'F':
			if foreignKey {
				slog.Error("Column entry malformed.", "log_code", "ec05b044", "Column", flags)
				return errors.Errorf("Found malformed column entry input (F used more than once): %s", flags)
			}
            foreignKey = true
		case 'N':
			if nullable {
				slog.Error("Column entry malformed.", "log_code", "41383bf1", "Column", flags)
				return errors.Errorf("Found malformed column entry input (N used more than once): %s", flags)
			}
            nullable = true
		default:
			slog.Error("Column entry malformed.", "log_code", "9288e4b5", "Column", flags)
			return errors.Errorf("Found malformed key flags/nullable on column entry input: %s", flags)
		}
	}

    column.PrimaryKey = primaryKey
    column.ForeignKey = foreignKey
    column.Nullable = nullable

    return nil
}
