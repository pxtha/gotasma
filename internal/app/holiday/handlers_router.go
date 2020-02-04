package holiday

import (
	"net/http"

	"github.com/gotasma/internal/app/auth"
	"github.com/gotasma/internal/pkg/http/router"
)

func (h *Handler) Routes() []router.Route {
	return []router.Route{
		{
			Path:        "/api/v1/holidays",
			Method:      http.MethodPost,
			Handler:     h.Create,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
		{
			Path:        "/api/v1/holidays/{holiday_id:[a-z0-9-\\-]+}",
			Method:      http.MethodDelete,
			Handler:     h.Delete,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
		{
			Path:        "/api/v1/holidays",
			Method:      http.MethodGet,
			Handler:     h.FindAll,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
	}
}
