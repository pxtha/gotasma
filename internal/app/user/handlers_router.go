package user

import (
	"net/http"

	"praslar.com/gotasma/internal/app/auth"
	"praslar.com/gotasma/internal/pkg/router"
)

func (h *Handler) Routes() []router.Route {
	return []router.Route{
		{
			Path:    "/api/v1/users/registration",
			Method:  http.MethodPost,
			Handler: h.Register,
		},
		{
			Path:        "/api/v1/users/dev",
			Method:      http.MethodPost,
			Handler:     h.CreateDev,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
		{
			Path:        "/api/v1/users/dev",
			Method:      http.MethodGet,
			Handler:     h.FindAllDev,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
		{
			Path:        "/api/v1/users/{id:[a-z0-9-\\-]+}",
			Method:      http.MethodDelete,
			Handler:     h.Delete,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
	}
}
