package main

import (
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

	log.Log.Info("Start Analytics service")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	conAnalytics, err := grpc.DialContext(ctx, "localhost:6200", grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Log.Fatal("Analytics service not started ", err.Error())
	}
	conVS, err := grpc.DialContext(context.Background(), "localhost:6201", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Log.Fatal("Analytics service not ready for stream ", err.Error())
	}
	clientAnalytics := gen.NewCameraWorkerClient(conAnalytics)
	videoStream := video.NewVideoStreamClient(conVS)
	//TODO clientAuth
	//TODO clientSH
	services := service.New(clientAnalytics, videoStream)
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
