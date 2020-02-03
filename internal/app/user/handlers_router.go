package user

import (
	"net/http"

	"praslar.com/gotasma/internal/app/auth"
	"praslar.com/gotasma/internal/pkg/http/router"
)

func (h *Handler) Routes() []router.Route {
	return []router.Route{
		{
			Path:    "/api/v1/users/registration",
			Method:  http.MethodPost,
			Handler: h.Register,
		},
		{
			Path:        "/api/v1/dev",
			Method:      http.MethodPost,
			Handler:     h.CreateDev,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
		{
			Path:        "/api/v1/dev",
			Method:      http.MethodGet,
			Handler:     h.FindAllDev,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
		{
			Path:        "/api/v1/dev/{user_id:[a-z0-9-\\-]+}",
			Method:      http.MethodDelete,
			Handler:     h.Delete,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
	}
}
