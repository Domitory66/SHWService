package routes

import (
	"shwservice/internal/delivery/http/echo/handlers/auth"
	"shwservice/internal/delivery/http/echo/middleware"
	errmidleware "shwservice/internal/delivery/http/echo/middleware/error"
	loggermiddleware "shwservice/internal/delivery/http/echo/middleware/logger"
	"shwservice/internal/delivery/http/echo/middleware/render"
	"shwservice/internal/usecase"
	"shwservice/pkg/logger"

	"github.com/labstack/echo/v4"
)

func SetRoutes(l logger.Logger, uc usecase.UseCase) *echo.Echo {
	router := echo.New()

	// middlewares
	log := loggermiddleware.New(l)
	errHandler := errmidleware.New()
	render := render.New("web/template")

	middleware := middleware.New(errHandler, log, &render)
	router.Use(middleware.Logging)
	router.Use(middleware.Handle)
	router.Renderer = middleware.Render

	router.Static("/assets", "./assets")

	// handlers
	authHandler := auth.New(uc.Auth)

	// auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/sign-in", authHandler.SignIn)
		auth.POST("/sign-up", authHandler.SignUp)
		auth.POST("/logout", authHandler.LogOut)
	}

	// user routes

	// device routes

	// home routes

	return router
}
