package auth

import (
	"context"

	"github.com/gotasma/internal/app/status"
	"github.com/gotasma/internal/app/types"
	"github.com/gotasma/internal/pkg/jwt"

	"github.com/sirupsen/logrus"
)

type (
	Authenticator interface {
		Auth(ctx context.Context, email, password string) (*types.User, error)
	}
	Service struct {
		jwtSigner     jwt.Signer
		authenticator Authenticator
	}
)

func NewService(signer jwt.Signer, authenticator Authenticator) *Service {
	return &Service{
		jwtSigner:     signer,
		authenticator: authenticator,
	}
}

func (s *Service) Auth(ctx context.Context, email, password string) (string, error) {
	user, err := s.authenticator.Auth(ctx, email, password)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to login with local, err: %v", err)
		return "", status.Auth().InvalidUserPassword
	}

	token, err := s.jwtSigner.Sign(userToClaims(user, jwt.DefaultLifeTime))
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to generate JWT token, err: %v", err)
		return "", err
	}
	return token, nil
}
