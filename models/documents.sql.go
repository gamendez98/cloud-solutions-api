// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: documents.sql

package models

import (
	"context"
	"database/sql"
)

const accountOwnsDocument = `-- name: AccountOwnsDocument :one
SELECT EXISTS(SELECT 1
              FROM documents
              WHERE account_id = $1
                AND id = $2)
`

type AccountOwnsDocumentParams struct {
	AccountID int32 `json:"accountId"`
	ID        int32 `json:"id"`
}

func (q *Queries) AccountOwnsDocument(ctx context.Context, arg AccountOwnsDocumentParams) (bool, error) {
	row := q.db.QueryRowContext(ctx, accountOwnsDocument, arg.AccountID, arg.ID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createDocument = `-- name: CreateDocument :one
INSERT INTO documents (name, text, file_path, embedding, account_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, name, text, file_path, embedding, account_id
`

type CreateDocumentParams struct {
	Name      string         `json:"name"`
	Text      sql.NullString `json:"text"`
	FilePath  sql.NullString `json:"filePath"`
	Embedding interface{}    `json:"embedding"`
	AccountID int32          `json:"accountId"`
}

// Create a new document
func (q *Queries) CreateDocument(ctx context.Context, arg CreateDocumentParams) (Document, error) {
	row := q.db.QueryRowContext(ctx, createDocument,
		arg.Name,
		arg.Text,
		arg.FilePath,
		arg.Embedding,
		arg.AccountID,
	)
	var i Document
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Name,
		&i.Text,
		&i.FilePath,
		&i.Embedding,
		&i.AccountID,
	)
	return i, err
}

const deleteDocument = `-- name: DeleteDocument :exec
DELETE
FROM documents
WHERE id = $1
`

// Delete a document by ID
func (q *Queries) DeleteDocument(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteDocument, id)
	return err
}

const getDocumentByID = `-- name: GetDocumentByID :one
SELECT id, created_at, name, text, file_path, embedding, account_id
FROM documents
WHERE id = $1
`

// Get a document by ID
func (q *Queries) GetDocumentByID(ctx context.Context, id int32) (Document, error) {
	row := q.db.QueryRowContext(ctx, getDocumentByID, id)
	var i Document
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Name,
		&i.Text,
		&i.FilePath,
		&i.Embedding,
		&i.AccountID,
	)
	return i, err
}

const getDocumentsByAccount = `-- name: GetDocumentsByAccount :many
SELECT id, created_at, name, text, file_path, embedding, account_id
FROM documents
WHERE account_id = $1
`

// Get all documents for a specific account
func (q *Queries) GetDocumentsByAccount(ctx context.Context, accountID int32) ([]Document, error) {
	rows, err := q.db.QueryContext(ctx, getDocumentsByAccount, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Document{}
	for rows.Next() {
		var i Document
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.Name,
			&i.Text,
			&i.FilePath,
			&i.Embedding,
			&i.AccountID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
