package userservice

import "google.golang.org/grpc"

type UserWorker struct {
	ClientConnection *grpc.ClientConn
}

func New(uri string) (UserWorker, error) {
	conUserService, err := grpc.NewClient(uri)
	if err != nil {
		return UserWorker{}, err
	}
	return UserWorker{ClientConnection: conUserService}, nil
}
