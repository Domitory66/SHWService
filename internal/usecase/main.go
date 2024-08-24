package usecase

import "context"

type Auth interface {
	SignIn(c context.Context) error
	SignUp(c context.Context) error
	LogOut(c context.Context) error
}

type UseCase struct {
	Auth
}

func New(auth Auth) UseCase {
	return UseCase{Auth: auth}
}
