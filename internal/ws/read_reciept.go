// internal/ws/read_receipt.go
package ws

import (
	"context"

	"github.com/suck-seed/yapp/internal/auth"
	dto "github.com/suck-seed/yapp/internal/dto/message"
	"github.com/suck-seed/yapp/internal/services"
)

type ReadReceiptFunction func(ctx context.Context, in *dto.InboundMessage) (*dto.OutboundMessage, error)

func MakeReadReceiptFunction(messageService services.IMessageService) ReadReceiptFunction {
	return func(ctx context.Context, in *dto.InboundMessage) (*dto.OutboundMessage, error) {
		if in.MessageID == nil {
			return nil, nil
		}

		userInfo := &auth.UserInfo{
			ID: in.UserID,
		}

		read, err := messageService.MarkMessageRead(ctx, userInfo, in.RoomID, *in.MessageID)
		if err != nil {
			return nil, err
		}

		return &dto.OutboundMessage{
			Type:      dto.MessageTypeRead,
			RoomID:    read.RoomID,
			AuthorID:  read.UserID,
			MessageID: &read.MessageID,
			ReadBy:    &read.UserID,
			ReadAt:    &read.ReadAt,
			SentAt:    read.ReadAt,
		}, nil
	}
}
