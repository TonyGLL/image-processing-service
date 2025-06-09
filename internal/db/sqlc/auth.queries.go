package db

import "context"

type AuthQuerier interface {
	GetUserPassword(ctx context.Context, username string) (string, error)
}

var _ AuthQuerier = (*Queries)(nil)
