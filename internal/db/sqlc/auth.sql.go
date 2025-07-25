// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: auth.sql

package db

import (
	"context"
)

const getUserPassword = `-- name: GetUserPassword :one
SELECT p.value FROM image_processing_schema.passwords p LEFT JOIN image_processing_schema.users u ON u.id = p.user_id WHERE u.username = $1
`

func (q *Queries) GetUserPassword(ctx context.Context, username string) (string, error) {
	row := q.queryRow(ctx, q.getUserPasswordStmt, getUserPassword, username)
	var value string
	err := row.Scan(&value)
	return value, err
}
