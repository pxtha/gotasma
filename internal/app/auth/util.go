package auth

import (
	"time"

	"github.com/gotasma/internal/app/types"
	"github.com/gotasma/internal/pkg/jwt"
)

func userToClaims(user *types.User, lifeTime time.Duration) jwt.Claims {
	return jwt.Claims{
		Role:      user.Role,
		UserID:    user.UserID,
		CreaterID: user.CreaterID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(lifeTime).Unix(),
			Id:        user.UserID,
			IssuedAt:  time.Now().Unix(),
			Issuer:    jwt.DefaultIssuer,
		},
	}
}

func claimsToUser(claims *jwt.Claims) *types.User {
	return &types.User{
		Role:      claims.Role,
		UserID:    claims.UserID,
		CreaterID: claims.CreaterID,
	}
}
