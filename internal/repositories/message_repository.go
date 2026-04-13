package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/suck-seed/yapp/internal/database"
	dto "github.com/suck-seed/yapp/internal/dto/message"
	"github.com/suck-seed/yapp/internal/models"
	"github.com/suck-seed/yapp/internal/utils"
)

type IMessageRepository interface {
	// Message creation flow
	CreateMessage(ctx context.Context, db database.DBRunner, message *models.Message) (*models.Message, error)
	AddMessageMention(ctx context.Context, db database.DBRunner, messageID uuid.UUID, userID uuid.UUID) error
	AddAttachment(ctx context.Context, db database.DBRunner, attachment *models.Attachment) (*models.Attachment, error)

	// Read
	GetMessageByID(ctx context.Context, db database.DBRunner, messageID uuid.UUID) (*models.Message, error)
	GetMessageDetailed(ctx context.Context, db database.DBRunner, messageID uuid.UUID) (*dto.MessageDetailed, error)
	GetMessagesByRoomID(ctx context.Context, db database.DBRunner, roomID uuid.UUID, limit int, offset int) ([]*models.Message, error)
	GetMessages(ctx context.Context, db database.DBRunner, params *dto.MessageQueryParams) ([]*dto.MessageDetailed, error)

	// Write
	UpdateMessageContent(ctx context.Context, db database.DBRunner, messageID uuid.UUID, content string) (*models.Message, error)
	SoftDeleteMessage(ctx context.Context, db database.DBRunner, messageID uuid.UUID) error
	UpdateMessage(ctx context.Context, db database.DBRunner, message *models.Message) (*models.Message, error)
	DeleteMessage(ctx context.Context, db database.DBRunner, message *models.Message) error

	// Reactions
	AddReaction(ctx context.Context, db database.DBRunner, reactionID uuid.UUID, messageID uuid.UUID, userID uuid.UUID, emoji string) error
	RemoveReaction(ctx context.Context, db database.DBRunner, messageID uuid.UUID, userID uuid.UUID, emoji string) (bool, error)
}

type messageRepository struct{}

func NewMessageRepository() IMessageRepository {
	return &messageRepository{}
}

// ── CreateMessage ─────────────────────────────────────────────────────────────

