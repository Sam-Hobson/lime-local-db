package state

import "log/slog"

type state struct {
	selectedDatabase string
}

func (s state) GetSelectedDb() string {
	return s.selectedDatabase
}

func (s *state) SetSelectedDb(db string) {
	slog.Info("Setting selected database.", "Log code", "c863f749", "Db", db)
	s.selectedDatabase = db
}

var globalState = &state{
	selectedDatabase: "",
}

func ApplicationState() *state {
	return globalState
}
