package db

import "context"

type UsersQuerier interface {
	CreateUser(ctx context.Context, username string) (int32, error)
	CreatePassword(ctx context.Context, arg CreatePasswordParams) error
	GetUserByUsername(ctx context.Context, username string) (GetUserByUsernameRow, error)
}

var _ UsersQuerier = (*Queries)(nil)
