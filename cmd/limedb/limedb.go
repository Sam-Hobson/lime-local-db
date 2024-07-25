package main

import (
	"fmt"

	"github.com/go-errors/errors"
	"github.com/huandu/go-sqlbuilder"

	. "github.com/sam-hobson/internal/cli"
)

func main() {
    // TODO: Load config if it exist

    sqlbuilder.DefaultFlavor = sqlbuilder.SQLite

	ProcessArgs()
}

func panicErr(err error) {
	if err != nil {
		fmt.Println(err.(*errors.Error).ErrorStack())
		panic(err)
	}

}
