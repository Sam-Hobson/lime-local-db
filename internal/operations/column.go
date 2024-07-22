package operations

import "fmt"

type ColumnDataType byte

func (c ColumnDataType) String() string {
	return columnDataTypeName[c]
}

const (
	ColumnIntDataType      = ColumnDataType(0)
	ColumnRealDataType     = ColumnDataType(1)
	ColumnTextDataType     = ColumnDataType(2)
	ColumnVarCharDataType  = ColumnDataType(3)
	ColumnBlobCharDataType = ColumnDataType(4)
	ColumnNullDataType     = ColumnDataType(5)
)

var columnDataTypeName = [...]string{
	"INT",
	"REAL",
	"TEXT",
	"VARCHAR",
	"BLOB",
	"NULL",
}

type Column struct {
	ColName       string
	DataType      ColumnDataType
	VarCharLength uint32
	PrimaryKey    bool
	ForeignKey    bool
	NotNull       bool
	DefaultVal    any // TODO: Update this to a real type
}

func (c *Column) String() string {
	return fmt.Sprintf("%+v", *c)
}

// TODO: Default values for colName and defaultVal
func NewColumn() *Column {
	return &Column{
		DataType:      ColumnNullDataType,
		VarCharLength: 0,
		PrimaryKey:    false,
		ForeignKey:    false,
		NotNull:       false,
	}
}
