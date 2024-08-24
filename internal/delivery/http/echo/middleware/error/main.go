package errmidleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type errorHandler struct {
}

func (eh errorHandler) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			if he, ok := err.(*echo.HTTPError); ok {
				c.JSON(he.Code, map[string]interface{}{"error": he.Message})
			} else {
				c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
			}
		}
		return nil
	}
}

func New() errorHandler {
	return errorHandler{}
}
