package videoservice

import "google.golang.org/grpc"

type VideoStreamWorker struct {
	ClientConnection *grpc.ClientConn
}

func New(uri string) (*VideoStreamWorker, error) {
	conVideoStreamService, err := grpc.NewClient(uri)
	if err != nil {
		return nil, err
	}
	return &VideoStreamWorker{ClientConnection: conVideoStreamService}, nil
}
