package auth

import (
	"time"

	"praslar.com/gotasma/internal/app/types"
	"praslar.com/gotasma/internal/pkg/jwt"
)

func userToClaims(user *types.User, lifeTime time.Duration) jwt.Claims {
	return jwt.Claims{
		Role:      user.Role,
		UserID:    user.UserID,
		ProjectID: user.ProjectID,
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
		ProjectID: claims.ProjectID,
		CreaterID: claims.CreaterID,
	}
}
