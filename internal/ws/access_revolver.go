package ws

import (
	"context"

	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/auth"
	"github.com/suck-seed/yapp/internal/services"
)

type AccessResolver func(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]uuid.UUID, error)

func MakeAccessResolver(roomService services.IRoomService) AccessResolver {
	return func(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]uuid.UUID, error) {
		return roomService.GetAccessibleRoomsForUser(ctx, &auth.UserInfo{ID: userID})
	}
}
