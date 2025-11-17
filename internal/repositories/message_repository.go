package repositories

import (
	"context"
	"fmt"
	"time"

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

	args := []any{params.RoomID}
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

	// LEFT JOIN is used for scenarios where a record in the left table can correspond to zero, one, or multiple (0...N) records in the right table
	// attachments a (0..N per message)
	// reactions r (0..N per message)

	query += `
	 )
		SELECT
		tm.id, tm.room_id, tm.author_id, tm.content, tm.mention_everyone,
		tm.sent_at, tm.edited_at, tm.created_at, tm.updated_at

		u.id, u.username, u.email, u.avatar_url

		a.id, a.message_id, a.url, a.file_name, a.file_type, a.created_at, a.updated_at

		r.id, r.emoji, r.user_id, ru.username, ru.avatar_url

		FROM target_messages as tm
		INNER JOIN users u ON tm.author_id=u.id
		LEFT JOIN attachments a ON tm.id = a.message_id
		LEFT JOIN reactions r ON tm.id = r.message_id
		LEFT JOIN users ru ON r.user_id = ru.user_id

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

	// before and after equally lina paryo so
	halfLimit := params.Limit / 2

	// query
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
					AND m.id <= 2
				)
			)
		ORDER BY m.sent_at DECS, m.id DESC
		LIMIT $3
		)

		UNION ALL (
		SELECT m.id, m.room_id, m.author_id, m.content, m.mention_everyone,
		m.sent_at, m.edited_at, m.created_at, m.updated_at
		FROM messages m
		WHERE m.room_id = $1
		    AND m.deleted_at IS NULL
			AND (
			    m.sent_at < (SELECT sent_at FROM messages WHERE id = $2)
				OR (
				m.sent_at = (SELECT sent_at FROM messages WHERE id = $2)
				AND m.id <= 2
				)
			)
			ORDER BY m.sent_at ASC, m.id ASC
			LIMIT $4
		)
	)

	SELECT
	tm.id, tm.room_id, tm.author_id, tm.content, tm.mention_everyone,
	tm.sent_at, tm.edited_at, tm.created_at, tm.updated_at

	u.id, u.username, u.email, u.avatar_url

	a.id, a.message_id, a.url, a.file_name, a.file_type, a.created_at, a.updated_at

	r.id, r.emoji, r.user_id, ru.username, ru.avatar_url

	FROM target_messages as tm
	INNER JOIN users u ON tm.author_id=u.id
	LEFT JOIN attachments a ON tm.id = a.message_id
	LEFT JOIN reactions r ON tm.id = r.message_id
	LEFT JOIN users ru ON r.user_id = ru.user_id

	ORDER BY tm.sent_at ASC, tm.id ASC
	`

	rows, err := r.db.Query(ctx, query, params.RoomID, params.Around, halfLimit, halfLimit)
	if err != nil {
		return nil, utils.ErrorFetchingMessages
	}

	defer rows.Close()

	return r.scanMessagesWithDetails(rows)

}

