package main

import (
	"context"
	"os"

	"github.com/alecthomas/kong"
	"github.com/bilte-co/bilte/cmd/web"
	"github.com/bilte-co/bilte/internal/logging"
)

type Cli struct {
	Web web.WebCmd `cmd:"" help:"Start the web server."`
}

func main() {
	// Display help if no args are provided instead of an error message
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "--help")
	}

	app := Cli{}

	logger := logging.NewLoggerFromEnv()
	cntx := context.Background()

	ctx := kong.Parse(&app,
		kong.Name("bilte"),
		kong.Description("bilte CLI"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
	)

	logging.WithLogger(cntx, logger)

	err := ctx.Run(&cntx)
	ctx.FatalIfErrorf(err)
}
