package ws

import (
	"context"

	"github.com/suck-seed/yapp/internal/dto"
	"github.com/suck-seed/yapp/internal/services"
)

type PersistFunction func(ctx context.Context, in *InboundMessage) (*OutboundMessage, error)

// MakePresistFunc : Performs various actions and pushes it to db
func MakePresistFunc(messageService services.IMessageService, userService services.IUserService) PersistFunction {
	return func(ctx context.Context, in *InboundMessage) (*OutboundMessage, error) {

		// send to messageService to handle
		saved, err := messageService.CreateMessage(context.Background(), &dto.CreateMessageReq{
			RoomId:          in.RoomID.String(),
			Content:         in.Content,
			MentionEveryone: in.MentionEveryone,
			Mentions:        in.Mentions,
		})

		if err != nil {
			return nil, err
		}

		return &OutboundMessage{
			Type:      MessageTypeText,
			ID:        saved.ID,
			RoomID:    saved.RoomId,
			AuthorID:  saved.AuthorId,
			Content:   saved.Content,
			Timestamp: saved.SentAt,
		}, nil

	}
}
