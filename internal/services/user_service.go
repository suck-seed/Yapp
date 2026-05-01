package services

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suck-seed/yapp/internal/auth"
	"github.com/suck-seed/yapp/internal/database"
	dto "github.com/suck-seed/yapp/internal/dto/user"
	"github.com/suck-seed/yapp/internal/models"
	"github.com/suck-seed/yapp/internal/repositories"
	"github.com/suck-seed/yapp/internal/utils"
)

type IUserService interface {
	Signup(c context.Context, req *dto.SignupUserReq) (*dto.SignupUserRes, error)
	Signin(c context.Context, req *dto.SigninUserReq) (*dto.SigninUserRes, error)

	GetUserMe(c context.Context, userInfo *auth.UserInfo) (*dto.UserMe, error)
	GetUserById(c context.Context, userID uuid.UUID) (*models.User, error)
	GetUserPublic(c context.Context, currentUserID uuid.UUID, targetUserID uuid.UUID) (*dto.UserPublic, error)
	GetMutualFriends(c context.Context, currentUserID uuid.UUID, targetUserID uuid.UUID) (*dto.MutualFriendRes, error)

	UpdateUserMe(c context.Context, userInfo *auth.UserInfo, req *dto.UpdateUserMeReq) (*dto.UserMe, error)
	UpdateUsername(c context.Context, userInfo *auth.UserInfo, req *dto.UpdateUsernameReq) (*dto.UserMe, error)
	UpdateEmail(c context.Context, userInfo *auth.UserInfo, req *dto.UpdateEmailReq) (*dto.UserMe, error)
	DeleteMe(c context.Context, userInfo *auth.UserInfo) error

	SendFriendRequest(c context.Context, userInfo *auth.UserInfo, req *dto.SendFriendRequestReq) (*dto.FriendRequestRes, error)
	RespondFriendRequest(c context.Context, userInfo *auth.UserInfo, requestID uuid.UUID, req *dto.RespondFriendRequestReq) error
	Unfriend(c context.Context, userInfo *auth.UserInfo, targetUserID uuid.UUID) error
	GetMyFriends(c context.Context, userInfo *auth.UserInfo) (*dto.FriendListRes, error)

	UpsertMyAppLink(c context.Context, userInfo *auth.UserInfo, req *dto.UpsertAppLinkReq) (*dto.UpsertAppLinkRes, error)
	DeleteMyAppLink(c context.Context, userInfo *auth.UserInfo, provider models.AppProvider) error
}

type userService struct {
	repositories.IUserRepository
	pool    *pgxpool.Pool
	timeout time.Duration
	mu      sync.RWMutex
}

func NewUserService(repository repositories.IUserRepository, pool *pgxpool.Pool) IUserService {
	return &userService{
		repository,
		pool,
		time.Duration(2) * time.Second,
		sync.RWMutex{},
	}
}

