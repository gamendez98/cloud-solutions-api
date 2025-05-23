-- Create a new document
-- name: CreateDocument :one
INSERT INTO documents (name, text, file_path, embedding, account_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, name, text, file_path, embedding, account_id;

-- Get a document by ID
-- name: GetDocumentByID :one
SELECT id, created_at, name, text, file_path, embedding, account_id
FROM documents
WHERE id = $1;

-- Delete a document by ID
-- name: DeleteDocument :exec
DELETE
FROM documents
WHERE id = $1;


-- Get all documents for a specific account
-- name: GetDocumentsByAccountID :many
SELECT id, created_at, name, text, file_path, embedding, account_id
FROM documents
WHERE account_id = $1
LIMIT $2 OFFSET $3;


-- name: AccountOwnsDocument :one
SELECT EXISTS(SELECT 1
              FROM documents
              WHERE account_id = $1
                AND id = $2);
