package main

import (
	"log/slog"

	"github.com/bosskrub9992/fuel-management-migrations/config"
	"github.com/bosskrub9992/fuel-management-migrations/slogger"
)

func main() {
	cfg := config.New()
	slog.SetDefault(slogger.New())

	slog.Info("test config", slog.Any("cfg", cfg))
}
