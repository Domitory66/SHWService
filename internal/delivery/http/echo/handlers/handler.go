package handlers

import (
	"github.com/labstack/echo/v4"
)

type Auth interface {
	SignIn(c echo.Context) error
	SignUp(c echo.Context) error
	LogOut(c echo.Context) error
}

func MainPage(c echo.Context) error {
	return c.Render(200, "login.html", nil)
}

type Handler struct {
	Auth
}

func New(auth Auth) Handler {
	return Handler{Auth: auth}
}
