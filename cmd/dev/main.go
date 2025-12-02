package main

import (
	"log/slog"

	"github.com/Nemagu/dnd/internal/config"
	"github.com/Nemagu/dnd/internal/logger/sl"
	"github.com/Nemagu/dnd/internal/port/http/web"
)

func main() {
	cfg := config.MustNewWebConfig()
	logLevel := slog.LevelInfo
	if cfg.Debug {
		logLevel = slog.LevelDebug
	}
	logger := sl.MustNewJSONLogger(logLevel)
	server := web.MustNewHTTPServer(logger, cfg)
	server.MustServe()
}
