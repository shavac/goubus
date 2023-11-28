package goubus

import (
	"log/slog"
)

var log = slog.Default()

func SetLoglevel(lvl slog.Level) {
	slog.SetLogLoggerLevel(slog.LevelWarn)
}
