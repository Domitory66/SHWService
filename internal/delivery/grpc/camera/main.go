package cameraservice

import (
	"google.golang.org/grpc"
)

type CameraWorker struct {
	ClientConnection *grpc.ClientConn
}

func New(uri string) (*CameraWorker, error) {
	conCameraService, err := grpc.NewClient(uri)
	if err != nil {
		return nil, err
	}
	return &CameraWorker{ClientConnection: conCameraService}, nil
}
