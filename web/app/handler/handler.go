package handler

import (
	"SmartHomeWebCam/SHWService/web/app/logger"
	"SmartHomeWebCam/SHWService/web/app/middleware"
	"SmartHomeWebCam/SHWService/web/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Services *service.Service
}

func New(s *service.Service) *Handler {
	return &Handler{Services: s}
}

func (h *Handler) SetupRoutes(l *logger.Logger) *gin.Engine {
	router := gin.New()
	router.Use(gin.LoggerWithWriter(l.Log.Writer()))
	router.Use(middleware.Middleware())
	router.LoadHTMLGlob("web/template/*/*.html")
	router.Static("/assets", "./assets")

	router.GET("/", h.mainPage)
	auth := router.Group("/auth")
	{

		auth.POST("/sign-in", h.signIN)
		auth.POST("/sign-up", h.signUp)
		auth.POST("/logout", h.logout)
	}

	api := router.Group("/api")
	{
		api.GET("/", h.getProfile)
		api.GET("/help", h.getHelp)

		list := api.Group("/listCameras")
		{
			list.POST("/addCamera", h.addCamera)
			list.GET("/add", h.showFormAdd)
			list.GET("/notFound", h.showFormNotFound)
			list.GET("/", h.getListCameras)
			list.GET("/:ip/:port", h.getCameraView)
			list.GET("/:ip/:port/video", h.Video)
			list.GET(":ip/:port/setProcess", h.SetProcess)
			list.GET("/:ip/:port/delete", h.deleteCamera)
		}
		//listSmart := api.Group("/listSmartDevices")

	}
	return router
}
