package auth

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
)

const WaitSignIn = time.Second * 10

func (auth auth) SignIn(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), WaitSignIn)
	defer cancel()
	auth.UseCase.SignIn(ctx)
	return nil
}
