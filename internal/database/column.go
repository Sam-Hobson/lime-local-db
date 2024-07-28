package database

import (
	"log/slog"
	"strings"
	"unicode"

	"github.com/go-errors/errors"
)

type Column struct {
	ColName    string
	DataType   ColumnDataType
	PrimaryKey bool
	ForeignKey bool
	NotNull    bool
	DefaultVal string
}

func (c *Column) String() string {
	sb := strings.Builder{}

	if c.PrimaryKey {
		sb.WriteRune('P')
	}
	if c.ForeignKey {
		sb.WriteRune('F')
	}
	if c.NotNull {
		sb.WriteRune('N')
	}

	sb.WriteRune(':')
	sb.WriteString(c.DataType.String())
	sb.WriteRune(':')
	sb.WriteString(c.ColName)

	if c.DefaultVal != "" {
		sb.WriteRune('{')
		sb.WriteString(c.DefaultVal)
		sb.WriteRune('}')
	}

	return sb.String()
}

func ParseColumnString(col string) (*Column, error) {
	parts := strings.Split(col, ":")

	if len(parts) != 3 {
		slog.Error("Column entry malformed.", "log_code", "b7144d7b", "Column", col)
		return nil, errors.Errorf("Found malformed column entry input: %s", col)
	}

	column := &Column{}

	keyFlags := parts[0]
	colType := parts[1]
	nameAndDefaultVal := parts[2]

	if err := parseColumnFlags(keyFlags, column); err != nil {
		return nil, err
	}
	if err := parseColumnDataType(colType, column); err != nil {
		return nil, err
	}
	if err := parseColumnNameAndDefaultVal(nameAndDefaultVal, column); err != nil {
		return nil, err
	}

	if (column.DataType == ColumnTextDataType) && (column.DefaultVal != "") {
		column.DefaultVal = "'" + column.DefaultVal + "'"
	}

	return column, nil
}

func parseColumnFlags(flags string, column *Column) error {
	var primaryKey = false
	var foreignKey = false
	var notNull = false

	for _, flag := range flags {
		switch unicode.ToUpper(flag) {
		case 'P':
			if primaryKey {
				slog.Error("Column entry malformed.", "log_code", "759534a7", "flags", flags)
				return errors.Errorf("Found malformed column entry input (P used more than once): %s", flags)
			}
			primaryKey = true
		case 'F':
			if foreignKey {
				slog.Error("Column entry malformed.", "log_code", "ec05b044", "flags", flags)
				return errors.Errorf("Found malformed column entry input (F used more than once): %s", flags)
			}
			foreignKey = true
		case 'N':
			if notNull {
				slog.Error("Column entry malformed.", "log_code", "41383bf1", "flags", flags)
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

func parseColumnDataType(dataType string, column *Column) error {
	if dt, err := NewDataType(dataType); err != nil {
		slog.Error("Failed to parse data type", "log_code", "7da9f304", "data_type", dataType)
		return err
	} else {
		column.DataType = dt
		return nil
	}
}

func parseColumnNameAndDefaultVal(nameAndDefaultVal string, column *Column) error {
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
