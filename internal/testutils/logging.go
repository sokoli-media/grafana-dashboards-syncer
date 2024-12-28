package testutils

import (
	"log/slog"
	"os"
)

var LoggerForTesting = slog.New(slog.NewJSONHandler(os.Stdout, nil))
