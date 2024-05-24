package main

import (
	"SmartHomeWebCam/SHWService/web/api/auth"
	gen "SmartHomeWebCam/SHWService/web/api/camera"
	video "SmartHomeWebCam/SHWService/web/api/video"
	"SmartHomeWebCam/SHWService/web/app"
	"SmartHomeWebCam/SHWService/web/app/handler"
	"SmartHomeWebCam/SHWService/web/app/logger"
	"SmartHomeWebCam/SHWService/web/service"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log := logger.NewLogger("./logs/PublicLog/", "network.log")

	if err := initConfig(); err != nil {
		log.Log.Fatalf("Error to initialize configs: %s", err.Error())
	}

	log.Log.Info("Start initialization services")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	conAnalytics, err := grpc.DialContext(ctx, viper.GetString("analyticsServiceAddr"), grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Log.Fatal("Analytics service not started ", err.Error())
	}
	conVS, err := grpc.DialContext(context.Background(), viper.GetString("analyticsVideoStream"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Log.Fatal("Analytics service not ready for stream ", err.Error())
	}
	conAuth, err := grpc.DialContext(context.Background(), viper.GetString("authServiceAddr"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Log.Fatal("Auth service not started ", err.Error())
	}
	clientAnalytics := gen.NewCameraWorkerClient(conAnalytics)
	videoStream := video.NewVideoStreamClient(conVS)
	authClient := auth.NewAuthClient(conAuth)
	services := service.New(clientAnalytics, videoStream, authClient)
	handlers := handler.New(services)

	srv := new(app.Server)
	go func() {
		if err := srv.Run(viper.GetString("publicAddr"), handlers.SetupRoutes(log)); err != nil {
			log.Log.Fatal("ERROR occurred while running http public Server: ", err.Error())
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Log.Println("Shutdown Server ... ")
}

func initConfig() error {
	viper.SetConfigName("config")
	viper.AddConfigPath("configs")

	return viper.ReadInConfig()
}
