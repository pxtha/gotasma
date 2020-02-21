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
			Path:        "/api/v1/projects/{project_id:[a-z0-9-\\-]+}",
			Method:      http.MethodGet,
			Handler:     h.FindByID,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
		{
			Path:        "/api/v1/projects/{project_id:[a-z0-9-\\-]+}",
			Method:      http.MethodDelete,
			Handler:     h.Delete,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
		//DONE Save project PUT
		{
			Path:        "/api/v1/projects/{project_id:[a-z0-9-\\-]+}",
			Method:      http.MethodPut,
			Handler:     h.Save,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
		//DONE Update project info POST
		{
			Path:        "/api/v1/projects/{project_id:[a-z0-9-\\-]+}",
			Method:      http.MethodPost,
			Handler:     h.Update,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
		//Manage devs of project
		{
			Path:        "/api/v1/projects/{project_id:[a-z0-9-\\-]+}/devs",
			Method:      http.MethodGet,
			Handler:     h.FindAllDevs,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
		{
			Path:        "/api/v1/projects/{project_id:[a-z0-9-\\-]+}/devs",
			Method:      http.MethodPost,
			Handler:     h.AddDev,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
		{
			Path:        "/api/v1/projects/{project_id:[a-z0-9-\\-]+}/devs",
			Method:      http.MethodDelete,
			Handler:     h.RemoveDev,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
		//Manage holidays of project
		{
			Path:        "/api/v1/projects/{project_id:[a-z0-9-\\-]+}/holidays",
			Method:      http.MethodPost,
			Handler:     h.AddHoliday,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
		{
			Path:        "/api/v1/projects/{project_id:[a-z0-9-\\-]+}/holidays",
			Method:      http.MethodDelete,
			Handler:     h.RemoveHoliday,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
		{
			Path:        "/api/v1/projects/{project_id:[a-z0-9-\\-]+}/holidays",
			Method:      http.MethodGet,
			Handler:     h.FindAllHolidays,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},

		//Manage tasks of project
		{
			Path:        "/api/v1/projects/{project_id:[a-z0-9-\\-]+}/tasks",
			Method:      http.MethodGet,
			Handler:     h.FindAllTasks,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
		{
			Path:        "/api/v1/projects/{project_id:[a-z0-9-\\-]+}/tasks",
			Method:      http.MethodPost,
			Handler:     h.AssignDev,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
	}
}
