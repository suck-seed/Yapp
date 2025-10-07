package services

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/auth"
	"github.com/suck-seed/yapp/internal/dto"
	"github.com/suck-seed/yapp/internal/models"
	"github.com/suck-seed/yapp/internal/repositories"
	"github.com/suck-seed/yapp/internal/utils"
)

type IUserService interface {
	Signup(c context.Context, req *dto.SignupUserReq) (*dto.SignupUserRes, error)
	Signin(c context.Context, req *dto.SigninUserReq) (*dto.SigninUserRes, error)

	GetUserMe(c context.Context) (*models.User, error)
	GetUserById(c context.Context, userId uuid.UUID) (*models.User, error)

	UpdateUserMe(c context.Context, req *dto.UpdateUserMeReq) (*models.User, error)
}

// userService : Behaves like a class, and implements IUserService's methods
type userService struct {
	repositories.IUserRepository
	timeout time.Duration
	mu      sync.RWMutex
}

// NewUserService : Constructor to return a new IUserService with all the user service methods
func NewUserService(repository repositories.IUserRepository) IUserService {
	return &userService{
		repository,
		time.Duration(2) * time.Second,
		sync.RWMutex{},
	}
}

// Methods
func (s *userService) Signup(c context.Context, req *dto.SignupUserReq) (*dto.SignupUserRes, error) {

	// interface that provides a way to control lifecycle, cancellation and prppaagation of requests
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// sanitize the inputs
	canonUsername, err := utils.SanitizeUsername(req.Username)
	if err != nil {
		return nil, utils.ErrorInvalidUsername
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

	// check username, email and number for existing records
	userByUsername, _ := s.IUserRepository.GetUserByUsername(ctx, canonUsername)
	userByEmail, _ := s.IUserRepository.GetUserByEmail(ctx, canonEmail)

	// userByNumber, _ := s.IUserRepository.GetUserByNumber(ctx, canonPhone)

	if userByUsername != nil {
		return nil, utils.ErrorUsernameExists
	}
	if userByEmail != nil {
		return nil, utils.ErrorEmailExists
	}

	// generate id
	id, err := uuid.NewV7()
	if err != nil {
		return nil, utils.ErrorInternal
	}

	// hash password
	password_hash, err := utils.HashPassword(canonPassword)

	if err != nil {
		return nil, utils.ErrorInternal
	}

	user := &models.User{
		ID:       id,
		Username: canonUsername,
		Email:    canonEmail,
		// PhoneNumber:  canonPhone,
		PasswordHash: password_hash,
		DisplayName:  canonDisplayName,
	}

	// calling the repo
	userCRES, err := s.IUserRepository.CreateUser(ctx, user)
	if err != nil {
		print(err)
		return nil, utils.ErrorCreatingUser
	}

	// create a response
	return &dto.SignupUserRes{
		ID:       userCRES.ID,
		Username: userCRES.Username,
	}, nil
}

func (s *userService) Signin(c context.Context, req *dto.SigninUserReq) (*dto.SigninUserRes, error) {

	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	user := &models.User{}

	canonEmail, err := utils.SanitizeEmail(req.Email)
	if err == nil {
		user, _ = s.IUserRepository.GetUserWithPasswordHashByEmail(ctx, canonEmail)
	}

	canonPassword, err := utils.SanitizePasswordPolicy(req.Password)
	if err != nil {
		return nil, utils.ErrorInvalidPassword
	}

	// Handle user not existing
	if user == nil {
		return nil, utils.ErrorUserNotFound
	}

	// hash req.Password and check if matches with pass from user
	err = utils.VerifyPassword(user.PasswordHash, canonPassword)
	if err != nil {
		return nil, utils.ErrorWrongPassword
	}

	// jwt
	signedToken, err := auth.GetSignedToken(user)
	if err != nil {
		return nil, utils.ErrorCreatingUser
	}

	userMe := dto.ToUserMe(*user)

	return &dto.SigninUserRes{
		AccessToken: signedToken,
		Success:     true,
		UserMe:      userMe,
	}, nil
}

func (s *userService) GetUserMe(c context.Context) (*models.User, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// Extract user info from context (already validated by middleware)
	userIdString, _, err := auth.CurrentUserFromContext(c)
	if err != nil {
		return nil, utils.ErrorUserNotFound
	}

	// parse uuid
	userId, err := uuid.Parse(userIdString)
	if err != nil {
		return nil, err
	}

	// Fetch the user from the repository
	user, err := s.IUserRepository.GetUserById(ctx, userId)
	if err != nil {
		return nil, utils.ErrorUserNotFound
	}

	if user == nil {
		return nil, utils.ErrorUserNotFound
	}

	return user, nil
}

func (s *userService) UpdateUserMe(c context.Context, req *dto.UpdateUserMeReq) (*models.User, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// Extract user info from context (already validated by middleware)
	userIdString, _, err := auth.CurrentUserFromContext(c)
	if err != nil {
		return nil, utils.ErrorUserNotFound
	}

	// parse uuid
	userId, err := uuid.Parse(userIdString)
	if err != nil {
		return nil, err
	}

	user, err := s.IUserRepository.UpdateUserById(ctx, userId, req)
	if err != nil {
		return nil, utils.ErrorUserNotFound
	}

	if user == nil {
		return nil, utils.ErrorUserNotFound
	}

	return user, nil
}

func (s *userService) GetUserById(c context.Context, userId uuid.UUID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// Fetch the user from the repository
	user, err := s.IUserRepository.GetUserById(ctx, userId)
	if err != nil {
		return nil, utils.ErrorUserNotFound
	}

	if user == nil {
		return nil, utils.ErrorUserNotFound
	}

	return user, nil
}
