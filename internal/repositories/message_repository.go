package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/dto"
	"github.com/suck-seed/yapp/internal/models"
)

type IMessageRepository interface {

	//	Message creation flow
	CreateMessage(ctx context.Context, message *models.Message) (*models.Message, error)
	AddMessageMention(ctx context.Context, messageId uuid.UUID, userID uuid.UUID) error
	AddAttachment(ctx context.Context, attachment *models.Attachment) (*models.Attachment, error)

	//	Additional
	GetMessageByID(ctx context.Context, messageID uuid.UUID) (*models.Message, error)
	GetMessagesByRoomID(ctx context.Context, roomID uuid.UUID, limit int, offset int) ([]*models.Message, error)

	GetMessages(ctx context.Context, queryParams *dto.MessageQueryParams) ([]*models.Message, error)

	UpdateMessage(ctx context.Context, message *models.Message) (*models.Message, error)
	DeleteMessage(ctx context.Context, message *models.Message) error
}

type messageRepository struct {
	db PGXTX
}

func NewMessageRepository(db PGXTX) IMessageRepository {

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
		message.ID,
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

	messageCRES := &models.Message{}

	err := row.Scan(
		&messageCRES.ID,
		&messageCRES.RoomId,
		&messageCRES.AuthorId,
		&messageCRES.Content,
		&messageCRES.SentAt,
		&messageCRES.EditedAt,
		&messageCRES.DeletedAt,
		&messageCRES.MentionEveryone,
		&messageCRES.CreatedAt,
		&messageCRES.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return messageCRES, nil

}

func (r *messageRepository) AddAttachment(ctx context.Context, attachment *models.Attachment) (*models.Attachment, error) {

	query := `
			INSERT INTO attachments (id, message_id, file_name, url, file_type, file_size)
			VALUES ($1, $2, $3, $4, $5, $6)

			RETURNING id, message_id, file_name, url, file_type, file_size, created_at, updated_at

			`

	row := r.db.QueryRow(
		ctx,
		query,
		attachment.AttachmentID,
		attachment.MessageID,
		attachment.FileName,
		attachment.URL,
		attachment.FileType,
		attachment.FileSize,
	)

	attachmentCRES := &models.Attachment{}
	err := row.Scan(
		&attachmentCRES.AttachmentID,
		&attachmentCRES.MessageID,
		&attachmentCRES.FileName,
		&attachmentCRES.URL,
		&attachmentCRES.FileType,
		&attachmentCRES.FileSize,
		&attachmentCRES.CreatedAt,
		&attachmentCRES.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return attachmentCRES, nil
}

func (r *messageRepository) AddMessageMention(ctx context.Context, messageId uuid.UUID, userID uuid.UUID) error {

	query := `
  				INSERT INTO message_mentions (message_id,user_id)
      			VALUES ($1, $2)
        		ON CONFLICT (message_id, user_id) DO NOTHING

   			`

	_, err := r.db.Exec(ctx, query, messageId, userID)

	return err
}

func (r *messageRepository) GetMessageByID(ctx context.Context, messageID uuid.UUID) (*models.Message, error) {

	query := `
		SELECT id, room_id, author_id, content, sent_at, edited_at, deleted_at, mention_everyone, created_at, updated_at
		FROM messages
		WHERE id = $1 AND deleted_at IS NULL
	`

	messageCRES := &models.Message{}

	err := r.db.QueryRow(ctx, query, messageID).Scan(
		&messageCRES.ID,
		&messageCRES.RoomId,
		&messageCRES.AuthorId,
		&messageCRES.Content,
		&messageCRES.SentAt,
		&messageCRES.EditedAt,
		&messageCRES.DeletedAt,
		&messageCRES.MentionEveryone,
		&messageCRES.CreatedAt,
		&messageCRES.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return messageCRES, nil
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

	messagesCRES := []*models.Message{}
	messageCRES := &models.Message{}

	for rows.Next() {

		err := rows.Scan(
			&messageCRES.ID,
			&messageCRES.RoomId,
			&messageCRES.AuthorId,
			&messageCRES.Content,
			&messageCRES.SentAt,
			&messageCRES.EditedAt,
			&messageCRES.DeletedAt,
			&messageCRES.MentionEveryone,
			&messageCRES.CreatedAt,
			&messageCRES.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		messagesCRES = append(messagesCRES, messageCRES)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messagesCRES, nil
}

func (r *messageRepository) GetMessages(ctx context.Context, queryParams *dto.MessageQueryParams) ([]*models.Message, error) {

	// query := `
	//            SELECT
	//                m.id, m.room_id, m.author_id, m.content,
	//                m.sent_at, m.edited_at, m.deleted_at,
	//                m.mention_everyone, m.created_at, m.updated_at
	//            FROM messages m
	//            WHERE m.room_id = $1
	//                AND m.deleted_at IS NULL
	//        `

	return []*models.Message{}, nil

}

// TODO : Implement message  update and delete
func (r *messageRepository) UpdateMessage(ctx context.Context, message *models.Message) (*models.Message, error) {

	return &models.Message{}, nil

}

func (r *messageRepository) DeleteMessage(ctx context.Context, message *models.Message) error {

	return nil
}
