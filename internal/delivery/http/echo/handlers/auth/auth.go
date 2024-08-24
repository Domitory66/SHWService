package auth

import (
	"shwservice/internal/usecase"

	"github.com/labstack/echo/v4"
)

type Auth interface {
	SignIn(c echo.Context) error
	SignUp(c echo.Context) error
	LogOut(c echo.Context) error
}

type auth struct {
	UseCase usecase.Auth
}

func New(uc usecase.Auth) Auth {
	return auth{UseCase: uc}
}