func (r *messageRepository) CreateMessage(ctx context.Context, db database.DBRunner, message *models.Message) (*models.Message, error) {
	query := `
		INSERT INTO messages (id, room_id, author_id, content, sent_at, mention_everyone)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, room_id, author_id, content, sent_at, edited_at, deleted_at, mention_everyone, created_at, updated_at
	`
	out := &models.Message{}
	err := db.QueryRow(ctx, query,
		message.ID, message.RoomID, message.AuthorID,
		message.Content, message.SentAt, message.MentionEveryone,
	).Scan(
		&out.ID, &out.RoomID, &out.AuthorID, &out.Content, &out.SentAt,
		&out.EditedAt, &out.DeletedAt, &out.MentionEveryone,
		&out.CreatedAt, &out.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ── AddAttachment ─────────────────────────────────────────────────────────────

func (r *messageRepository) AddAttachment(ctx context.Context, db database.DBRunner, attachment *models.Attachment) (*models.Attachment, error) {
	query := `
		INSERT INTO attachments (id, message_id, file_name, url, file_type, file_size)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, message_id, file_name, url, file_type, file_size, created_at, updated_at
	`
	out := &models.Attachment{}
	err := db.QueryRow(ctx, query,
		attachment.ID, attachment.MessageID, attachment.FileName,
		attachment.URL, attachment.FileType, attachment.FileSize,
	).Scan(
		&out.ID, &out.MessageID, &out.FileName, &out.URL,
		&out.FileType, &out.FileSize, &out.CreatedAt, &out.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ── AddMessageMention ─────────────────────────────────────────────────────────

func (r *messageRepository) AddMessageMention(ctx context.Context, db database.DBRunner, messageID uuid.UUID, userID uuid.UUID) error {
	_, err := db.Exec(ctx, `
		INSERT INTO message_mentions (message_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (message_id, user_id) DO NOTHING
	`, messageID, userID)
	return err
}

// ── GetMessageByID ────────────────────────────────────────────────────────────

func (r *messageRepository) GetMessageByID(ctx context.Context, db database.DBRunner, messageID uuid.UUID) (*models.Message, error) {
	query := `
		SELECT id, room_id, author_id, content, sent_at, edited_at, deleted_at, mention_everyone, created_at, updated_at
		FROM messages
		WHERE id = $1 AND deleted_at IS NULL
	`
	out := &models.Message{}
	err := db.QueryRow(ctx, query, messageID).Scan(
		&out.ID, &out.RoomID, &out.AuthorID, &out.Content, &out.SentAt,
		&out.EditedAt, &out.DeletedAt, &out.MentionEveryone,
		&out.CreatedAt, &out.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ── GetMessageDetailed ────────────────────────────────────────────────────────
// Returns a single message with author, attachments, and reactions joined in.

func (r *messageRepository) GetMessageDetailed(ctx context.Context, db database.DBRunner, messageID uuid.UUID) (*dto.MessageDetailed, error) {
	query := `
		SELECT
			m.id, m.room_id, m.author_id, m.content, m.mention_everyone,
			m.sent_at, m.edited_at, m.created_at, m.updated_at,

			u.id, u.username, u.email, u.avatar_url,

			a.id, a.message_id, a.url, a.file_name, a.file_type, a.created_at, a.updated_at,

			r.id, r.emoji, r.user_id, ru.username, ru.avatar_url,

			mu.id, mu.username, mu.email, mu.avatar_url
		FROM messages m
		INNER JOIN users u ON m.author_id = u.id
		LEFT JOIN attachments a ON m.id = a.message_id
		LEFT JOIN reactions r ON m.id = r.message_id
		LEFT JOIN users ru ON r.user_id = ru.id
		LEFT JOIN message_mentions mm ON m.id = mm.message_id
		LEFT JOIN users mu ON mm.user_id = mu.id
		WHERE m.id = $1 AND m.deleted_at IS NULL
	`
	rows, err := db.Query(ctx, query, messageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages, err := r.scanMessagesWithDetails(rows)
	if err != nil {
		return nil, err
	}
	if len(messages) == 0 {
		return nil, pgx.ErrNoRows
	}
	return messages[0], nil
}

// ── GetMessagesByRoomID ───────────────────────────────────────────────────────

func (r *messageRepository) GetMessagesByRoomID(ctx context.Context, db database.DBRunner, roomID uuid.UUID, limit int, offset int) ([]*models.Message, error) {
	query := `
		SELECT id, room_id, author_id, content, sent_at, edited_at, deleted_at, mention_everyone, created_at, updated_at
		FROM messages
		WHERE room_id = $1 AND deleted_at IS NULL
		ORDER BY sent_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := db.Query(ctx, query, roomID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		m := &models.Message{}
		if err := rows.Scan(
			&m.ID, &m.RoomID, &m.AuthorID, &m.Content, &m.SentAt,
			&m.EditedAt, &m.DeletedAt, &m.MentionEveryone,
			&m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, rows.Err()
}

// ── GetMessages ───────────────────────────────────────────────────────────────
func (r *messageRepository) GetMessages(ctx context.Context, db database.DBRunner, params *dto.MessageQueryParams) ([]*dto.MessageDetailed, error) {
	if params.Around != nil {
		return r.getMessagesAround(ctx, db, params)
	}

	query := `
		WITH target_messages AS (
			SELECT m.id, m.room_id, m.author_id, m.content, m.mention_everyone,
				   m.sent_at, m.edited_at, m.created_at, m.updated_at
			FROM messages m
			WHERE m.room_id = $1
			  AND m.deleted_at IS NULL
	`

	args := []any{params.RoomID}
	argIdx := 1

	if params.Before != nil {
		argIdx++
		query += fmt.Sprintf(`
			AND (
				m.sent_at < (SELECT sent_at FROM messages WHERE id = $%d)
				OR (
					m.sent_at = (SELECT sent_at FROM messages WHERE id = $%d)
					AND m.id < $%d
				)
			)
			ORDER BY m.sent_at DESC, m.id DESC
		`, argIdx, argIdx, argIdx)
		args = append(args, params.Before)
	} else if params.After != nil {
		argIdx++
		query += fmt.Sprintf(`
			AND (
				m.sent_at > (SELECT sent_at FROM messages WHERE id = $%d)
				OR (
					m.sent_at = (SELECT sent_at FROM messages WHERE id = $%d)
					AND m.id > $%d
				)
			)
			ORDER BY m.sent_at ASC, m.id ASC
		`, argIdx, argIdx, argIdx)
		args = append(args, params.After)
	} else {
		query += `ORDER BY m.sent_at DESC, m.id DESC`
	}

	argIdx++
	query += fmt.Sprintf(` LIMIT $%d`, argIdx)
	args = append(args, params.Limit)

	query += `
		)
		SELECT
			tm.id, tm.room_id, tm.author_id, tm.content, tm.mention_everyone,
			tm.sent_at, tm.edited_at, tm.created_at, tm.updated_at,

			u.id, u.username, u.email, u.avatar_url,

			a.id, a.message_id, a.url, a.file_name, a.file_type, a.created_at, a.updated_at,

			r.id, r.emoji, r.user_id, ru.username, ru.avatar_url,

			mu.id, mu.username, mu.email, mu.avatar_url
		FROM target_messages tm
		INNER JOIN users u ON tm.author_id = u.id
		LEFT JOIN attachments a ON tm.id = a.message_id
		LEFT JOIN reactions r ON tm.id = r.message_id
		LEFT JOIN users ru ON r.user_id = ru.id
		LEFT JOIN message_mentions mm ON tm.id = mm.message_id
		LEFT JOIN users mu ON mm.user_id = mu.id
		ORDER BY tm.sent_at ASC, tm.id ASC
	`

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, utils.ErrorFetchingMessages
	}
	defer rows.Close()

	return r.scanMessagesWithDetails(rows)
}

// ── getMessagesAround ─────────────────────────────────────────────────────────
func (r *messageRepository) getMessagesAround(ctx context.Context, db database.DBRunner, params *dto.MessageQueryParams) ([]*dto.MessageDetailed, error) {
	halfLimit := params.Limit / 2

	query := `
		WITH target_messages AS (
			(
				SELECT m.id, m.room_id, m.author_id, m.content, m.mention_everyone,
					   m.sent_at, m.edited_at, m.created_at, m.updated_at
				FROM messages m
				WHERE m.room_id = $1
				  AND m.deleted_at IS NULL
				  AND (
					  m.sent_at < (SELECT sent_at FROM messages WHERE id = $2)
					  OR (
						  m.sent_at = (SELECT sent_at FROM messages WHERE id = $2)
						  AND m.id <= $2
					  )
				  )
				ORDER BY m.sent_at DESC, m.id DESC
				LIMIT $3
			)
			UNION ALL
			(
				SELECT m.id, m.room_id, m.author_id, m.content, m.mention_everyone,
					   m.sent_at, m.edited_at, m.created_at, m.updated_at
				FROM messages m
				WHERE m.room_id = $1
				  AND m.deleted_at IS NULL
				  AND (
					  m.sent_at > (SELECT sent_at FROM messages WHERE id = $2)
					  OR (
						  m.sent_at = (SELECT sent_at FROM messages WHERE id = $2)
						  AND m.id > $2
					  )
				  )
				ORDER BY m.sent_at ASC, m.id ASC
				LIMIT $4
			)
		)
		SELECT
			tm.id, tm.room_id, tm.author_id, tm.content, tm.mention_everyone,
			tm.sent_at, tm.edited_at, tm.created_at, tm.updated_at,

			u.id, u.username, u.email, u.avatar_url,

			a.id, a.message_id, a.url, a.file_name, a.file_type, a.created_at, a.updated_at,

			r.id, r.emoji, r.user_id, ru.username, ru.avatar_url,

			mu.id, mu.username, mu.email, mu.avatar_url
		FROM target_messages tm
		INNER JOIN users u ON tm.author_id = u.id
		LEFT JOIN attachments a ON tm.id = a.message_id
		LEFT JOIN reactions r ON tm.id = r.message_id
		LEFT JOIN users ru ON r.user_id = ru.id
		LEFT JOIN message_mentions mm ON tm.id = mm.message_id
		LEFT JOIN users mu ON mm.user_id = mu.id
		ORDER BY tm.sent_at ASC, tm.id ASC
	`

	rows, err := db.Query(ctx, query, params.RoomID, params.Around, halfLimit, halfLimit)
	if err != nil {
		return nil, utils.ErrorFetchingMessages
	}
	defer rows.Close()

	return r.scanMessagesWithDetails(rows)
}

// ── UpdateMessageContent ──────────────────────────────────────────────────────

func (r *messageRepository) UpdateMessageContent(ctx context.Context, db database.DBRunner, messageID uuid.UUID, content string) (*models.Message, error) {
	query := `
		UPDATE messages
		SET content = $1, edited_at = now(), updated_at = now()
		WHERE id = $2 AND deleted_at IS NULL
		RETURNING id, room_id, author_id, content, sent_at, edited_at, deleted_at, mention_everyone, created_at, updated_at
	`
	out := &models.Message{}
	err := db.QueryRow(ctx, query, content, messageID).Scan(
		&out.ID, &out.RoomID, &out.AuthorID, &out.Content, &out.SentAt,
		&out.EditedAt, &out.DeletedAt, &out.MentionEveryone,
		&out.CreatedAt, &out.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ── SoftDeleteMessage ─────────────────────────────────────────────────────────

func (r *messageRepository) SoftDeleteMessage(ctx context.Context, db database.DBRunner, messageID uuid.UUID) error {
	_, err := db.Exec(ctx, `
		UPDATE messages SET deleted_at = now(), updated_at = now() WHERE id = $1
	`, messageID)
	return err
}

// ── AddReaction ───────────────────────────────────────────────────────────────
// ON CONFLICT DO NOTHING makes this idempotent — duplicate reactions are silently ignored.

func (r *messageRepository) AddReaction(ctx context.Context, db database.DBRunner, reactionID uuid.UUID, messageID uuid.UUID, userID uuid.UUID, emoji string) error {
	_, err := db.Exec(ctx, `
		INSERT INTO reactions (id, message_id, user_id, emoji)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (message_id, user_id, emoji) DO NOTHING
	`, reactionID, messageID, userID, emoji)
	return err
}

// ── RemoveReaction ────────────────────────────────────────────────────────────
// Returns (true, nil) if a row was deleted, (false, nil) if reaction didn't exist.

func (r *messageRepository) RemoveReaction(ctx context.Context, db database.DBRunner, messageID uuid.UUID, userID uuid.UUID, emoji string) (bool, error) {
	tag, err := db.Exec(ctx, `
		DELETE FROM reactions WHERE message_id = $1 AND user_id = $2 AND emoji = $3
	`, messageID, userID, emoji)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() == 1, nil
}

// ── scanMessagesWithDetails ───────────────────────────────────────────────────
// Collapses the JOIN-expanded rows (one row per attachment per reaction) back into
// MessageDetailed structs, deduplicating attachments and grouping reactions by emoji.
func (r *messageRepository) scanMessagesWithDetails(rows interface {
	Next() bool
	Err() error
	Scan(dest ...any) error
}) ([]*dto.MessageDetailed, error) {
	messageMap := make(map[uuid.UUID]*dto.MessageDetailed)
	messageOrder := []uuid.UUID{}

	for rows.Next() {
		var (
			message models.Message
			author  dto.UserBasic

			attachmentID    *uuid.UUID
			attachmentMsgID *uuid.UUID
			attURL          *string
			attFileName     *string
			attFileType     *string
			attCreatedAt    *time.Time
			attUpdatedAt    *time.Time

			reactionID       *uuid.UUID
			emoji            *string
			reactorUserID    *uuid.UUID
			reactorUsername  *string
			reactorAvatarURL *string

			mentionUserID    *uuid.UUID
			mentionUsername  *string
			mentionEmail     *string
			mentionAvatarURL *string
		)

		if err := rows.Scan(
			&message.ID, &message.RoomID, &message.AuthorID, &message.Content, &message.MentionEveryone,
			&message.SentAt, &message.EditedAt, &message.CreatedAt, &message.UpdatedAt,

			&author.ID, &author.Username, &author.Email, &author.AvatarURL,

			&attachmentID, &attachmentMsgID, &attURL, &attFileName, &attFileType, &attCreatedAt, &attUpdatedAt,

			&reactionID, &emoji, &reactorUserID, &reactorUsername, &reactorAvatarURL,

			&mentionUserID, &mentionUsername, &mentionEmail, &mentionAvatarURL,
		); err != nil {
			return nil, utils.ErrorFetchingMessages
		}

		msgDetailed, exists := messageMap[message.ID]
		if !exists {
			msgDetailed = &dto.MessageDetailed{
				Message:     message,
				Author:      author,
				Attachments: []dto.AttachmentResponseMinimal{},
				Reactions:   []dto.ReactionGroup{},
				Mentions:    []dto.UserBasic{},
			}
			messageMap[message.ID] = msgDetailed
			messageOrder = append(messageOrder, message.ID)
		}

		if attachmentID != nil {
			found := false
			for _, a := range msgDetailed.Attachments {
				if a.ID == *attachmentID {
					found = true
					break
				}
			}
			if !found {
				msgDetailed.Attachments = append(msgDetailed.Attachments, dto.AttachmentResponseMinimal{
					ID:        *attachmentID,
					MessageID: *attachmentMsgID,
					URL:       *attURL,
					FileName:  *attFileName,
					FileType:  attFileType,
					CreatedAt: *attCreatedAt,
					UpdatedAt: *attUpdatedAt,
				})
			}
		}

		if reactionID != nil && emoji != nil {
			emojiFound := false
			for i := range msgDetailed.Reactions {
				if msgDetailed.Reactions[i].Emoji != *emoji {
					continue
				}
				emojiFound = true

				reactorFound := false
				for _, reactor := range msgDetailed.Reactions[i].Reactors {
					if reactorUserID != nil && reactor.ID == *reactorUserID {
						reactorFound = true
						break
					}
				}

				if !reactorFound && reactorUserID != nil && reactorUsername != nil {
					msgDetailed.Reactions[i].Reactors = append(msgDetailed.Reactions[i].Reactors, dto.UserBasic{
						ID:        *reactorUserID,
						Username:  *reactorUsername,
						AvatarURL: reactorAvatarURL,
					})
					msgDetailed.Reactions[i].Count++
				}
				break
			}

			if !emojiFound {
				group := dto.ReactionGroup{
					Emoji:    *emoji,
					Count:    0,
					Reactors: []dto.UserBasic{},
				}

				if reactorUserID != nil && reactorUsername != nil {
					group.Reactors = append(group.Reactors, dto.UserBasic{
						ID:        *reactorUserID,
						Username:  *reactorUsername,
						AvatarURL: reactorAvatarURL,
					})
					group.Count = 1
				}

				msgDetailed.Reactions = append(msgDetailed.Reactions, group)
			}
		}

		if mentionUserID != nil && mentionUsername != nil && mentionEmail != nil {
			found := false
			for _, mentionedUser := range msgDetailed.Mentions {
				if mentionedUser.ID == *mentionUserID {
					found = true
					break
				}
			}

			if !found {
				msgDetailed.Mentions = append(msgDetailed.Mentions, dto.UserBasic{
					ID:        *mentionUserID,
					Username:  *mentionUsername,
					Email:     *mentionEmail,
					AvatarURL: mentionAvatarURL,
				})
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, utils.ErrorMessageRowsIteration
	}

	result := make([]*dto.MessageDetailed, 0, len(messageOrder))
	for _, id := range messageOrder {
		result = append(result, messageMap[id])
	}

	return result, nil
}

// ── TODO stubs ────────────────────────────────────────────────────────────────

func (r *messageRepository) UpdateMessage(ctx context.Context, db database.DBRunner, message *models.Message) (*models.Message, error) {
	return &models.Message{}, nil
}

func (r *messageRepository) DeleteMessage(ctx context.Context, db database.DBRunner, message *models.Message) error {
	return nil
}
