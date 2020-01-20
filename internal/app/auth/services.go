package auth

import (
	"context"

	"github.com/sirupsen/logrus"

	"praslar.com/gotasma/internal/app/status"
	"praslar.com/gotasma/internal/app/types"
	"praslar.com/gotasma/internal/pkg/jwt"
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

func (s *Service) Auth(ctx context.Context, email, password string) (string, *types.User, error) {
	user, err := s.authenticator.Auth(ctx, email, password)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to login with local, err: %v", err)
		return "", nil, status.Auth().InvalidUserPassword
	}

	token, err := s.jwtSigner.Sign(userToClaims(user, jwt.DefaultLifeTime))
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to generate JWT token, err: %v", err)
		return "", nil, err
	}
	return token, user, nil
}
