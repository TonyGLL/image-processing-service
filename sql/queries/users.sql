-- name: CreateUser :one
INSERT INTO image_processing_schema.users (username) VALUES($1) RETURNING id;

-- name: CreatePassword :exec
INSERT INTO image_processing_schema.passwords (value, user_id) VALUES ($1, $2);

-- name: GetUserByUsername :one
SELECT id, username FROM image_processing_schema.users WHERE username = $1;