package main

import (
	"fmt"
	"github.com/go-errors/errors"

    . "github.com/sam-hobson/internal/cli"
)

func main() {
    flags := GetFlags()

    ProcessArgs(flags)
}

func panicErr(err error) {
	if err != nil {
		fmt.Println(err.(*errors.Error).ErrorStack())
		panic(err)
	}
}
