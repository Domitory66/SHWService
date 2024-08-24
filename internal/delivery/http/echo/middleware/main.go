package middleware

import (
	"io"

	"github.com/labstack/echo/v4"
)

type HandleError interface {
	Handle(next echo.HandlerFunc) echo.HandlerFunc
}

type Logger interface {
	Logging(next echo.HandlerFunc) echo.HandlerFunc
}

type Render interface {
	Render(io.Writer, string, interface{}, echo.Context) error
}

type middleware struct {
	HandleError
	Logger
	Render
}

func New(he HandleError, l Logger, r Render) middleware {
	return middleware{HandleError: he, Logger: l, Render: r}
}
