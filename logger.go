package main

import (
	"log/slog"
	"os"
	"sync"
)

var singletonLogger *slog.Logger
var once sync.Once

func GetLogger() *slog.Logger {
	once.Do(func() {
		logOptions := &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}
		if IsDev() {
			logOptions.Level = slog.LevelDebug
		}
		singletonLogger = slog.New(slog.NewJSONHandler(os.Stdout, logOptions))
	})
	return singletonLogger
}
