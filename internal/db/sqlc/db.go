// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package db

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.createImageStmt, err = db.PrepareContext(ctx, createImage); err != nil {
		return nil, fmt.Errorf("error preparing query CreateImage: %w", err)
	}
	if q.createImageOptionsStmt, err = db.PrepareContext(ctx, createImageOptions); err != nil {
		return nil, fmt.Errorf("error preparing query CreateImageOptions: %w", err)
	}
	if q.createPasswordStmt, err = db.PrepareContext(ctx, createPassword); err != nil {
		return nil, fmt.Errorf("error preparing query CreatePassword: %w", err)
	}
	if q.createUserStmt, err = db.PrepareContext(ctx, createUser); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUser: %w", err)
	}
	if q.getAllImagesStmt, err = db.PrepareContext(ctx, getAllImages); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllImages: %w", err)
	}
	if q.getImageByIdStmt, err = db.PrepareContext(ctx, getImageById); err != nil {
		return nil, fmt.Errorf("error preparing query GetImageById: %w", err)
	}
	if q.getUserByUsernameStmt, err = db.PrepareContext(ctx, getUserByUsername); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserByUsername: %w", err)
	}
	if q.getUserPasswordStmt, err = db.PrepareContext(ctx, getUserPassword); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserPassword: %w", err)
	}
	if q.updateImageCropOptionsStmt, err = db.PrepareContext(ctx, updateImageCropOptions); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateImageCropOptions: %w", err)
	}
	if q.updateImageResizeOptionsStmt, err = db.PrepareContext(ctx, updateImageResizeOptions); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateImageResizeOptions: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.createImageStmt != nil {
		if cerr := q.createImageStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createImageStmt: %w", cerr)
		}
	}
	if q.createImageOptionsStmt != nil {
		if cerr := q.createImageOptionsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createImageOptionsStmt: %w", cerr)
		}
	}
	if q.createPasswordStmt != nil {
		if cerr := q.createPasswordStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createPasswordStmt: %w", cerr)
		}
	}
	if q.createUserStmt != nil {
		if cerr := q.createUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createUserStmt: %w", cerr)
		}
	}
	if q.getAllImagesStmt != nil {
		if cerr := q.getAllImagesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllImagesStmt: %w", cerr)
		}
	}
	if q.getImageByIdStmt != nil {
		if cerr := q.getImageByIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getImageByIdStmt: %w", cerr)
		}
	}
	if q.getUserByUsernameStmt != nil {
		if cerr := q.getUserByUsernameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserByUsernameStmt: %w", cerr)
		}
	}
	if q.getUserPasswordStmt != nil {
		if cerr := q.getUserPasswordStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserPasswordStmt: %w", cerr)
		}
	}
	if q.updateImageCropOptionsStmt != nil {
		if cerr := q.updateImageCropOptionsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateImageCropOptionsStmt: %w", cerr)
		}
	}
	if q.updateImageResizeOptionsStmt != nil {
		if cerr := q.updateImageResizeOptionsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateImageResizeOptionsStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                           DBTX
	tx                           *sql.Tx
	createImageStmt              *sql.Stmt
	createImageOptionsStmt       *sql.Stmt
	createPasswordStmt           *sql.Stmt
	createUserStmt               *sql.Stmt
	getAllImagesStmt             *sql.Stmt
	getImageByIdStmt             *sql.Stmt
	getUserByUsernameStmt        *sql.Stmt
	getUserPasswordStmt          *sql.Stmt
	updateImageCropOptionsStmt   *sql.Stmt
	updateImageResizeOptionsStmt *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                           tx,
		tx:                           tx,
		createImageStmt:              q.createImageStmt,
		createImageOptionsStmt:       q.createImageOptionsStmt,
		createPasswordStmt:           q.createPasswordStmt,
		createUserStmt:               q.createUserStmt,
		getAllImagesStmt:             q.getAllImagesStmt,
		getImageByIdStmt:             q.getImageByIdStmt,
		getUserByUsernameStmt:        q.getUserByUsernameStmt,
		getUserPasswordStmt:          q.getUserPasswordStmt,
		updateImageCropOptionsStmt:   q.updateImageCropOptionsStmt,
		updateImageResizeOptionsStmt: q.updateImageResizeOptionsStmt,
	}
}
