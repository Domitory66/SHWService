package service

import (
	"SmartHomeWebCam/SHWService/web/api/auth"
	gen "SmartHomeWebCam/SHWService/web/api/camera"
	video "SmartHomeWebCam/SHWService/web/api/video"
)

type Service struct {
	gen.CameraWorkerClient
	video.VideoStreamClient
	auth.AuthClient
}

func New(camClient gen.CameraWorkerClient, stream video.VideoStreamClient, authClient auth.AuthClient) *Service {
	return &Service{CameraWorkerClient: camClient, VideoStreamClient: stream, AuthClient: authClient}
}
