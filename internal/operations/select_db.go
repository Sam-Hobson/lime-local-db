package operations

import "log/slog"

func SelectDb(db string) {
    slog.Info("Setting selected database.", "log_code", "cc2de9b6", "db", db)
	mutableApplicationState.setSelectedDatabase(db)
}
