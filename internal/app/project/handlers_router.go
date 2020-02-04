package project

import (
	"net/http"

	"github.com/gotasma/internal/app/auth"
	"github.com/gotasma/internal/pkg/http/router"
)

func (h *Handler) Routes() []router.Route {
	return []router.Route{
		{
			Path:        "/api/v1/projects",
			Method:      http.MethodPost,
			Handler:     h.Create,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
		{
			Path:        "/api/v1/projects",
			Method:      http.MethodGet,
			Handler:     h.FindAll,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
		{
			Path:        "/api/v1/projects/{project_id:[a-z0-9-\\-]+}/devs",
			Method:      http.MethodPut,
			Handler:     h.AddProjectToUser,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
		{
			Path:        "/api/v1/projects/{project_id:[a-z0-9-\\-]+}/devs",
			Method:      http.MethodDelete,
			Handler:     h.RemoveUserFromProject,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
	}
}
