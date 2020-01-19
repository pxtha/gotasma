package user

import (
	"net/http"

	"praslar.com/gotasma/internal/pkg/router"
)

func (h *Handler) Routes() []router.Route {
	return []router.Route{
		{
			Path:    "/api/v1/users/registration",
			Method:  http.MethodPost,
			Handler: h.Register,
		},
	}
}