// accept the needed methods through the row interface
func (r *messageRepository) scanMessagesWithDetails(rows interface {
	Next() bool
	Err() error
	Scan(dest ...any) error
}) ([]*dto.MessageDetailed, error) {
	messageMap := make(map[uuid.UUID]*dto.MessageDetailed)
	messageOrder := []uuid.UUID{}

	// iterate over the rows
	for rows.Next() {

		// define variables
		var (
			message models.Message
			author  dto.UserBasic

			attachmentID    *uuid.UUID
			attachmentMsgID *uuid.UUID
			url             *string
			fileName        *string
			fileType        *string
			attCreatedAt    *time.Time
			attUpdatedAt    *time.Time

			// REACTION (COUNT NEEDED SO DOING IT LIKE THIS)
			reactionID       *uuid.UUID
			emoji            *string
			reactorUserID    *uuid.UUID
			reactorUsername  *string
			reactorAvatarURL *string
		)

		// gotta be in the same order
		err := rows.Scan(
			// tm
			&message.ID,
			&message.RoomId,
			&message.AuthorId,
			&message.Content,
			&message.MentionEveryone,
			&message.SentAt,
			&message.EditedAt,
			&message.CreatedAt,
			&message.UpdatedAt,

			// author
			&author.ID,
			&author.Username,
			&author.Email,
			&author.AvatarURL,

			// attachment
			&attachmentID,
			&attachmentMsgID,
			&url,
			&fileName,
			&fileType,
			&attCreatedAt,
			&attUpdatedAt,

			// reaction
			&reactionID,
			&emoji,
			&reactorUserID,
			&reactorUsername,
			&reactorAvatarURL,
		)

		if err != nil {
			return nil, utils.ErrorFetchingMessages
		}

		// Check if message already exist in map or not
		msgDetailed, ok := messageMap[message.ID]

		// if doesnt exist
		if !ok {
			msgDetailed = &dto.MessageDetailed{
				Message: message,
				Author:  author,

				// we will handles these down below
				Attachments: []dto.AttachmentResponseMinimal{},
				Reactions:   []dto.ReactionGroup{},
				Mentions:    []dto.UserBasic{},
			}

			// add above messageDetailed on messagemap
			messageMap[message.ID] = msgDetailed

			// we order them serially on the way they came from db (ASC/DES)
			messageOrder = append(messageOrder, message.ID)
		}

		// ATTACHMENTS if exists and not added already
		if attachmentID != nil {

			// check if already exists
			found := false
			for _, attachment := range msgDetailed.Attachments {

				if attachment.ID == *attachmentID {
					found = true
					break
				}
			}

			if !found {
				attachment := &dto.AttachmentResponseMinimal{
					ID:        *attachmentID,
					MessageID: *attachmentMsgID,
					URL:       *url,
					FileName:  *fileName,
					FileType:  fileType,
					CreatedAt: *attCreatedAt,
					UpdatedAt: *attCreatedAt,
				}

				// add the attachment to the Attachments in msgDetailed
				msgDetailed.Attachments = append(msgDetailed.Attachments, *attachment)
			}

		}

		// REACTION IF EXISTS
		if reactionID != nil && emoji != nil {
			found := false
			for i := range msgDetailed.Reactions {

				// is this emoji already added
				if msgDetailed.Reactions[i].Emoji == *emoji {
					found = true

					// is this current userAdded in the userlist?
					userFound := false
					// loop through the reactors
					for _, reactor := range msgDetailed.Reactions[i].Reactors {
						// loop over reactors
						if reactorUserID != nil && reactor.ID == *reactorUserID {
							userFound = true
							break
						}
					}

					if userFound && reactorUserID != nil {

						// add user
						currentReactor := &dto.UserBasic{
							ID:        *reactorUserID,
							Username:  *reactorUsername,
							AvatarURL: reactorAvatarURL,
						}
						msgDetailed.Reactions[i].Reactors = append(msgDetailed.Reactions[i].Reactors, *currentReactor)

					}
					break
				}

			}

			// NOT FOUND, NEW EMOJI
			if !found {
				reactionGroup := &dto.ReactionGroup{
					Emoji:    *emoji,
					Count:    1,
					Reactors: []dto.UserBasic{},
				}

				// add user
				if reactorUserID != nil {
					currentReactor := &dto.UserBasic{
						ID:        *reactorUserID,
						Username:  *reactorUsername,
						AvatarURL: reactorAvatarURL,
					}
					reactionGroup.Reactors = append(reactionGroup.Reactors, *currentReactor)

				}

				// add this emoji group to Reaction in msgDetailed
				msgDetailed.Reactions = append(msgDetailed.Reactions, *reactionGroup)
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, utils.ErrorMessageRowsIteration
	}

	// convert map to ordered slice
	result := make([]*dto.MessageDetailed, 0, len(messageOrder))
	for _, msgID := range messageOrder {
		result = append(result, messageMap[msgID])
	}

	return result, nil

}

// TODO : Implement message  update and delete
func (r *messageRepository) UpdateMessage(ctx context.Context, message *models.Message) (*models.Message, error) {

	return &models.Message{}, nil

}

func (r *messageRepository) DeleteMessage(ctx context.Context, message *models.Message) error {

	return nil
}
