package cmd

import (
	"log/slog"
	"os"

	"github.com/theunrepentantgeek/make-graph/internal/config"
)

// CLI represents the top-level command-line interface.
// This is a stub that will be fleshed out in a later task.
type CLI struct {
	Verbose bool   `kong:"short='v',help='Enable verbose output.'"`
	Config  string `kong:"short='c',type='existingfile',help='Path to configuration file.'"`
}

// Flags holds the resolved flags for subcommands.
type Flags struct {
	Verbose bool
	Log     *slog.Logger
	Config  *config.Config
}

// CreateLogger creates a logger based on CLI flags.
func (c *CLI) CreateLogger() *slog.Logger {
	level := slog.LevelInfo
	if c.Verbose {
		level = slog.LevelDebug
	}

	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
}

// CreateConfig loads the configuration file if specified.
func (c *CLI) CreateConfig() (*config.Config, error) {
	return config.New(), nil
}
