package database

import (
	"strings"

	"github.com/go-errors/errors"
)

type ColumnDataType byte

func (c ColumnDataType) String() string {
	return columnDataTypeName[c]
}

const (
	ColumnIntDataType      = ColumnDataType(0)
	ColumnRealDataType     = ColumnDataType(1)
	ColumnTextDataType     = ColumnDataType(2)
	ColumnBlobCharDataType = ColumnDataType(3)
	ColumnNullDataType     = ColumnDataType(4)
)

var columnDataTypeName = [...]string{
	"INT",
	"REAL",
	"TEXT",
	"BLOB",
	"NULL",
}

func NewDataType(dataType string) (ColumnDataType, error) {
	switch strings.ToUpper(dataType) {
	case "INT":
		return ColumnIntDataType, nil
	case "REAL":
		return ColumnRealDataType, nil
	case "TEXT":
		return ColumnTextDataType, nil
	case "BLOB":
		return ColumnBlobCharDataType, nil
	case "NULL":
		return ColumnNullDataType, nil
	default:
		return ColumnNullDataType, errors.Errorf("Data type %s does not exist", dataType)
	}
}
