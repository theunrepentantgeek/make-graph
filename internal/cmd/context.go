package cmd

import (
	"log/slog"

	"github.com/theunrepentantgeek/make-graph/internal/config"
)

// Flags contains shared state passed through kong commands.
type Flags struct {
	Verbose bool
	Log     *slog.Logger
	Config  *config.Config
}
