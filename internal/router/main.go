package router

import (
	"github.com/labstack/echo/v4"
	"github.com/r3labs/sse/v2"
)

func NewRouter(e *echo.Echo, sseServer *sse.Server) *echo.Echo {
	e.GET("/sse", func(c echo.Context) error {
		country := c.Request().Header.Get("CF-IPCountry")

		go func(c string) {
			if c == "" {
				return
			}

			sseServer.Publish("feed", &sse.Event{
				Data: []byte("Your country is: " + c),
			})
		}(country)

		go func() {
			<-c.Request().Context().Done() // Received Browser Disconnection
			return
		}()

		sseServer.ServeHTTP(c.Response(), c.Request())

		return nil
	})

	return e
}
