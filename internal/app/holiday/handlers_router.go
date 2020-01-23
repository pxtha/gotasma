package holiday

import (
	"net/http"

	"praslar.com/gotasma/internal/app/auth"
	"praslar.com/gotasma/internal/pkg/router"
)

func (h *Handler) Routes() []router.Route {
	return []router.Route{
		{
			Path:        "/api/v1/holiday",
			Method:      http.MethodPost,
			Handler:     h.Create,
			Middlewares: []router.Middleware{auth.RequiredAuthMiddleware},
		},
	}
}
