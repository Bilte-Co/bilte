package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bilte-co/bilte/internal/logging"
	"github.com/bilte-co/bilte/internal/router"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/r3labs/sse/v2"
)

type WebCmd struct{}

func (cmd *WebCmd) Run(ctx *context.Context) error {
	logger := logging.NewLoggerFromEnv()

	err := godotenv.Load()
	if err != nil {
		logger.Debug("🤯 failed to load environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	e := echo.New()
	e.Use(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))

	if appEnv == "production" {
		e.HideBanner = true
		e.HidePort = true
		e.Pre(middleware.HTTPSRedirect())

		e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
			Level: 5,
		}))

		e.Use(middleware.Recover())
		e.Use(middleware.Secure())
	}

	e.Use(middleware.Logger())

	e.Static("/", "static")

	server := sse.New()       // create SSE broadcaster server
	server.AutoReplay = false // do not replay messages for each new subscriber that connects

	_ = server.CreateStream("feed")

	go func(s *sse.Server) {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.Publish("feed", &sse.Event{
					Data: []byte(time.Now().Format(time.RFC3339Nano)),
				})
			}
		}
	}(server)

	e = router.NewRouter(e, server)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))

	return nil
}
