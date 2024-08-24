package loggermiddleware

import (
	"shwservice/pkg/logger"

	"github.com/labstack/echo/v4"
)

type log struct {
	Log logger.Logger
}

func (l log) Logging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		l.Log.Info("Request received:",
			logger.Args{
				"URI:":    c.Request().RequestURI,
				"Body":    c.Request().Body,
				"IP":      c.RealIP(),
				"Cookies": c.Request().Cookies(),
			})
		return nil
	}
}

func New(l logger.Logger) log {
	return log{Log: l}
}
