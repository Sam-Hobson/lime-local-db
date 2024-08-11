package types

import (
	"strings"
	"unicode"

	"github.com/go-errors/errors"
	"github.com/sam-hobson/internal/util"
)

type Column struct {
	Name           string
	DataType       ColumnDataType
	PrimaryKey     bool
	ForeignKey     bool
	NotNull        bool
	AutoIncrememnt bool
	DefaultVal     string
}

func (c *Column) String() string {
	sb := strings.Builder{}

	if c.PrimaryKey {
		sb.WriteRune('P')
	}
	if c.ForeignKey {
		sb.WriteRune('F')
	}
	if c.AutoIncrememnt {
		sb.WriteRune('A')
	}
	if c.NotNull {
		sb.WriteRune('N')
	}

	sb.WriteRune(':')
	sb.WriteString(c.DataType.String())
	sb.WriteRune(':')
	sb.WriteString(c.Name)

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
		util.Log("b7144d7b").Error("Column entry malformed.", "Column", col)
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
	var autoIncrememnt = false

	for _, flag := range flags {
		switch unicode.ToUpper(flag) {
		case 'P':
			if primaryKey {
				util.Log("759534a7").Warn("Column entry malformed.", "Flags", flags)
				return errors.Errorf("Found malformed column entry input (P used more than once): %s", flags)
			}
			primaryKey = true
		case 'F':
			if foreignKey {
				util.Log("ec05b044").Warn("Column entry malformed.", "Flags", flags)
				return errors.Errorf("Found malformed column entry input (F used more than once): %s", flags)
			}
			foreignKey = true
		case 'N':
			if notNull {
				util.Log("41383bf1").Warn("Column entry malformed.", "Flags", flags)
				return errors.Errorf("Found malformed column entry input (N used more than once): %s", flags)
			}
			notNull = true
		case 'A':
			if autoIncrememnt {
				util.Log("98fccd5e").Warn("Column entry malformed.", "Flags", flags)
				return errors.Errorf("Found malformed column entry input (A used more than once): %s", flags)
			}
			autoIncrememnt = true
		default:
			util.Log("9288e4b5").Error("Column entry malformed.", "Key flags", flags)
			return errors.Errorf("Found malformed key flags on column entry input: %s", flags)
		}
	}

	column.PrimaryKey = primaryKey
	column.ForeignKey = foreignKey
	column.NotNull = notNull
	column.AutoIncrememnt = autoIncrememnt

	return nil
}

func parseColumnDataType(dataType string, column *Column) error {
	if dt, err := NewDataType(dataType); err != nil {
		util.Log("7da9f304").Error("Failed to parse data type", "Data type", dataType)
		return err
	} else {
		column.DataType = dt
		return nil
	}
}

func parseColumnNameAndDefaultVal(nameAndDefaultVal string, column *Column) error {
	if nameAndDefaultVal == "" {
		util.Log("ea5cd3fa").Error("Column entry malformed.", "Name", nameAndDefaultVal)
		return errors.Errorf("Found malformed column name input (a column name must be provided): %s", nameAndDefaultVal)
	}

	startDefaultValIndex := strings.IndexRune(nameAndDefaultVal, '{')

	if startDefaultValIndex == -1 {
		column.Name = nameAndDefaultVal
		return nil
	}

	endDefaultValIndex := strings.IndexRune(nameAndDefaultVal, '}')

	if (startDefaultValIndex == 0) || (endDefaultValIndex != len(nameAndDefaultVal)-1) {
		util.Log("d63efff4").Error("Column entry malformed.", "Name", nameAndDefaultVal)
		return errors.Errorf("Found malformed column name input: %s", nameAndDefaultVal)
	}

	column.Name = nameAndDefaultVal[:startDefaultValIndex]
	column.DefaultVal = nameAndDefaultVal[startDefaultValIndex+1 : endDefaultValIndex]

	return nil
}