func (s *userService) Signup(c context.Context, req *dto.SignupUserReq) (*dto.SignupUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	canonUsername, err := utils.SanitizeUsername(req.Username)
	if err != nil {
		return nil, utils.ErrorInvalidUserName
	}
	canonPassword, err := utils.SanitizePasswordPolicy(req.Password)
	if err != nil {
		return nil, utils.ErrorInvalidPassword
	}
	canonEmail, err := utils.SanitizeEmail(req.Email)
	if err != nil {
		return nil, utils.ErrorInvalidEmail
	}
	canonDisplayName, err := utils.SanitizeDisplayName(req.DisplayName)
	if err != nil {
		return nil, utils.ErrorInvalidDisplayName
	}

	userByUsername, _ := s.IUserRepository.GetUserByUsername(ctx, runner, canonUsername)
	userByEmail, _ := s.IUserRepository.GetUserByEmail(ctx, runner, canonEmail)

	if userByUsername != nil {
		return nil, utils.ErrorUsernameExists
	}
	if userByEmail != nil {
		return nil, utils.ErrorEmailExists
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, utils.ErrorInternal
	}

	passwordHash, err := utils.HashPassword(canonPassword)
	if err != nil {
		return nil, utils.ErrorInternal
	}

	user := &models.User{
		ID:           id,
		Username:     canonUsername,
		Email:        canonEmail,
		PasswordHash: passwordHash,
		DisplayName:  canonDisplayName,
	}

	userCRES, err := s.IUserRepository.CreateUser(ctx, runner, user)
	if err != nil {
		return nil, utils.ErrorCreatingUser
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	return &dto.SignupUserRes{
		ID:       userCRES.ID,
		Username: userCRES.Username,
		Email:    userCRES.Email,
	}, nil
}

func (s *userService) Signin(c context.Context, req *dto.SigninUserReq) (*dto.SigninUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	canonEmail, err := utils.SanitizeEmail(req.Email)
	if err != nil {
		return nil, utils.ErrorInvalidEmail
	}

	canonPassword, err := utils.SanitizePasswordPolicy(req.Password)
	if err != nil {
		return nil, utils.ErrorInvalidPassword
	}

	user, err := s.IUserRepository.GetUserWithPasswordHashByEmail(ctx, runner, canonEmail)
	if err != nil {
		return nil, utils.ErrorUserNotFound
	}

	err = utils.VerifyPassword(user.PasswordHash, canonPassword)
	if err != nil {
		return nil, utils.ErrorWrongPassword
	}

	signedToken, err := auth.GetSignedToken(user)
	if err != nil {
		return nil, utils.ErrorCreatingUser
	}

	userMe, err := s.buildUserMe(ctx, runner, user)
	if err != nil {
		return nil, err
	}

	return &dto.SigninUserRes{
		AccessToken: signedToken,
		Success:     true,
		UserMe:      *userMe,
	}, nil
}

func (s *userService) GetUserMe(c context.Context, userInfo *auth.UserInfo) (*dto.UserMe, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	user, err := s.IUserRepository.GetUserById(ctx, runner, userInfo.ID)
	if err != nil {
		return nil, utils.ErrorUserNotFound
	}

	return s.buildUserMe(ctx, runner, user)
}

func (s *userService) GetUserById(c context.Context, userID uuid.UUID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	user, err := s.IUserRepository.GetUserById(ctx, runner, userID)
	if err != nil {
		return nil, utils.ErrorUserNotFound
	}

	return user, nil
}

func (s *userService) GetUserPublic(c context.Context, currentUserID uuid.UUID, targetUserID uuid.UUID) (*dto.UserPublic, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	user, err := s.IUserRepository.GetUserById(ctx, runner, targetUserID)
	if err != nil {
		return nil, utils.ErrorUserNotFound
	}

	public := dto.ToUserPublic(*user)

	friendCount, err := s.IUserRepository.CountFriends(ctx, runner, targetUserID)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	public.FriendCount = friendCount

	isFriend := false
	if currentUserID != uuid.Nil && currentUserID != targetUserID {
		isFriend, err = s.IUserRepository.AreFriends(ctx, runner, currentUserID, targetUserID)
		if err != nil {
			return nil, utils.ErrorInternal
		}
	}
	public.IsFriend = isFriend

	if currentUserID != uuid.Nil && currentUserID != targetUserID {
		mutualCount, err := s.IUserRepository.CountMutualFriends(ctx, runner, currentUserID, targetUserID)
		if err != nil {
			return nil, utils.ErrorInternal
		}
		public.MutualFriendCount = mutualCount
	}

	links, err := s.IUserRepository.GetUserAppLinks(ctx, runner, targetUserID, true)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	public.AppLinks = mapAppLinks(links)

	return &public, nil
}

func (s *userService) GetMutualFriends(c context.Context, currentUserID uuid.UUID, targetUserID uuid.UUID) (*dto.MutualFriendRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	_, err = s.IUserRepository.GetUserById(ctx, runner, targetUserID)
	if err != nil {
		return nil, utils.ErrorUserNotFound
	}

	users, err := s.IUserRepository.GetMutualFriends(ctx, runner, currentUserID, targetUserID)
	if err != nil {
		return nil, utils.ErrorInternal
	}

	out := make([]*dto.UserPublic, 0, len(users))
	for _, u := range users {
		current := dto.ToUserPublic(*u)
		links, err := s.IUserRepository.GetUserAppLinks(ctx, runner, u.ID, true)
		if err != nil {
			return nil, utils.ErrorInternal
		}
		current.AppLinks = mapAppLinks(links)
		out = append(out, &current)
	}

	return &dto.MutualFriendRes{
		Users: out,
		Total: len(out),
	}, nil
}

func (s *userService) UpdateUserMe(c context.Context, userInfo *auth.UserInfo, req *dto.UpdateUserMeReq) (*dto.UserMe, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	if req.DisplayName == nil && req.Description == nil && req.AvatarURL == nil && req.AvatarThumbnailURL == nil && req.FriendPolicy == nil {
		return nil, utils.ErrorNoFieldsToUpdate
	}

	fields := make(map[string]any)

	if req.DisplayName != nil {
		canonDisplayName, err := utils.SanitizeDisplayName(*req.DisplayName)
		if err != nil {
			return nil, utils.ErrorInvalidDisplayName
		}
		fields["display_name"] = canonDisplayName
	}

	if req.Description != nil {
		canonDescription, err := utils.SanitizeText(req.Description)
		if err != nil {
			return nil, err
		}
		fields["description"] = canonDescription
	}

	if req.AvatarURL != nil {
		fields["avatar_url"] = req.AvatarURL
	}

	if req.AvatarThumbnailURL != nil {
		fields["avatar_thumbnail_url"] = req.AvatarThumbnailURL
	}

	if req.FriendPolicy != nil {
		fields["friend_policy"] = *req.FriendPolicy
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	user, err := s.IUserRepository.UpdateUserById(ctx, runner, userInfo.ID, fields)
	if err != nil {
		return nil, utils.ErrorUserNotFound
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()

	return s.buildUserMe(ctx, database.NewConnWrapper(conn), user)
}

func (s *userService) UpdateUsername(c context.Context, userInfo *auth.UserInfo, req *dto.UpdateUsernameReq) (*dto.UserMe, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	canonUsername, err := utils.SanitizeUsername(req.NewUsername)
	if err != nil {
		return nil, utils.ErrorInvalidUserName
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	existing, _ := s.IUserRepository.GetUserByUsername(ctx, runner, canonUsername)
	if existing != nil && existing.ID != userInfo.ID {
		return nil, utils.ErrorUsernameExists
	}

	updated, err := s.IUserRepository.UpdateUsername(ctx, runner, userInfo.ID, canonUsername)
	if err != nil {
		return nil, utils.ErrorUserNotFound
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()

	return s.buildUserMe(ctx, database.NewConnWrapper(conn), updated)
}

func (s *userService) UpdateEmail(c context.Context, userInfo *auth.UserInfo, req *dto.UpdateEmailReq) (*dto.UserMe, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	canonEmail, err := utils.SanitizeEmail(req.NewEmail)
	if err != nil {
		return nil, utils.ErrorInvalidEmail
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	existing, _ := s.IUserRepository.GetUserByEmail(ctx, runner, canonEmail)
	if existing != nil && existing.ID != userInfo.ID {
		return nil, utils.ErrorEmailExists
	}

	updated, err := s.IUserRepository.UpdateEmail(ctx, runner, userInfo.ID, canonEmail)
	if err != nil {
		return nil, utils.ErrorUserNotFound
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()

	return s.buildUserMe(ctx, database.NewConnWrapper(conn), updated)
}

func (s *userService) DeleteMe(c context.Context, userInfo *auth.UserInfo) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	err = s.IUserRepository.DeleteUserById(ctx, runner, userInfo.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return utils.ErrorUserNotFound
		}
		return utils.ErrorInternal
	}

	if err := runner.Commit(ctx); err != nil {
		return utils.ErrorInternal
	}

	return nil
}

func (s *userService) SendFriendRequest(c context.Context, userInfo *auth.UserInfo, req *dto.SendFriendRequestReq) (*dto.FriendRequestRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	if req.ReceiverID == userInfo.ID {
		return nil, utils.ErrorInvalidInput
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	receiver, err := s.IUserRepository.GetUserById(ctx, runner, req.ReceiverID)
	if err != nil {
		return nil, utils.ErrorUserNotFound
	}

	alreadyFriends, err := s.IUserRepository.AreFriends(ctx, runner, userInfo.ID, req.ReceiverID)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	if alreadyFriends {
		return nil, utils.ErrorFriendAlreadyExists
	}

	outgoingExists, err := s.IUserRepository.DoesFriendRequestExist(ctx, runner, userInfo.ID, req.ReceiverID)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	if outgoingExists {
		return nil, utils.ErrorFriendRequestAlreadyExists
	}

	incomingExists, err := s.IUserRepository.DoesFriendRequestExist(ctx, runner, req.ReceiverID, userInfo.ID)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	if incomingExists {
		return nil, utils.ErrorFriendRequestAlreadyExists
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, utils.ErrorInternal
	}

	saved, err := s.IUserRepository.CreateFriendRequest(ctx, runner, &models.FriendRequest{
		ID:         id,
		SenderID:   userInfo.ID,
		ReceiverID: req.ReceiverID,
	})
	if err != nil {
		return nil, utils.ErrorInternal
	}

	sender, err := s.IUserRepository.GetUserById(ctx, runner, userInfo.ID)
	if err != nil {
		return nil, utils.ErrorUserNotFound
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	return &dto.FriendRequestRes{
		ID:        saved.ID,
		Sender:    dto.ToUserPublic(*sender),
		Receiver:  dto.ToUserPublic(*receiver),
		CreatedAt: saved.CreatedAt,
	}, nil
}

func (s *userService) RespondFriendRequest(c context.Context, userInfo *auth.UserInfo, requestID uuid.UUID, req *dto.RespondFriendRequestReq) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	friendRequest, err := s.IUserRepository.GetFriendRequestByID(ctx, runner, requestID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return utils.ErrorFriendRequestNotFound
		}
		return utils.ErrorInternal
	}

	if friendRequest.ReceiverID != userInfo.ID {
		return utils.ErrorUnauthorizedToHandleFriendRequest
	}

	if req.Action == "decline" {
		_, err := s.IUserRepository.DeleteFriendRequestByID(ctx, runner, requestID)
		if err != nil {
			return utils.ErrorInternal
		}

		if err := runner.Commit(ctx); err != nil {
			return utils.ErrorInternal
		}
		return nil
	}

	_, err = s.IUserRepository.CreateFriendship(ctx, runner, friendRequest.SenderID, friendRequest.ReceiverID)
	if err != nil {
		return utils.ErrorInternal
	}

	_, err = s.IUserRepository.DeleteFriendRequestByID(ctx, runner, requestID)
	if err != nil {
		return utils.ErrorInternal
	}

	if err := runner.Commit(ctx); err != nil {
		return utils.ErrorInternal
	}

	return nil
}

func (s *userService) Unfriend(c context.Context, userInfo *auth.UserInfo, targetUserID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	err = s.IUserRepository.DeleteFriendship(ctx, runner, userInfo.ID, targetUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return utils.ErrorFriendNotFound
		}
		return utils.ErrorInternal
	}

	if err := runner.Commit(ctx); err != nil {
		return utils.ErrorInternal
	}

	return nil
}

func (s *userService) GetMyFriends(c context.Context, userInfo *auth.UserInfo) (*dto.FriendListRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	users, err := s.IUserRepository.ListFriends(ctx, runner, userInfo.ID)
	if err != nil {
		return nil, utils.ErrorInternal
	}

	out := make([]*dto.UserPublic, 0, len(users))
	for _, u := range users {
		current := dto.ToUserPublic(*u)

		friendCount, err := s.IUserRepository.CountFriends(ctx, runner, u.ID)
		if err != nil {
			return nil, utils.ErrorInternal
		}
		current.FriendCount = friendCount
		current.IsFriend = true

		mutualCount, err := s.IUserRepository.CountMutualFriends(ctx, runner, userInfo.ID, u.ID)
		if err != nil {
			return nil, utils.ErrorInternal
		}
		current.MutualFriendCount = mutualCount

		links, err := s.IUserRepository.GetUserAppLinks(ctx, runner, u.ID, true)
		if err != nil {
			return nil, utils.ErrorInternal
		}
		current.AppLinks = mapAppLinks(links)

		out = append(out, &current)
	}

	return &dto.FriendListRes{
		Users: out,
		Total: len(out),
	}, nil
}

func (s *userService) UpsertMyAppLink(c context.Context, userInfo *auth.UserInfo, req *dto.UpsertAppLinkReq) (*dto.UpsertAppLinkRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	id, err := uuid.NewV7()
	if err != nil {
		return nil, utils.ErrorInternal
	}

	saved, err := s.IUserRepository.UpsertAppLink(ctx, runner, &models.UserAppLink{
		ID:            id,
		UserID:        userInfo.ID,
		Provider:      req.Provider,
		URL:           req.URL,
		ShowOnProfile: req.Show,
	})
	if err != nil {
		return nil, utils.ErrorInternal
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	return &dto.UpsertAppLinkRes{
		ID:            saved.ID,
		UserID:        saved.UserID,
		Provider:      saved.Provider,
		URL:           saved.URL,
		ShowOnProfile: saved.ShowOnProfile,
		CreatedAt:     saved.CreatedAt,
		UpdatedAt:     saved.UpdatedAt,
	}, nil
}

func (s *userService) DeleteMyAppLink(c context.Context, userInfo *auth.UserInfo, provider models.AppProvider) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	err = s.IUserRepository.DeleteAppLink(ctx, runner, userInfo.ID, provider)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return utils.ErrorAppLinkNotFound
		}
		return utils.ErrorInternal
	}

	if err := runner.Commit(ctx); err != nil {
		return utils.ErrorInternal
	}

	return nil
}

func (s *userService) buildUserMe(ctx context.Context, runner database.DBRunner, user *models.User) (*dto.UserMe, error) {
	res := dto.ToUserMe(*user)

	links, err := s.IUserRepository.GetUserAppLinks(ctx, runner, user.ID, false)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	res.AppLinks = mapAppLinks(links)

	return &res, nil
}

func mapAppLinks(links []*models.UserAppLink) []dto.AppLink {
	out := make([]dto.AppLink, 0, len(links))
	for _, link := range links {
		out = append(out, dto.AppLink{
			Provider: link.Provider,
			URL:      link.URL,
			Show:     link.ShowOnProfile,
		})
	}
	return out
}
