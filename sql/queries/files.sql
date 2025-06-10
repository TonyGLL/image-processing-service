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
    i.id AS image_id,
    i.url,
    io.id AS option_id,
    io.format,
    io.quality
FROM image_processing_schema.images i
INNER JOIN image_processing_schema.images_options io ON io.image_id = i.id
ORDER BY i.id OFFSET $1 LIMIT $2;

