package auth

import (
	"context"
	"net/http"

	"github.com/gotasma/internal/app/status"
	"github.com/gotasma/internal/app/types"
	"github.com/gotasma/internal/pkg/http/respond"
	"github.com/gotasma/internal/pkg/jwt"
	"github.com/sirupsen/logrus"
)

type (
	contextKey string
)

const (
	authContextKey contextKey = "r_auth_user"
)

//UserInfoMiddleware verifier login
func UserInfoMiddleware(verifier jwt.Verifier) func(http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//TODO front end work. - Set Header [Authorizatio-token]
			key := r.Header.Get("Authorization")

			if key == "" {
				inner.ServeHTTP(w, r)
				return
			}

			claims, err := verifier.Verify(key)
			if err != nil {
				logrus.WithContext(r.Context()).Errorf("invalid JWT token, err: %v", err)
				inner.ServeHTTP(w, r)
				return
			}

			newCtx := NewContext(r.Context(), claimsToUser(claims))

			r = r.WithContext(newCtx)
			inner.ServeHTTP(w, r)
		})
	}
}

func NewContext(ctx context.Context, user *types.User) context.Context {
	return context.WithValue(ctx, authContextKey, user)
}

func FromContext(ctx context.Context) *types.User {
	if v, ok := ctx.Value(authContextKey).(*types.User); ok {
		return v
	}
	return nil
}

func RequiredAuthMiddleware(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user := FromContext(r.Context()); user == nil {
			respond.JSON(w, http.StatusUnauthorized, status.Policy().Unauthorized)
			return
		}
		inner.ServeHTTP(w, r)
	})
}
