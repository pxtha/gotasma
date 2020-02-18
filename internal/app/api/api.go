package api

import (
	"net/http"

	"github.com/gotasma/internal/app/auth"
	"github.com/gotasma/internal/pkg/http/middleware"
	"github.com/gotasma/internal/pkg/http/router"
)

const (
	get     = http.MethodGet
	post    = http.MethodPost
	put     = http.MethodPut
	delete  = http.MethodDelete
	options = http.MethodOptions
)

func NewRouter() (http.Handler, error) {
	//=================Policy-Role Base=========
	policySrv, err := newPolicyService()
	if err != nil {
		return nil, err
	}
	//=================User=====================
	userSrv, err := newUserService(policySrv)
	if err != nil {
		return nil, err
	}
	userHandler := newUserHandler(userSrv)
	//===============Holiday====================
	holidaySrv, err := newHolidayService(policySrv)
	if err != nil {
		return nil, err
	}
	holidayHandler := newHolidayHandler(holidaySrv)
	//===============Project====================
	projectSrv, err := newProjectService(policySrv, userSrv, holidaySrv)
	if err != nil {
		return nil, err
	}
	projectHandler := newProjectHandler(projectSrv)
	//===============Sub handler================
	jwtSignVerifier := newJWTSignVerifier()
	userInfoMiddleware := auth.UserInfoMiddleware(jwtSignVerifier)
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
	routes = append(routes, holidayHandler.Routes()...)
	routes = append(routes, projectHandler.Routes()...)

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
