package operations

import "log/slog"

func SelectCol(col string) {
    slog.Info("Setting selected column.", "log_code", "ab0777a5", "Column", col)
    mutableApplicationState.setSelectedColumn(col)
}
