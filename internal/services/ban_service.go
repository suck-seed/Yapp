package services

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suck-seed/yapp/internal/database"
	dto "github.com/suck-seed/yapp/internal/dto/hall"
	"github.com/suck-seed/yapp/internal/repositories"
)

type IBanService interface {
	BanUser(ctx context.Context, hallID uuid.UUID, userID uuid.UUID) (*dto.BanUserRes, error)
	UnbanUser(ctx context.Context, hallID uuid.UUID, userID uuid.UUID) (*dto.UnbanRes, error)
	GetBanByID(ctx context.Context, hallID uuid.UUID, userID uuid.UUID, banID uuid.UUID) (*dto.BanSummaryRes, error)
	GetAllHallBans(ctx context.Context, hallID uuid.UUID, userID uuid.UUID) (*dto.AllBannedUserRes, error)
}

type banService struct {
	repositories.IBanRepsitory
	repositories.IUserRepository
	repositories.IHallRepository
	repositories.IRoleRepository

	pool    *pgxpool.Pool
	timeout time.Duration
	mu      sync.RWMutex
}

func NewBanService(banRepo repositories.IBanRepsitory, userRepo repositories.IUserRepository, hallRepo repositories.IHallRepository, roleRepo repositories.IRoleRepository, pool *pgxpool.Pool) IBanService {
	return &banService{
		banRepo,
		userRepo,
		hallRepo,
		roleRepo,
		pool,
		time.Duration(2) * time.Second,
		sync.RWMutex{},
	}
}

func (s *banService) BanUser(ctx context.Context, hallID uuid.UUID, userID uuid.UUID) (*dto.BanUserRes, error) {

	return nil, nil
}
func (s *banService) UnbanUser(ctx context.Context, hallID uuid.UUID, userID uuid.UUID) (*dto.UnbanRes, error) {

	return nil, nil
}
func (s *banService) GetBanByID(ctx context.Context, hallID uuid.UUID, userID uuid.UUID, banID uuid.UUID) (*dto.BanSummaryRes, error) {

	return nil, nil
}
func (s *banService) GetAllHallBans(ctx context.Context, hallID uuid.UUID, userID uuid.UUID) (*dto.AllBannedUserRes, error) {

	return nil, nil
}

// helper functions
// canBanUsers - check if current user is hall owner? or has admin role ? or has permission to ban members
func (s *banService) canBanUsers(ctx context.Context, db database.DBRunner, hallID uuid.UUID, userID uuid.UUID) (bool, error) {

	return false, nil
}

// hasHigherPermission : helps check role hierarchy, to detemrine if the ban is possible
// help in determining is ban is possible if, by checking their both banner permission
func (s *banService) hasHigherPermission(ctx context.Context, db database.DBRunner, hallID uuid.UUID, bannerID uuid.UUID, banTargetId uuid.UUID) (bool, error) {

	return false, nil
}
