package api

import (
	"net/http"

	"praslar.com/gotasma/internal/app/auth"
	"praslar.com/gotasma/internal/pkg/http/middleware"
	"praslar.com/gotasma/internal/pkg/router"
)

const (
	get     = http.MethodGet
	post    = http.MethodPost
	put     = http.MethodPut
	delete  = http.MethodDelete
	options = http.MethodOptions
)

func NewRouter() (http.Handler, error) {

	userSrv, err := newUserService()
	if err != nil {
		return nil, err
	}
	userHandler := newUserHandler(userSrv)

	jwtSignVerifier := newJWTSignVerifier()

	userInfoMiddleware := auth.UserInfoMiddleware(jwtSignVerifier)

	authHandler := newAuthHandler(jwtSignVerifier, userSrv)
	indexHandler := NewIndexHandler()

	routes := []router.Route{
		{
			Path:        "/",
			Method:      get,
			Handler:     indexHandler.ServeHTTP,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
	}

	routes = append(routes, userHandler.Routes()...)
	routes = append(routes, authHandler.Routes()...)

	conf := router.LoadConfigFromEnv()
	conf.Routes = routes
	conf.Middlewares = []router.Middleware{
		userInfoMiddleware,
	}

	r, err := router.New(conf)

	if err != nil {
		return nil, err
	}

	return middleware.CORS(r), nil
}
