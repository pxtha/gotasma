package auth

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"praslar.com/gotasma/internal/app/types"
	"praslar.com/gotasma/internal/pkg/jwt"
)

const (
	TokenCookieName = "_r_token"
	UserCookieName  = "_r_user"
)

func userToClaims(user *types.User, lifeTime time.Duration) jwt.Claims {
	return jwt.Claims{
		Role:   int(user.Role),
		UserID: user.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(lifeTime).Unix(),
			Id:        user.UserID,
			IssuedAt:  time.Now().Unix(),
			Issuer:    jwt.DefaultIssuer,
			Subject:   user.UserID,
		},
	}
}

func createTokenCookie(token string, r *http.Request) *http.Cookie {
	return &http.Cookie{
		Name:     TokenCookieName,
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		Domain:   r.Host,
		Path:     "/",
		Secure:   false,
		HttpOnly: false, // allow client to access this cookie
	}
}
func userInfoCookieValue(user *types.User) (string, error) {
	b, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func createUserInfoCookie(user *types.User, r *http.Request) (*http.Cookie, error) {
	info, err := userInfoCookieValue(user)
	if err != nil {
		return nil, err
	}
	return &http.Cookie{
		Name:     UserCookieName,
		Value:    info,
		Expires:  time.Now().Add(24 * time.Hour),
		Domain:   r.Host,
		Path:     "/",
		HttpOnly: false, // allow client to access this cookie
	}, nil
}

func setCookies(w http.ResponseWriter, r *http.Request, token string, user *types.User) {
	tokenCookie := createTokenCookie(token, r)
	http.SetCookie(w, tokenCookie)
	// set user info to cookie
	userCookie, err := createUserInfoCookie(user.Strip(), r)
	if err != nil {
		logrus.WithContext(r.Context()).Errorf("failed to create user cookie, err: %v", err)
		return
	}
	http.SetCookie(w, userCookie)
}
