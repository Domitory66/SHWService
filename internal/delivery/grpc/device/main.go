package deviceservice

import "google.golang.org/grpc"

type DeviceWorker struct {
	ClientConnection *grpc.ClientConn
}

func New(uri string) (*DeviceWorker, error) {
	conVideoStreamService, err := grpc.NewClient(uri)
	if err != nil {
		return nil, err
	}
	return &DeviceWorker{ClientConnection: conVideoStreamService}, nil
}
