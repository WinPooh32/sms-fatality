// Code generated by sqlc. DO NOT EDIT.

package messages

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
	if q.createMessageStmt, err = db.PrepareContext(ctx, createMessage); err != nil {
		return nil, fmt.Errorf("error preparing query CreateMessage: %w", err)
	}
	if q.deleteMessageStmt, err = db.PrepareContext(ctx, deleteMessage); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteMessage: %w", err)
	}
	if q.getMessagesStmt, err = db.PrepareContext(ctx, getMessages); err != nil {
		return nil, fmt.Errorf("error preparing query GetMessages: %w", err)
	}
	if q.listPhonesStmt, err = db.PrepareContext(ctx, listPhones); err != nil {
		return nil, fmt.Errorf("error preparing query ListPhones: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.createMessageStmt != nil {
		if cerr := q.createMessageStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createMessageStmt: %w", cerr)
		}
	}
	if q.deleteMessageStmt != nil {
		if cerr := q.deleteMessageStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteMessageStmt: %w", cerr)
		}
	}
	if q.getMessagesStmt != nil {
		if cerr := q.getMessagesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getMessagesStmt: %w", cerr)
		}
	}
	if q.listPhonesStmt != nil {
		if cerr := q.listPhonesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listPhonesStmt: %w", cerr)
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
	db                DBTX
	tx                *sql.Tx
	createMessageStmt *sql.Stmt
	deleteMessageStmt *sql.Stmt
	getMessagesStmt   *sql.Stmt
	listPhonesStmt    *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                tx,
		tx:                tx,
		createMessageStmt: q.createMessageStmt,
		deleteMessageStmt: q.deleteMessageStmt,
		getMessagesStmt:   q.getMessagesStmt,
		listPhonesStmt:    q.listPhonesStmt,
	}
}
