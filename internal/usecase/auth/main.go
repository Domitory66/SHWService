package auth

import (
	userservice "shwservice/internal/delivery/grpc/user"
	"shwservice/pkg/logger"
)

type authUseCase struct {
	Logger     logger.Logger
	UserWorker userservice.UserWorker
}

func New(l logger.Logger, uw userservice.UserWorker) authUseCase {
	return authUseCase{Logger: l, UserWorker: uw}
}
