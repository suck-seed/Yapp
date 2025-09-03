package ws

import (
	"context"

	"github.com/suck-seed/yapp/internal/services"
)

// MakePresistFunc : Performs various actions and pushes it to db
func MakePresistFunc(messageService services.IMessageService, userService services.IUserService) PersistFunc {
	return func(ctx context.Context, in *InboundMessage) (*OutboundMessage, error) {

	}
}
