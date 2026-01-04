package ws

import (
	"context"

	dto "github.com/suck-seed/yapp/internal/dto/message"
	"github.com/suck-seed/yapp/internal/services"
)

type PersistFunction func(ctx context.Context, in *dto.InboundMessage) (*dto.OutboundMessage, error)

// MakePresistFunction : Performs various actions and pushes it to db
func MakePresistFunction(messageService services.IMessageService, userService services.IUserService) PersistFunction {
	return func(ctx context.Context, in *dto.InboundMessage) (*dto.OutboundMessage, error) {

		// send to messageService to handle
		saved, err := messageService.CreateMessage(context.Background(), &dto.CreateMessageReq{
			RoomID:          in.RoomID,
			AuthorID:        in.UserID,
			Content:         in.Content,
			SentAt:          in.SentAt,
			Attachments:     in.Attachments,
			MentionEveryone: in.MentionEveryone,
			Mentions:        in.Mentions,
		})

		if err != nil {
			return nil, err
		}

		return &dto.OutboundMessage{
			Type: in.Type,

			ID:       saved.ID,
			RoomID:   saved.RoomID,
			AuthorID: saved.AuthorID,
			Content:  saved.Content,
			SentAt:   saved.SentAt,

			CreatedAt: saved.CreatedAt,
			EditedAt:  saved.EditedAt,
			DeletedAt: saved.DeletedAt,
			UpdatedAt: saved.UpdatedAt,

			MentionsEveryone: saved.MentionsEveryone,
			Mentions:         saved.Mentions,
			Attachments:      saved.Attachments,

			Error: "",
		}, nil

	}
}
