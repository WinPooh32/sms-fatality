// Code generated by sqlc. DO NOT EDIT.
// source: version.sql

package db

import (
	"context"
)

const getVersion = `-- name: GetVersion :one
SELECT version_id, is_applied FROM goose_db_version
ORDER BY id DESC
LIMIT 1
`

type GetVersionRow struct {
	VersionID int64 `json:"version_id"`
	IsApplied bool  `json:"is_applied"`
}

func (q *Queries) GetVersion(ctx context.Context) (GetVersionRow, error) {
	row := q.db.QueryRowContext(ctx, getVersion)
	var i GetVersionRow
	err := row.Scan(&i.VersionID, &i.IsApplied)
	return i, err
}