package executors

import cp "github.com/bigkevmcd/go-configparser"

type Executor interface {
	Priority() int
	Execute(state *ExecutionState) (*ExecutionState, error)
}

type ExecutionState struct {
	Config *cp.ConfigParser
}
