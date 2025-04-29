package git

import (
	"context"

	"github.com/bilte-co/bilte/internal/logging"
	"github.com/joho/godotenv"
)

type GitCmd struct {
	Directory string `help:"Directory to search for git repositories." short:"d"`
	APIKey    string `help:"API key for authentication." short:"k"`
}

func (cmd *GitCmd) Run(ctx *context.Context) error {
	logger := logging.NewLoggerFromEnv()

	err := godotenv.Load()
	if err != nil {
		logger.Debug("🤯 failed to load environment variables")
	}

	if cmd.Directory != "" {
		logger.Info("🗂️ Directory provided, searching for git repositories...", "directory", cmd.Directory)

		// call Local function
		if err := cmd.Local(cmd.Directory); err != nil {
			logger.Error("🛑 Error searching local repositories", "error", err)
			return err
		}
	}

	if cmd.APIKey != "" {
		logger.Info("🌎 API Key provided, searching through available remote repositories...")

		// call Remote function
		if err := cmd.Remote(cmd.APIKey); err != nil {
			logger.Error("🛑 Error searching remote repositories", "error", err)
			return err
		}
	}

	return nil
}

func (cmd *GitCmd) Local(dir string) error {
	return nil
}

func (cmd *GitCmd) Remote(apiKey string) error {
	return nil
}
