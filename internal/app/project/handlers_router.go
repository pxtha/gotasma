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
			Handler:     h.FindAllProjects,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
		{
			Path:        "/api/v1/projects/{project_id:[a-z0-9-\\-]+}",
			Method:      http.MethodDelete,
			Handler:     h.Delete,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
		{
			Path:   "/api/v1/projects/{project_id:[a-z0-9-\\-]+}/devs",
			Method: http.MethodGet,
			//Find all devs of this project
			Handler:     h.FindAllDevs,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
		{
			Path:        "/api/v1/projects/{project_id:[a-z0-9-\\-]+}/devs",
			Method:      http.MethodPut,
			Handler:     h.AddDevs,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
		{
			Path:        "/api/v1/projects/{project_id:[a-z0-9-\\-]+}/devs",
			Method:      http.MethodDelete,
			Handler:     h.RemoveDev,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
	}
}
