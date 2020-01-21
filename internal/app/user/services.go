package user

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"praslar.com/gotasma/internal/app/status"
	"praslar.com/gotasma/internal/app/types"
	"praslar.com/gotasma/internal/pkg/db"
	"praslar.com/gotasma/internal/pkg/uuid"
	"praslar.com/gotasma/internal/pkg/validator"
)

type (
	Repository interface {
		Create(context.Context, *types.User) (string, error)
		FindByEmail(ctx context.Context, email string) (*types.User, error)
	}

	Service struct {
		repo Repository
	}
)

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Register(ctx context.Context, req *types.RegisterRequest) (*types.User, error) {

	if err := validator.Validate(req); err != nil {
		return nil, err
	}

	existingUser, err := s.repo.FindByEmail(ctx, req.Email)

	if err != nil && !db.IsErrNotFound(err) {
		logrus.WithContext(ctx).Errorf("failed to check existing user by email, err: %v", err)
		return nil, fmt.Errorf("failed to check existing user by email: %w", err)
	}

	if existingUser != nil {
		logrus.Infof("email already registered")
		logrus.WithContext(ctx).Debug("email already registered")
		return nil, status.User().DuplicatedEmail
	}

	password, err := s.generatePassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password: %w", err)
	}

	user := &types.User{
		Password:  password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Role:      req.Role,
		UserID:    uuid.New(),
	}

	if _, err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}
	
	return user.Strip(), nil
}

func (s *Service) generatePassword(pass string) (string, error) {
	rs, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate password: %w", err)
	}
	return string(rs), nil
}

func (s *Service) Auth(ctx context.Context, email, password string) (*types.User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil && !db.IsErrNotFound(err) {
		logrus.WithContext(ctx).Errorf("failed to check existing user by email, err: %v", err)
		return nil, status.Gen().Internal
	}
	if db.IsErrNotFound(err) {
		logrus.WithContext(ctx).Debugf("user not found, email: %s", email)
		return nil, status.Auth().InvalidUserPassword
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logrus.WithContext(ctx).Error("invalid password")
		return nil, status.Auth().InvalidUserPassword
	}
	return user.Strip(), nil
}
