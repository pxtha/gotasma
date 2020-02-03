package project

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gotasma/internal/app/types"
	"github.com/gotasma/internal/pkg/http/respond"
)

type (
	service interface {
		Create(ctx context.Context, req *types.CreateProjectRequest) (*types.Project, error)
		FindAll(ctx context.Context) ([]*types.Project, error)
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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req types.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	project, err := h.srv.Create(r.Context(), &req)
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: project,
	})
}

func (h *Handler) FindAll(w http.ResponseWriter, r *http.Request) {
	projects, err := h.srv.FindAll(r.Context())

	if err != nil {
		respond.Error(w, err, http.StatusBadRequest)
		return
	}

	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: projects,
	})
}
