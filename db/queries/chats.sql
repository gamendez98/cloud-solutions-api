-- Get chat by ID
-- name: GetChatByID :one
SELECT id, created_at, messages, account_id
FROM chats
WHERE id = $1;

-- List all chats by account_id
-- name: ListChatsByAccountID :many
SELECT id, created_at, messages, account_id
FROM chats
WHERE account_id = $1
ORDER BY created_at DESC;

-- Create a new chat
-- name: CreateChat :one
INSERT INTO chats (messages, account_id)
VALUES ($1, $2)
RETURNING id, created_at, messages, account_id;

-- Update chat's messages
-- name: UpdateChatMessages :one
UPDATE chats
SET messages = $1
WHERE id = $2
RETURNING id, created_at, messages, account_id;

-- Delete chat by ID
-- name: DeleteChat :exec
DELETE
FROM chats
WHERE id = $1;