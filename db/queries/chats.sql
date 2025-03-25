-- Get chat by ID
-- name: GetChatByID :one
SELECT *
FROM chats
WHERE id = $1;

-- List all chats by account_id
-- name: ListChatsByAccountID :many
SELECT *
FROM chats
WHERE account_id = $1
ORDER BY created_at DESC;

-- Create a new chat
-- name: CreateChat :one
INSERT INTO chats (messages, account_id)
VALUES ($1, $2)
RETURNING *;

-- Update chat's messages
-- name: UpdateChatMessages :one
UPDATE chats
SET messages = $1
WHERE id = $2
RETURNING *;

-- Delete chat by ID
-- name: DeleteChat :exec
DELETE
FROM chats
WHERE id = $1;


-- name: AccountOwnsChat :one
SELECT EXISTS(SELECT 1
              FROM chats
              WHERE account_id = $1
                AND id = $2);


-- name: GetChatsByAccountID :many
SELECT *
FROM chats
WHERE account_id = $1;


-- name: AddMessageToChat :one
UPDATE chats
SET messages = messages || @newMessage::jsonb
WHERE id = @chatID
RETURNING *;


-- name: MarkAsReadByID :exec
UPDATE chats
SET unread_messages= false
WHERE id = $1;


-- name: MarkAsUnreadByID :exec
UPDATE chats
SET unread_messages = true
WHERE id = $1;


-- name: IsUnread :one
SELECT unread_messages
FROM chats
WHERE id = $1;

