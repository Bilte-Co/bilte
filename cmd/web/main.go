package web

import (
	"context"
	"html/template"
	"os"
	"strings"

	"github.com/bilte-co/bilte/internal/logging"
	"github.com/bilte-co/bilte/internal/router"
	"github.com/bilte-co/bilte/internal/templates"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	sloggin "github.com/samber/slog-gin"
)

type WebCmd struct{}

func safeHTML(s string) template.HTML {
	return template.HTML(s)
}

func staticCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/public") {
			c.Header("Cache-Control", "public, max-age=31536000")
		}
		c.Next()
	}
}

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

	isProduction := appEnv == "production"

	if isProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	natsOpts := []nats.Option{}

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	natsToken := os.Getenv("NATS_TOKEN")
	if natsToken != "" {
		logger.Info("ℹ️ using NATS token")
		natsOpts = append(natsOpts, nats.Token(natsToken))
	}

	nc, err := nats.Connect(natsURL, natsOpts...)
	if err != nil {
		return err
	}

	r := gin.Default()
	config := sloggin.Config{
		WithSpanID:  true,
		WithTraceID: true,
	}
	r.Use(sloggin.NewWithConfig(logger, config))

	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{"127.0.0.1"})

	if gin.Mode() == gin.ReleaseMode {
		r.Use(staticCacheMiddleware())
		r.Use(gin.Recovery())
	}

	ginHtmlRenderer := r.HTMLRender
	r.HTMLRender = &templates.HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}

	r.Static("/public", "static")

	r = router.NewRouter(r, &isProduction, nc)

	r.Run(":" + port)

	return nil
}
