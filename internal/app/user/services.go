package user

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"praslar.com/gotasma/internal/app/auth"
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
		FindAllDev(ctx context.Context, createrID string) ([]*types.User, error)
		Delete(cxt context.Context, id string) error
	}

	PolicyService interface {
		Validate(ctx context.Context, obj string, act string) error
	}

	Service struct {
		repo   Repository
		policy PolicyService
	}
)

func New(repo Repository, policy PolicyService) *Service {
	return &Service{
		repo:   repo,
		policy: policy,
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
		Role:      types.PM,
		UserID:    uuid.New(),
	}

	if _, err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return user.Strip(), nil
}

func (s *Service) CreateDev(ctx context.Context, req *types.RegisterRequest) (*types.User, error) {

	//Only PM can create dev
	if err := s.policy.Validate(ctx, types.PolicyObjectDev, types.PolicyActionDevCreate); err != nil {
		return nil, err
	}

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

	pm := auth.FromContext(ctx)

	user := &types.User{
		Password:  password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Role:      req.Role,
		UserID:    uuid.New(),
		CreaterID: pm.UserID,
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

func (s *Service) FindAllDev(ctx context.Context) ([]*types.UserInfo, error) {
	if err := s.policy.Validate(ctx, types.PolicyObjectUser, types.PolicyActionUserReadList); err != nil {
		return nil, err
	}

	pm := auth.FromContext(ctx)
	users, err := s.repo.FindAllDev(ctx, pm.UserID)
	info := make([]*types.UserInfo, 0)
	for _, usr := range users {
		info = append(info, &types.UserInfo{
			Email:     usr.Email,
			FirstName: usr.FirstName,
			LastName:  usr.LastName,
			Role:      usr.Role,
			ProjectID: usr.ProjectID,
			CreaterID: usr.CreaterID,
			UserID:    usr.UserID,
		})
	}
	return info, err
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.policy.Validate(ctx, types.PolicyObjectDeleteDev, types.PolicyActionDevDelete); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
