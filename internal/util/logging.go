package util

import (
	"fmt"
	"log/slog"
	"os"
)

var sessionId int64 = -1

func SetSessionId(id int64) {
	sessionId = id
}

func groupLogId(logHash string) slog.Attr {
	if sessionId == -1 {
		fatal("Cannot proceed as session id was not initialised.")
	}

	return slog.Group("Id", slog.Int64("Session", sessionId), slog.String("Log", logHash))
}

func Log(logHash string) *slog.Logger {
    return slog.With(groupLogId(logHash))
}

func fatal(s string, args ...any) {
	fmt.Printf(s+"\n", args...)
	os.Exit(1)
}
