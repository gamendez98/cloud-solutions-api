// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: chats.sql

package models

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/sqlc-dev/pqtype"
)

const accountOwnsChat = `-- name: AccountOwnsChat :one
SELECT EXISTS(SELECT 1
              FROM chats
              WHERE account_id = $1
                AND id = $2)
`

type AccountOwnsChatParams struct {
	AccountID int32 `json:"accountId"`
	ID        int32 `json:"id"`
}

func (q *Queries) AccountOwnsChat(ctx context.Context, arg AccountOwnsChatParams) (bool, error) {
	row := q.db.QueryRowContext(ctx, accountOwnsChat, arg.AccountID, arg.ID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const addMessageToChat = `-- name: AddMessageToChat :one
UPDATE chats
SET messages = messages || $1::jsonb
WHERE id = $2
    RETURNING id, created_at, messages, account_id, unread_messages
`

type AddMessageToChatParams struct {
	Newmessage json.RawMessage `json:"newmessage"`
	Chatid     int32           `json:"chatid"`
}

func (q *Queries) AddMessageToChat(ctx context.Context, arg AddMessageToChatParams) (Chat, error) {
	row := q.db.QueryRowContext(ctx, addMessageToChat, arg.Newmessage, arg.Chatid)
	var i Chat
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Messages,
		&i.AccountID,
		&i.UnreadMessages,
	)
	return i, err
}

const createChat = `-- name: CreateChat :one
INSERT INTO chats (messages, account_id)
VALUES ($1, $2) RETURNING id, created_at, messages, account_id, unread_messages
`

type CreateChatParams struct {
	Messages  pqtype.NullRawMessage `json:"messages"`
	AccountID int32                 `json:"accountId"`
}

// Create a new chat
func (q *Queries) CreateChat(ctx context.Context, arg CreateChatParams) (Chat, error) {
	row := q.db.QueryRowContext(ctx, createChat, arg.Messages, arg.AccountID)
	var i Chat
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Messages,
		&i.AccountID,
		&i.UnreadMessages,
	)
	return i, err
}

const deleteChat = `-- name: DeleteChat :exec
DELETE
FROM chats
WHERE id = $1
`

// Delete chat by ID
func (q *Queries) DeleteChat(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteChat, id)
	return err
}

const getChatByID = `-- name: GetChatByID :one
SELECT id, created_at, messages, account_id, unread_messages
FROM chats
WHERE id = $1
`

// Get chat by ID
func (q *Queries) GetChatByID(ctx context.Context, id int32) (Chat, error) {
	row := q.db.QueryRowContext(ctx, getChatByID, id)
	var i Chat
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Messages,
		&i.AccountID,
		&i.UnreadMessages,
	)
	return i, err
}

const getChatsByAccountID = `-- name: GetChatsByAccountID :many
SELECT id, created_at, messages, account_id, unread_messages
FROM chats
WHERE account_id = $1
`

func (q *Queries) GetChatsByAccountID(ctx context.Context, accountID int32) ([]Chat, error) {
	rows, err := q.db.QueryContext(ctx, getChatsByAccountID, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Chat{}
	for rows.Next() {
		var i Chat
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.Messages,
			&i.AccountID,
			&i.UnreadMessages,
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

const isUnread = `-- name: IsUnread :one
SELECT unread_messages
FROM chats
WHERE id = $1
`

func (q *Queries) IsUnread(ctx context.Context, id int32) (sql.NullBool, error) {
	row := q.db.QueryRowContext(ctx, isUnread, id)
	var unread_messages sql.NullBool
	err := row.Scan(&unread_messages)
	return unread_messages, err
}

const listChatsByAccountID = `-- name: ListChatsByAccountID :many
SELECT id, created_at, messages, account_id, unread_messages
FROM chats
WHERE account_id = $1
ORDER BY created_at DESC
`

// List all chats by account_id
func (q *Queries) ListChatsByAccountID(ctx context.Context, accountID int32) ([]Chat, error) {
	rows, err := q.db.QueryContext(ctx, listChatsByAccountID, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Chat{}
	for rows.Next() {
		var i Chat
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.Messages,
			&i.AccountID,
			&i.UnreadMessages,
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

const markAsReadByID = `-- name: MarkAsReadByID :exec
UPDATE chats
SET unread_messages= false
WHERE id = $1
`

func (q *Queries) MarkAsReadByID(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, markAsReadByID, id)
	return err
}

const markAsUnreadByID = `-- name: MarkAsUnreadByID :exec
UPDATE chats
SET unread_messages = true
WHERE id = $1
`

func (q *Queries) MarkAsUnreadByID(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, markAsUnreadByID, id)
	return err
}

const updateChatMessages = `-- name: UpdateChatMessages :one
UPDATE chats
SET messages = $1
WHERE id = $2 RETURNING id, created_at, messages, account_id, unread_messages
`

type UpdateChatMessagesParams struct {
	Messages pqtype.NullRawMessage `json:"messages"`
	ID       int32                 `json:"id"`
}

// Update chat's messages
func (q *Queries) UpdateChatMessages(ctx context.Context, arg UpdateChatMessagesParams) (Chat, error) {
	row := q.db.QueryRowContext(ctx, updateChatMessages, arg.Messages, arg.ID)
	var i Chat
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Messages,
		&i.AccountID,
		&i.UnreadMessages,
	)
	return i, err
}
