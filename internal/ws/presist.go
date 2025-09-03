package ws

import (
	"context"

	"github.com/suck-seed/yapp/internal/services"
)

func MakePresistFunc(messageService services.IMessageService, userService services.IUserService) PersistFunc {
	return func(ctx context.Context, in *InboundMessage) (*OutboundMessage, error) {

	}
}
