package web

import (
	"context"
	"html/template"
	"os"
	"strings"

	"github.com/bilte-co/bilte/internal/logging"
	"github.com/bilte-co/bilte/internal/router"
	"github.com/bilte-co/bilte/internal/templates"
	"github.com/joho/godotenv"
	sloggin "github.com/samber/slog-gin"

	"github.com/gin-gonic/gin"
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

	// engine.HTMLRender = gintemplrenderer.Default

	ginHtmlRenderer := r.HTMLRender
	r.HTMLRender = &templates.HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}

	// e := echo.New()
	// e.Use(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
	// 	RedirectCode: http.StatusMovedPermanently,
	// }))

	// if isProduction {
	// 	e.HideBanner = true
	// 	e.HidePort = true
	// 	e.Pre(middleware.HTTPSRedirect())

	// 	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
	// 		Level: 5,
	// 	}))

	// 	e.Use(middleware.Recover())
	// 	e.Use(middleware.Secure())
	// }
	//

	r.Static("/public", "static")

	// server := sse.New()       // create SSE broadcaster server
	// server.AutoReplay = false // do not replay messages for each new subscriber that connects

	// _ = server.CreateStream("feed")

	// go func(s *sse.Server) {
	// 	ticker := time.NewTicker(5 * time.Second)
	// 	defer ticker.Stop()

	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			s.Publish("feed", &sse.Event{
	// 				Data: []byte(time.Now().UTC().Format(time.RFC3339)),
	// 			})
	// 		}
	// 	}
	// }(server)

	r = router.NewRouter(r, &isProduction)

	// e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
	r.Run(":" + port)

	return nil
}
