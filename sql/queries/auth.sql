-- name: GetUserPassword :one
SELECT p.value FROM image_processing_schema.passwords p LEFT JOIN image_processing_schema.users u ON u.id = p.user_id WHERE u.username = $1;