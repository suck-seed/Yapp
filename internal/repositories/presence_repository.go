// internal/repositories/presence_repository.go
package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/suck-seed/yapp/internal/models"
)

type IPresenceRepository interface {
	MarkConnected(ctx context.Context, userID uuid.UUID, connectionID string, ttl time.Duration) (*models.UserPresence, error)
	RefreshConnection(ctx context.Context, userID uuid.UUID, connectionID string, ttl time.Duration) error
	MarkDisconnected(ctx context.Context, userID uuid.UUID, connectionID string) (*models.UserPresence, error)

	SetManualStatus(ctx context.Context, userID uuid.UUID, status models.PresenceStatus, ttl time.Duration) (*models.UserPresence, error)
	GetUserPresence(ctx context.Context, userID uuid.UUID) (*models.UserPresence, error)
	GetManyPresences(ctx context.Context, userIDs []uuid.UUID) ([]*models.UserPresence, error)

	SetTyping(ctx context.Context, roomID uuid.UUID, userID uuid.UUID, ttl time.Duration) error
	StopTyping(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) error
	GetTypingUsers(ctx context.Context, roomID uuid.UUID) ([]uuid.UUID, error)
}

type presenceRepository struct {
	client *redis.Client
}

func NewPresenceRepository(client *redis.Client) IPresenceRepository {
	return &presenceRepository{client: client}
}

func presenceHashKey(userID uuid.UUID) string {
	return fmt.Sprintf("presence:user:%s", userID.String())
}

func presenceHeartbeatKey(userID uuid.UUID) string {
	return fmt.Sprintf("presence:user:%s:heartbeat", userID.String())
}

func presenceConnectionsKey(userID uuid.UUID) string {
	return fmt.Sprintf("presence:user:%s:connections", userID.String())
}

func typingRoomKey(roomID uuid.UUID) string {
	return fmt.Sprintf("typing:room:%s", roomID.String())
}

func typingUserKey(roomID uuid.UUID, userID uuid.UUID) string {
	return fmt.Sprintf("typing:room:%s:user:%s", roomID.String(), userID.String())
}

func (r *presenceRepository) MarkConnected(ctx context.Context, userID uuid.UUID, connectionID string, ttl time.Duration) (*models.UserPresence, error) {
	now := time.Now().UTC()

	current, _ := r.GetUserPresence(ctx, userID)
	status := models.PresenceStatusOnline

	if current != nil {
		if current.Status == models.PresenceStatusAway || current.Status == models.PresenceStatusBusy {
			status = current.Status
		}
	}

	pipe := r.client.TxPipeline()

	pipe.SAdd(ctx, presenceConnectionsKey(userID), connectionID)
	pipe.Expire(ctx, presenceConnectionsKey(userID), ttl)
	pipe.Set(ctx, presenceHeartbeatKey(userID), "1", ttl)

	pipe.HSet(ctx, presenceHashKey(userID), map[string]any{
		"user_id":      userID.String(),
		"status":       string(status),
		"updated_at":   now.Format(time.RFC3339Nano),
		"last_seen_at": now.Format(time.RFC3339Nano),
	})

	if _, err := pipe.Exec(ctx); err != nil {
		return nil, err
	}

	return r.GetUserPresence(ctx, userID)
}

func (r *presenceRepository) RefreshConnection(ctx context.Context, userID uuid.UUID, connectionID string, ttl time.Duration) error {
	pipe := r.client.TxPipeline()

	pipe.SAdd(ctx, presenceConnectionsKey(userID), connectionID)
	pipe.Expire(ctx, presenceConnectionsKey(userID), ttl)
	pipe.Set(ctx, presenceHeartbeatKey(userID), "1", ttl)

	_, err := pipe.Exec(ctx)
	return err
}

func (r *presenceRepository) MarkDisconnected(ctx context.Context, userID uuid.UUID, connectionID string) (*models.UserPresence, error) {
	now := time.Now().UTC()

	if err := r.client.SRem(ctx, presenceConnectionsKey(userID), connectionID).Err(); err != nil {
		return nil, err
	}

	count, err := r.client.SCard(ctx, presenceConnectionsKey(userID)).Result()
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return r.GetUserPresence(ctx, userID)
	}

	pipe := r.client.TxPipeline()

	pipe.Del(ctx, presenceHeartbeatKey(userID))
	pipe.HSet(ctx, presenceHashKey(userID), map[string]any{
		"user_id":      userID.String(),
		"status":       string(models.PresenceStatusOffline),
		"updated_at":   now.Format(time.RFC3339Nano),
		"last_seen_at": now.Format(time.RFC3339Nano),
	})

	if _, err := pipe.Exec(ctx); err != nil {
		return nil, err
	}

	return r.GetUserPresence(ctx, userID)
}

