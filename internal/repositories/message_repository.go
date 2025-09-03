package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/models"
)

type IMessageRepository interface {
	CreateMessage(ctx context.Context, message *models.Message) (*models.Message, error)
	GetMessageByID(ctx context.Context, messageID uuid.UUID) (*models.Message, error)
	GetMessagesByRoomID(ctx context.Context, roomID uuid.UUID, limit int, offset int) ([]*models.Message, error)
	GetRoomMessages(ctx context.Context, roomID uuid.UUID, before *time.Time, limit int) ([]*models.Message, error)
	UpdateMessage(ctx context.Context, message *models.Message) (*models.Message, error)
	DeleteMessage(ctx context.Context, message *models.Message) error
}

type messageRepository struct {
	db PGXTX
}

func NewMessageReposiroty(db PGXTX) IMessageRepository {

	return &messageRepository{
		db: db,
	}
}

func (r *messageRepository) CreateMessage(ctx context.Context, message *models.Message) (*models.Message, error) {

	query := `
			INSERT INTO messages (id, room_id, author_id, content, sent_at, edited_at, deleted_at, mention_everyone, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)

			RETURNING id, room_id, author_id, content, sent_at, edited_at, deleted_at, mention_everyone, created_at, updated_at

			`

	row := r.db.QueryRow(
		ctx,
		query,
		message.MessageId,
		message.RoomId,
		message.AuthorId,
		message.Content,
		message.SentAt,
		message.EditedAt,
		message.DeletedAt,
		message.MentionEveryone,
		message.CreatedAt,
		message.UpdatedAt,
	)

	createdMessage := &models.Message{}

	err := row.Scan(
		&createdMessage.MessageId,
		&createdMessage.RoomId,
		&createdMessage.AuthorId,
		&createdMessage.Content,
		&createdMessage.SentAt,
		&createdMessage.EditedAt,
		&createdMessage.DeletedAt,
		&createdMessage.MentionEveryone,
		&createdMessage.CreatedAt,
		&createdMessage.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return createdMessage, nil

}

func (r *messageRepository) GetMessageByID(ctx context.Context, messageID uuid.UUID) (*models.Message, error) {

	query := `
		SELECT id, room_id, author_id, content, sent_at, edited_at, deleted_at, mention_everyone, created_at, updated_at
		FROM messages
		WHERE id = $1 AND deleted_at IS NULL
	`

	message := &models.Message{}

	err := r.db.QueryRow(ctx, query, messageID).Scan(
		&message.MessageId,
		&message.RoomId,
		&message.AuthorId,
		&message.Content,
		&message.SentAt,
		&message.EditedAt,
		&message.DeletedAt,
		&message.MentionEveryone,
		&message.CreatedAt,
		&message.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return message, nil
}

func (r *messageRepository) GetMessagesByRoomID(ctx context.Context, roomID uuid.UUID, limit int, offset int) ([]*models.Message, error) {

	query := `
		SELECT id, room_id, author_id, content, sent_at, edited_at, deleted_at, mention_everyone, created_at, updated_at
		FROM messages
		WHERE room_id = $1 AND deleted_at IS NULL
		ORDER BY sent_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, roomID, limit, offset)
	if err != nil {
		return nil, err
	}

	messages := []*models.Message{}
	message := &models.Message{}

	for rows.Next() {

		err := rows.Scan(
			&message.MessageId,
			&message.RoomId,
			&message.AuthorId,
			&message.Content,
			&message.SentAt,
			&message.EditedAt,
			&message.DeletedAt,
			&message.MentionEveryone,
			&message.CreatedAt,
			&message.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		messages = append(messages, message)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *messageRepository) GetRoomMessages(ctx context.Context, roomID uuid.UUID, before *time.Time, limit int) ([]*models.Message, error) {

	var query string
	var args []interface{}

	if before != nil {
		query = `
			SELECT id, room_id, author_id, content, sent_at, edited_at, deleted_at, mention_everyone, created_at, updated_at
			FROM messages
			WHERE room_id = $1 AND deleted_at IS NULL AND sent_at < $2
			ORDER BY sent_at DESC
			LIMIT $3
		`

		args = []interface{}{roomID, *before, limit}

	} else {
		query = `
			SELECT id, room_id, author_id, content, sent_at, edited_at, deleted_at, mention_everyone, created_at, updated_at
			FROM messages
			WHERE room_id = $1 AND deleted_at IS NULL
			ORDER BY sent_at DESC
		`

		args = []interface{}{roomID, limit}

	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	messages := []*models.Message{}
	message := &models.Message{}

	for rows.Next() {

		err := rows.Scan(
			&message.MessageId,
			&message.RoomId,
			&message.AuthorId,
			&message.Content,
			&message.SentAt,
			&message.EditedAt,
			&message.DeletedAt,
			&message.MentionEveryone,
			&message.CreatedAt,
			&message.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		messages = append(messages, message)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil

}

// TODO : Implement message  update and delete
func (r *messageRepository) UpdateMessage(ctx context.Context, message *models.Message) (*models.Message, error) {

	return &models.Message{}, nil

}

func (r *messageRepository) DeleteMessage(ctx context.Context, message *models.Message) error {

	return nil
}
