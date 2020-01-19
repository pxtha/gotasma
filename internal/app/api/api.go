package api

import (
	"net/http"

	"github.com/sirupsen/logrus"

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

	authHandler := newAuthHandler(jwtSignVerifier, userSrv)
	indexHandler := NewIndexHandler()

	routes := []router.Route{
		{
			Path:    "/",
			Method:  get,
			Handler: indexHandler.ServeHTTP,
		},
	}

	routes = append(routes, userHandler.Routes()...)
	routes = append(routes, authHandler.Routes()...)

	conf := router.LoadConfigFromEnv()
	conf.Routes = routes
	//conf.NotFoundHandler = indexHandler

	logrus.Info(conf)
	r, err := router.New(conf)
	if err != nil {
		return nil, err
	}

	return r, nil
}