func (r *presenceRepository) SetManualStatus(ctx context.Context, userID uuid.UUID, status models.PresenceStatus, ttl time.Duration) (*models.UserPresence, error) {
	now := time.Now().UTC()

	pipe := r.client.TxPipeline()

	if status == models.PresenceStatusOffline {
		pipe.Del(ctx, presenceHeartbeatKey(userID))
	} else {
		pipe.Set(ctx, presenceHeartbeatKey(userID), "1", ttl)
	}

	pipe.HSet(ctx, presenceHashKey(userID), map[string]any{
		"user_id":      userID.String(),
		"status":       string(status),
		"updated_at":   now.Format(time.RFC3339Nano),
		"last_seen_at": now.Format(time.RFC3339Nano),
	})

	if _, err := pipe.Exec(ctx); err != nil {
		return nil, err
	}

	return r.GetUserPresence(ctx, userID)
}

func (r *presenceRepository) GetUserPresence(ctx context.Context, userID uuid.UUID) (*models.UserPresence, error) {
	values, err := r.client.HGetAll(ctx, presenceHashKey(userID)).Result()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	if len(values) == 0 {
		return &models.UserPresence{
			UserID:    userID,
			Status:    models.PresenceStatusOffline,
			UpdatedAt: now,
		}, nil
	}

	status := models.PresenceStatus(values["status"])
	heartbeatExists, err := r.client.Exists(ctx, presenceHeartbeatKey(userID)).Result()
	if err != nil {
		return nil, err
	}

	if heartbeatExists == 0 {
		status = models.PresenceStatusOffline
	}

	updatedAt := now
	if values["updated_at"] != "" {
		if parsed, err := time.Parse(time.RFC3339Nano, values["updated_at"]); err == nil {
			updatedAt = parsed
		}
	}

	var lastSeenAt *time.Time
	if values["last_seen_at"] != "" {
		if parsed, err := time.Parse(time.RFC3339Nano, values["last_seen_at"]); err == nil {
			lastSeenAt = &parsed
		}
	}

	return &models.UserPresence{
		UserID:     userID,
		Status:     status,
		LastSeenAt: lastSeenAt,
		UpdatedAt:  updatedAt,
	}, nil
}

func (r *presenceRepository) GetManyPresences(ctx context.Context, userIDs []uuid.UUID) ([]*models.UserPresence, error) {
	out := make([]*models.UserPresence, 0, len(userIDs))

	for _, userID := range userIDs {
		p, err := r.GetUserPresence(ctx, userID)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}

	return out, nil
}

func (r *presenceRepository) SetTyping(ctx context.Context, roomID uuid.UUID, userID uuid.UUID, ttl time.Duration) error {
	pipe := r.client.TxPipeline()

	pipe.SAdd(ctx, typingRoomKey(roomID), userID.String())
	pipe.Set(ctx, typingUserKey(roomID, userID), "1", ttl)
	pipe.Expire(ctx, typingRoomKey(roomID), ttl)

	_, err := pipe.Exec(ctx)
	return err
}

func (r *presenceRepository) StopTyping(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) error {
	pipe := r.client.TxPipeline()

	pipe.SRem(ctx, typingRoomKey(roomID), userID.String())
	pipe.Del(ctx, typingUserKey(roomID, userID))

	_, err := pipe.Exec(ctx)
	return err
}

func (r *presenceRepository) GetTypingUsers(ctx context.Context, roomID uuid.UUID) ([]uuid.UUID, error) {
	raw, err := r.client.SMembers(ctx, typingRoomKey(roomID)).Result()
	if err != nil {
		return nil, err
	}

	out := make([]uuid.UUID, 0, len(raw))

	for _, current := range raw {
		userID, err := uuid.Parse(current)
		if err != nil {
			continue
		}

		exists, err := r.client.Exists(ctx, typingUserKey(roomID, userID)).Result()
		if err != nil {
			return nil, err
		}

		if exists == 0 {
			_ = r.client.SRem(ctx, typingRoomKey(roomID), current).Err()
			continue
		}

		out = append(out, userID)
	}

	return out, nil
}
