package main

import (
	"fmt"
	"log/slog"
	"os"

	. "github.com/sam-hobson/internal/command"
	. "github.com/sam-hobson/internal/config"

	"github.com/go-errors/errors"
)

func main() {
	home := os.Getenv("HOME")
	args := os.Args[1:]

	config, err := ParseConfig(home)
	state := &ExecutionState{Config: config}

	if err != nil && args[0] != "setup" {
		panicErr(err)
	}

	executors, err := ProcessArgs(args)
	panicErr(err)

	for _, exec := range executors {
		slog.Info("Executing executor.", "Priority", exec.Priority())

		state, err = exec.Execute(state)
		panicErr(err)
	}

}

func panicErr(err error) {
	if err != nil {
		fmt.Println(err.(*errors.Error).ErrorStack())
		panic(err)
	}
}
