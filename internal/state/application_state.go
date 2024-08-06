package state

import (
	"github.com/sam-hobson/internal/util"
)

type state struct {
	selectedDatabase string
	limedbHome       string
}

func (s *state) GetSelectedDb() string {
	return s.selectedDatabase
}

func (s *state) SetSelectedDb(databaseName string) {
	util.Log("5c4edaeb").Info("Setting selected database.", "Database name", databaseName)
	s.selectedDatabase = databaseName
}

func (s *state) GetLimedbHome() string {
	return s.limedbHome
}

func (s *state) SetLimedbHome(home string) {
	util.Log("0d7bad86").Info("Setting limedb home.", "Home", home)
	s.limedbHome = home
}

var globalState = &state{
	selectedDatabase: "",
	limedbHome:       "",
}

func ApplicationState() *state {
	return globalState
}
