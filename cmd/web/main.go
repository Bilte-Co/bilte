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
)

type WebCmd struct{}

func (cmd *WebCmd) Run(ctx *context.Context) error {
	logger := logging.NewLoggerFromEnv()

	err := godotenv.Load()
	if err != nil {
		logger.Debug().Err(err).Msg("🤯 failed to load environment variables")
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

		// e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(100))))

		e.Use(middleware.Recover())
		e.Use(middleware.Secure())
	}

	e.Use(middleware.Logger())

	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "the request has timed out",
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			_ = c.String(http.StatusRequestTimeout, "the request has timed out")
		},
		Timeout: 30 * time.Second,
	}))

	e.Static("/", "static")

	e = router.NewRouter(e)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))

	return nil
}
