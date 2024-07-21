package cli

import (
	"log/slog"
	"math"
	"strconv"
	"strings"

	"github.com/go-errors/errors"
	op "github.com/sam-hobson/internal/operations"
	"github.com/spf13/cobra"
)

const NewDbCommandUsage = `new-db [Database name] [Database columns]...

A Database column is formatted as follows:

[Key flags][Not null]:[Column type]:[Column name]{[Default value]}


[Key flags] is either P for primary key, or F for foreign key.
[Not null] is N if the column should not be NULL.

[Column type] is the data type of the column, one of:
- INT
- REAL
- STR          -- Text field of unlimited length.
- STR{length}  -- Text field of fixed length.
- BLOB

[Column type] can be left blank, in which case the column will have a dynamic type.

[Default value] is the value of the column if a value is not specified.
`

const NewDbCommandExample = `limedb new-db petdb P:STR:name{default} N:STR:gender{F} N::breed{Dog}

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

	keyFlags := parts[0]
	colType := parts[1]
	nameAndDefaultVal := parts[2]

	if err := parseKeyFlags(keyFlags, column); err != nil {
		return nil, err
	}
	if err := parseDataType(colType, column); err != nil {
		return nil, err
	}
	if err := parseNameAndDefaultVal(nameAndDefaultVal, column); err != nil {
		return nil, err
	}

	return column, nil
}

func parseKeyFlags(flags string, column *op.Column) error {
	var primaryKey = false
	var foreignKey = false
	var notNull = false

	flags = strings.ToUpper(flags)

	for _, flag := range flags {
		switch flag {
		case 'P':
			if primaryKey {
				slog.Error("Column entry malformed.", "log_code", "759534a7", "Key_flags", flags)
				return errors.Errorf("Found malformed column entry input (P used more than once): %s", flags)
			}
			primaryKey = true
		case 'F':
			if foreignKey {
				slog.Error("Column entry malformed.", "log_code", "ec05b044", "Key_flags", flags)
				return errors.Errorf("Found malformed column entry input (F used more than once): %s", flags)
			}
			foreignKey = true
		case 'N':
			if notNull {
				slog.Error("Column entry malformed.", "log_code", "41383bf1", "Key_flags", flags)
				return errors.Errorf("Found malformed column entry input (N used more than once): %s", flags)
			}
			notNull = true
		default:
			slog.Error("Column entry malformed.", "log_code", "9288e4b5", "Key_flags", flags)
			return errors.Errorf("Found malformed key flags on column entry input: %s", flags)
		}
	}

	column.PrimaryKey = primaryKey
	column.ForeignKey = foreignKey
	column.NotNull = notNull

	return nil
}

func parseDataType(dataType string, column *op.Column) error {
	dataType = strings.ToUpper(dataType)

	// re := regexp.MustCompile(`STR{(.*)}`)

	if dataType == "INT" {
		column.DataType = op.ColumnIntDataType
	} else if dataType == "REAL" {
		column.DataType = op.ColumnRealDataType
	} else if dataType == "STR" {
		column.DataType = op.ColumnTextDataType
	} else if dataType == "BLOB" {
		column.DataType = op.ColumnBlobCharDataType
	} else if strings.HasPrefix(dataType, "STR{") && strings.HasSuffix(dataType, "}") {
		endLenIndex := strings.IndexRune(dataType, '}')

		if endLenIndex != len(dataType)-1 {
			slog.Error("Column entry malformed.", "log_code", "5acc18ef", "Data_types", dataType)
			return errors.Errorf("Found malformed column data type input: %s", dataType)
		}

		lengthStr := dataType[len("STR{") : endLenIndex]

        if lengthStr == "" {
            column.DataType = op.ColumnTextDataType
            return nil
        }

		length, err := strconv.ParseUint(lengthStr, 10, 32)

		if err != nil {
			slog.Error("Column entry malformed.", "log_code", "699712b6", "Data_types", dataType)
			return errors.Errorf("Found malformed column data type input (invalid STR length): %s", dataType)
		}

		column.DataType = op.ColumnVarCharDataType
		column.VarCharLength = uint32(length)
	} else if dataType == "" {
        column.DataType = op.ColumnNullDataType
	} else {
		slog.Error("Column entry malformed.", "log_code", "aaf65198", "Data_types", dataType)
		return errors.Errorf("Found malformed column data type input: %s", dataType)
	}

	return nil
}

func parseNameAndDefaultVal(nameAndDefaultVal string, column *op.Column) error {
    if nameAndDefaultVal == "" {
		slog.Error("Column entry malformed.", "log_code", "ea5cd3fa", "Name", nameAndDefaultVal)
		return errors.Errorf("Found malformed column name input (a column name must be provided): %s", nameAndDefaultVal)
    }

	startDefaultValIndex := strings.IndexRune(nameAndDefaultVal, '{')

	if startDefaultValIndex == -1 {
		column.ColName = nameAndDefaultVal
		return nil
	}

	endDefaultValIndex := strings.IndexRune(nameAndDefaultVal, '}')

	if (startDefaultValIndex == 0) || (endDefaultValIndex != len(nameAndDefaultVal)-1) {
		slog.Error("Column entry malformed.", "log_code", "d63efff4", "Name", nameAndDefaultVal)
		return errors.Errorf("Found malformed column name input: %s", nameAndDefaultVal)
	}

	column.ColName = nameAndDefaultVal[:startDefaultValIndex]
	column.DefaultVal = nameAndDefaultVal[startDefaultValIndex+1 : endDefaultValIndex]

	return nil
}
