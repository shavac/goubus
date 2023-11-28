package goubus

import (
	"log/slog"
)

var log = slog.Default()

func SetLogLevel(lvl slog.Level) {
	slog.SetLogLoggerLevel(lvl)
}
