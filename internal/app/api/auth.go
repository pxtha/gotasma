package api

import (
	"github.com/gotasma/internal/app/auth"
	envconfig "github.com/gotasma/internal/pkg/env"
	"github.com/gotasma/internal/pkg/jwt"
)

func newAuthHandler(signer jwt.Signer, authenticator auth.Authenticator) *auth.Handler {
	srv := auth.NewService(signer, authenticator)
	return auth.NewHandler(srv)
}

func newJWTSignVerifier() jwt.SignVerifier {
	var conf jwt.Config
	envconfig.Load(&conf)
	return jwt.New(conf)
}
