package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/dto"
	"github.com/suck-seed/yapp/internal/models"
	"github.com/suck-seed/yapp/internal/utils"
)

type IMessageRepository interface {

	//	Message creation flow
	CreateMessage(ctx context.Context, message *models.Message) (*models.Message, error)
	AddMessageMention(ctx context.Context, messageId uuid.UUID, userID uuid.UUID) error
	AddAttachment(ctx context.Context, attachment *models.Attachment) (*models.Attachment, error)

	//	Additional
	GetMessageByID(ctx context.Context, messageID uuid.UUID) (*models.Message, error)
	GetMessagesByRoomID(ctx context.Context, roomID uuid.UUID, limit int, offset int) ([]*models.Message, error)

	GetMessages(ctx context.Context, params *dto.MessageQueryParams) ([]*dto.MessageDetailed, error)

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
			INSERT INTO messages (id, room_id, author_id, content, sent_at, mention_everyone)
			VALUES ($1, $2, $3, $4, $5, $6)

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
		message.MentionEveryone,
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
		attachment.ID,
		attachment.MessageID,
		attachment.FileName,
		attachment.URL,
		attachment.FileType,
		attachment.FileSize,
	)

	attachmentCRES := &models.Attachment{}
	err := row.Scan(
		&attachmentCRES.ID,
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

	for rows.Next() {

		messageCRES := &models.Message{}

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

func (r *messageRepository) GetMessages(ctx context.Context, params *dto.MessageQueryParams) ([]*dto.MessageDetailed, error) {

	if params.Around != nil {
		return r.getMessagesAround(ctx, params)
	}
	query := `

	WITH target_messages AS (

	    SELECT m.id, m.room_id, m.author_id, m.content, m.mention_everyone,
					m.sent_at, m.edited_at, m.created_at, m.updated_at
		FROM messages m
		WHERE m.room_id = $1
		AND m.deleted_at IS NULL
		-- continued
	`

	args := []interface{}{params.RoomID}
	argCount := 1

	// CURSOR BASED PAGINATION
	if params.Before != nil {
		argCount++

		//  FETCHING ANYTHING BEFORE OR == params.Before
		//  AND whose uuid < params.Before
		// UUID7 are time based type shit
		query += fmt.Sprintf(`

		AND (
		    m.sent_at < (SELECT sent_at FROM messages WHERE id = $%d)
            OR
            m.sent_at = (SELECT sent_at FROM messages WHERE id = $%d)
			AND
			m.id < $%d
		)

		`, argCount, argCount, argCount)

		args = append(args, params.Before)

		// ORDERING
		query += `ORDER BY m.sent_at DESC, m.id DESC`
	} else if params.After != nil {
		argCount++

		query += fmt.Sprintf(`

		AND (
		    m.sent_at > (SELECT sent_at FROM messages WHERE id = $%d)
            OR
            m.sent_at = (SELECT sent_at FROM messages WHERE id = $%d)
			AND
			m.id > $%d
		)

		`, argCount, argCount, argCount)

		args = append(args, params.After)

		// SORTING BASED ON TIME AND UUID
		query += `ORDER BY m.sent_at ASC, m.id ASC`

	} else {

		// JUST SORT THE RECENT ONES
		query += `ORDER BY m.sent_at DESC, m.id DESC`

	}

	// add LIMIT
	argCount++
	query += fmt.Sprintf(` LIMIT $%d`, argCount)
	args = append(args, params.Limit)

	// Fetch and add Join for other data too

	// tm for Message
	// u for userbasic

	query += `
	 )
		SELECT
		tm.id, tm.room_id, tm.author_id, tm.content, tm.mention_everyone,
		tm.sent_at, tm.edited_at, tm.created_at, tm.updated_at

		u.id, u.username, u.email, u.avatar_url

		a.id, a.message_id, a.url, a.file_name, a.file_type, a.created_at, a.updated_at

		r.id, r.emoji, r.user_id

		FROM target_messages as tm
		INNER JOIN users u ON tm.author_id=u.id
		LEFT JOIN attachments a ON tm.id = a.message_id
		LEFT JOIN reactions r ON tm.id = r.message_id

		ORDER BY tm.sent_at ASC, tm.id ASC
	`

	// PROCESSING WITH DB
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, utils.ErrorFetchingMessages
	}

	defer rows.Close()

	return r.scanMessagesWithDetails(rows)

	// rows.Next()
	// rows.Err()
	// rows.Scan()

}

// GET MESSAGES AROUND A MESSAGE
func (r *messageRepository) getMessagesAround(ctx context.Context, params *dto.MessageQueryParams) ([]*dto.MessageDetailed, error) {

	return []*dto.MessageDetailed{}, nil

}

// accept the needed methods through the row interface
func (r *messageRepository) scanMessagesWithDetails(rows interface {
	Next() bool
	Err() error
	Scan(dest ...interface{}) error
}) ([]*dto.MessageDetailed, error) {
	messageMap := make(map[uuid.UUID]*dto.MessageDetailed)
	messageOrder := []uuid.UUID{}

	// iterate over the rows
	for rows.Next() {

		// define variables
		var (
			message    models.Message
			author     dto.UserBasic
			attachment *dto.AttachmentResponseMinimal
			mention    *dto.UserBasic
		)

	}
	return nil, nil

}

// TODO : Implement message  update and delete
func (r *messageRepository) UpdateMessage(ctx context.Context, message *models.Message) (*models.Message, error) {

	return &models.Message{}, nil

}

func (r *messageRepository) DeleteMessage(ctx context.Context, message *models.Message) error {

	return nil
}
