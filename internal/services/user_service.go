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

	GetUserByID(c context.Context, userId *uuid.UUID) (*models.User, error)
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

	// canonPhone, err := utils.SanitizePhoneE164(req.PhoneNumber)
	// if err != nil {
	// 	return nil, utils.ErrorInvalidPhoneNumber
	// }

	canonDisplayName, err := utils.SanitizeDisplayName(req.DisplayName)
	if err != nil {
		return nil, utils.ErrorInvalidDisplayName
	}

	// to assign username to displayName if null
	// if canonDisplayName == nil {
	// 	canonDisplayName = &canonUsername
	// }

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
	// if userByNumber != nil {
	// 	return nil, utils.ErrorNumberExists
	// }

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
	r, err := s.IUserRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, utils.ErrorCreatingUser
	}

	// create a response

	return &dto.SignupUserRes{
		ID:       r.ID.String(),
		Username: r.Username,
	}, nil

}

func (s *userService) Signin(c context.Context, req *dto.SigninUserReq) (*dto.SigninUserRes, error) {

	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	user := &models.User{}

	canonEmail, err := utils.SanitizeEmail(req.Email)
	if err == nil {
		user, _ = s.IUserRepository.GetUserByEmail(ctx, canonEmail)
	}

	canonUsername, err := utils.SanitizeUsername(req.Email)
	if err == nil {
		user, _ = s.IUserRepository.GetUserByUsername(ctx, canonUsername)
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

	return &dto.SigninUserRes{
		AccessToken: signedToken,
		ID:          user.ID.String(),
		Username:    user.Username,
	}, nil
}

func (s *userService) GetUserByID(c context.Context, userId *uuid.UUID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// Actually fetch the user from the repository
	user, err := s.IUserRepository.GetUserById(ctx, userId)
	if err != nil {
		return nil, utils.ErrorUserNotFound
	}

	if user == nil {
		return nil, utils.ErrorUserNotFound
	}

	return user, nil
}
