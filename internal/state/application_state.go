package state

import (
	"github.com/sam-hobson/internal/util"
)

type state struct {
	selectedDatabase string
}

func (s state) GetSelectedDb() string {
	return s.selectedDatabase
}

func (s *state) SetSelectedDb(db string) {
    util.Log("5c4edaeb").Info("Setting selected database.", "Db", db)
	s.selectedDatabase = db
}

var globalState = &state{
	selectedDatabase: "",
}

func ApplicationState() *state {
	return globalState
}
