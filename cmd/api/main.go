package main

import (
	"errors"
	"log"
	"net"
	"os"
	"shwservice/internal/common/config"
	cameraworker "shwservice/internal/delivery/grpc/camera"
	deviceworker "shwservice/internal/delivery/grpc/device"
	userworker "shwservice/internal/delivery/grpc/user"
	videoworker "shwservice/internal/delivery/grpc/video"
	"shwservice/internal/delivery/http/echo/routes"
	"shwservice/internal/delivery/http/server"
	"shwservice/internal/usecase"
	"shwservice/internal/usecase/auth"
	"shwservice/pkg/logger"

	"github.com/joho/godotenv"
)

var (
	errInitConfig         = errors.New("can't init config")
	errUserService        = errors.New("initialization user service failed")
	errCameraService      = errors.New("initialization camera service failed")
	errVideoStreamService = errors.New("initialization video stream service failed")
	errDeviceService      = errors.New("initialization device service failed")
	errInitLogger         = errors.New("can't init logger")
)

var (
	msgInfoStartWithoutLocalLogger  = "Start service without local logger."
	msgInfoStartWithoutRemoteLogger = "Start service without remote logger."
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg, err := config.New(os.Getenv("configPath"))
	if err != nil {
		log.Fatalf(errInitConfig.Error(), err.Error())
	}

	f, err := os.Open(cfg.Logger.Path)
	if err != nil {
		log.Print(msgInfoStartWithoutLocalLogger)
	}

	remoteLogger, err := net.Dial("tcp", cfg.Logger.Address)
	if err != nil {
		log.Print(msgInfoStartWithoutRemoteLogger)
	}

	loggerWriter, err := logger.New(f, remoteLogger, cfg.IsDebug)
	if err != nil {
		log.Fatalf(errInitLogger.Error(), err.Error())
	}

	// Connect to user service
	userWorker, err := userworker.New(cfg.Addresses.User)
	if err != nil {
		loggerWriter.Error(errUserService, logger.Args{"Error": err.Error()})
	}

	// Connect to camera service
	_, err = cameraworker.New(cfg.Addresses.Camera)
	if err != nil {
		loggerWriter.Error(errCameraService, logger.Args{"Error": err.Error()})
	}

	// Connect to video stream service
	_, err = videoworker.New(cfg.Addresses.Video)
	if err != nil {
		loggerWriter.Error(errVideoStreamService, logger.Args{"Error": err.Error()})
	}

	// Connect to user service
	_, err = deviceworker.New(cfg.Addresses.User)
	if err != nil {
		loggerWriter.Error(errDeviceService, logger.Args{"Error": err.Error()})
	}

	authUseCase := auth.New(loggerWriter, userWorker)
	useCase := usecase.New(authUseCase)

	server := new(server.Server)
	if err := server.Run(cfg.SHWSAddress, routes.SetRoutes(loggerWriter, useCase)); err != nil {
		log.Fatal(err)
	}
}
