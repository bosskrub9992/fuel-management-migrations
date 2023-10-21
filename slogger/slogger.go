package slogger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/lmittmann/tint"
)

func New() *slog.Logger {
	slogHandler := tint.NewHandler(os.Stdout, &tint.Options{
		AddSource:  true,
		Level:      slog.LevelDebug,
		TimeFormat: time.StampMilli,
		NoColor:    false,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)
				a.Value = slog.StringValue(fmt.Sprintf("%s:%d", filepath.Base(source.File), source.Line))
				return a
			}
			return a
		},
	})
	return slog.New(slogHandler)
}
