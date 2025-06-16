package db

import "context"

type FilesQuerier interface {
	CreateImage(ctx context.Context, url string) (int32, error)
	CreateImageOptions(ctx context.Context, arg CreateImageOptionsParams) error
	GetAllImages(ctx context.Context, arg GetAllImagesParams) ([]GetAllImagesRow, error)
	GetImageById(ctx context.Context, id int32) (GetImageByIdRow, error)
	UpdateImageResizeOptions(ctx context.Context, arg UpdateImageResizeOptionsParams) error
	UpdateImageCropOptions(ctx context.Context, arg UpdateImageCropOptionsParams) error
}

var _ FilesQuerier = (*Queries)(nil)
