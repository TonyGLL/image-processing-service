-- name: CreateImage :one
INSERT INTO image_processing_schema.images (url) VALUES($1) RETURNING id;

-- name: CreateImageOptions :exec
INSERT INTO 
image_processing_schema.images_options (
    resize_width,
    resize_height,
    crop_width,
    crop_height,
    crop_x,
    crop_y,
    rotate,
    format,
    grayscale,
    sepia,
    image_id
) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);

-- name: GetAllImages :many
SELECT 
    i.id,
    i.url,
    jsonb_build_object(
        'id', io.id,
        'resize_width', io.resize_width,
        'resize_height', io.resize_height,
        'crop_width', io.crop_width,
        'crop_height', io.crop_height,
        'crop_x', io.crop_x,
        'crop_y', io.crop_y,
        'rotate', io.rotate,
        'format', io.format,
        'grayscale', io.grayscale,
        'sepia', io.sepia
    ) AS transformations
FROM image_processing_schema.images i
INNER JOIN image_processing_schema.images_options io ON io.image_id = i.id
ORDER BY i.id OFFSET $1 LIMIT $2;

-- name: GetImageById :one
SELECT 
    i.id,
    i.url,
    jsonb_build_object(
        'id', io.id,
        'resize_width', io.resize_width,
        'resize_height', io.resize_height,
        'crop_width', io.crop_width,
        'crop_height', io.crop_height,
        'crop_x', io.crop_x,
        'crop_y', io.crop_y,
        'rotate', io.rotate,
        'format', io.format,
        'grayscale', io.grayscale,
        'sepia', io.sepia
    ) AS transformations
FROM image_processing_schema.images i
INNER JOIN image_processing_schema.images_options io ON io.image_id = i.id
WHERE i.id = $1
LIMIT 1;

-- name: UpdateImageResizeOptions :exec
UPDATE image_processing_schema.images_options
SET resize_width = $1,
    resize_height = $2,
    updated_at = now()
WHERE image_id = $3;

-- name: UpdateImageCropOptions :exec
UPDATE image_processing_schema.images_options
SET crop_width = $1,
    crop_height = $2,
    crop_x = $3,
    crop_y = $4,
    updated_at = now()
WHERE image_id = $5;