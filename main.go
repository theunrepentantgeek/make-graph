package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/rotisserie/eris"

	"github.com/theunrepentantgeek/make-graph/internal/cmd"
)

func main() {
	var cli cmd.CLI

	ctx := kong.Parse(
		&cli,
		kong.UsageOnError(),
	)

	log := cli.CreateLogger()
	cfg, err := cli.CreateConfig()
	if err != nil {
		log.Error(eris.ToString(err, true))
		ctx.Exit(1)
	}

	flags := &cmd.Flags{
		Verbose: cli.Verbose,
		Log:     log,
		Config:  cfg,
	}

	err = ctx.Run(flags)
	if err != nil {
		log.Error(eris.ToString(err, true))
		ctx.Exit(1)
	}

	os.Exit(0)
}
