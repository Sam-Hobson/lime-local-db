package main

import (
	"fmt"
	. "github.com/sam-hobson/internal/command"
	. "github.com/sam-hobson/internal/config"
	"os"

	"github.com/go-errors/errors"
)

func main() {
	_, err := ParseConfig()
	panicErr(err)
	err = ProcessArgs(os.Args[1:])
	panicErr(err)
}

func panicErr(err error) {
	if err != nil {
		fmt.Println(err.(*errors.Error).ErrorStack())
		panic(err)
	}
}
