package db

import "context"

type FilesQuerier interface {
	CreateImage(ctx context.Context, url string) (int32, error)
	CreateImageOptions(ctx context.Context, arg CreateImageOptionsParams) error
	GetAllImages(ctx context.Context, arg GetAllImagesParams) ([]GetAllImagesRow, error)
}

var _ FilesQuerier = (*Queries)(nil)
