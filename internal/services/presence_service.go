package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	dto "github.com/suck-seed/yapp/internal/dto/user"
	"github.com/suck-seed/yapp/internal/models"
	"github.com/suck-seed/yapp/internal/repositories"
	"github.com/suck-seed/yapp/internal/utils"
)

type IPresenceService interface {
	MarkConnected(ctx context.Context, userID uuid.UUID, connectionID uuid.UUID) (*dto.UserPresenceRes, error)
	RefreshConnection(ctx context.Context, userID uuid.UUID, connectionID uuid.UUID) error
	MarkDisconnected(ctx context.Context, userID uuid.UUID, connectionID uuid.UUID) (*dto.UserPresenceRes, error)

	SetManualStatus(ctx context.Context, userID uuid.UUID, status models.PresenceStatus) (*dto.UserPresenceRes, error)
	GetUserPresence(ctx context.Context, userID uuid.UUID) (*dto.UserPresenceRes, error)
	GetManyPresences(ctx context.Context, userIDs []uuid.UUID) ([]*dto.UserPresenceRes, error)

	SetTyping(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) error
	StopTyping(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) error
	GetTypingUsers(ctx context.Context, roomID uuid.UUID) ([]uuid.UUID, error)
}

type presenceService struct {
	repositories.IPresenceRepository

	timeout time.Duration
	ttl     time.Duration
}

func NewPresenceService(presenceRepo repositories.IPresenceRepository) IPresenceService {
	return &presenceService{
		IPresenceRepository: presenceRepo,
		timeout:             2 * time.Second,
		ttl:                 90 * time.Second,
	}
}

func presenceToRes(p *models.UserPresence) *dto.UserPresenceRes {
	return &dto.UserPresenceRes{
		UserID:     p.UserID,
		Status:     p.Status,
		LastSeenAt: p.LastSeenAt,
		UpdatedAt:  p.UpdatedAt,
	}
}

func isValidPresenceStatus(status models.PresenceStatus) bool {
	switch status {
	case models.PresenceStatusOnline,
		models.PresenceStatusOffline,
		models.PresenceStatusAway,
		models.PresenceStatusBusy:
		return true
	default:
		return false
	}
}

func (s *presenceService) MarkConnected(ctx context.Context, userID uuid.UUID, connectionID uuid.UUID) (*dto.UserPresenceRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	p, err := s.IPresenceRepository.MarkConnected(ctx, userID, connectionID, s.ttl)
	if err != nil {
		return nil, utils.ErrorInternal
	}

	return presenceToRes(p), nil
}

func (s *presenceService) RefreshConnection(ctx context.Context, userID uuid.UUID, connectionID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if err := s.IPresenceRepository.RefreshConnection(ctx, userID, connectionID, s.ttl); err != nil {
		return utils.ErrorInternal
	}

	return nil
}

func (s *presenceService) MarkDisconnected(ctx context.Context, userID uuid.UUID, connectionID uuid.UUID) (*dto.UserPresenceRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	p, err := s.IPresenceRepository.MarkDisconnected(ctx, userID, connectionID)
	if err != nil {
		return nil, utils.ErrorInternal
	}

	return presenceToRes(p), nil
}

func (s *presenceService) SetManualStatus(ctx context.Context, userID uuid.UUID, status models.PresenceStatus) (*dto.UserPresenceRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if !isValidPresenceStatus(status) {
		return nil, utils.ErrorInvalidInput
	}

	p, err := s.IPresenceRepository.SetManualStatus(ctx, userID, status, s.ttl)
	if err != nil {
		return nil, utils.ErrorInternal
	}

	return presenceToRes(p), nil
}

func (s *presenceService) GetUserPresence(ctx context.Context, userID uuid.UUID) (*dto.UserPresenceRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	p, err := s.IPresenceRepository.GetUserPresence(ctx, userID)
	if err != nil {
		return nil, utils.ErrorInternal
	}

	return presenceToRes(p), nil
}

func (s *presenceService) GetManyPresences(ctx context.Context, userIDs []uuid.UUID) ([]*dto.UserPresenceRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	rows, err := s.IPresenceRepository.GetManyPresences(ctx, userIDs)
	if err != nil {
		return nil, utils.ErrorInternal
	}

	out := make([]*dto.UserPresenceRes, 0, len(rows))
	for _, row := range rows {
		out = append(out, presenceToRes(row))
	}

	return out, nil
}

func (s *presenceService) SetTyping(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if err := s.IPresenceRepository.SetTyping(ctx, roomID, userID, 5*time.Second); err != nil {
		return utils.ErrorInternal
	}

	return nil
}

func (s *presenceService) StopTyping(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if err := s.IPresenceRepository.StopTyping(ctx, roomID, userID); err != nil {
		return utils.ErrorInternal
	}

	return nil
}

func (s *presenceService) GetTypingUsers(ctx context.Context, roomID uuid.UUID) ([]uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	users, err := s.IPresenceRepository.GetTypingUsers(ctx, roomID)
	if err != nil {
		return nil, utils.ErrorInternal
	}

	return users, nil
}
