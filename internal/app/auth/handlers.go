package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gotasma/internal/app/types"
	"github.com/gotasma/internal/pkg/http/respond"
)

type (
	service interface {
		Auth(ctx context.Context, username, password string) (string, error)
	}
	Handler struct {
		srv service
	}
)

func NewHandler(srv service) *Handler {
	return &Handler{
		srv: srv,
	}
}

func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Email    string
		Password string
	}{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	token, err := h.srv.Auth(r.Context(), req.Email, req.Password)
	if err != nil {
		respond.Error(w, err, http.StatusUnauthorized)
		return
	}

	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: map[string]interface{}{
			"token": token,
		},
	})
}
