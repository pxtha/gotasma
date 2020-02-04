package project

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gotasma/internal/app/types"
	"github.com/gotasma/internal/pkg/http/respond"
)

type (
	service interface {
		Create(ctx context.Context, req *types.CreateProjectRequest) (*types.Project, error)
		FindAll(ctx context.Context) ([]*types.Project, error)
		AddProjectToUser(ctx context.Context, userID []string, projectID string) (int, error)
		RemoveUserFromProject(ctx context.Context, userID string, projectID string) error
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

func (h *Handler) AddProjectToUser(w http.ResponseWriter, r *http.Request) {
	projectID := mux.Vars(r)["project_id"]
	if projectID == "" {
		respond.Error(w, fmt.Errorf("Project ID is not valid"), http.StatusBadRequest)
		return
	}

	var req types.AddUsersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	assignedUser, err := h.srv.AddProjectToUser(r.Context(), req.UserIDs, projectID)
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: assignedUser,
	})
}

func (h *Handler) RemoveUserFromProject(w http.ResponseWriter, r *http.Request) {
	projectID := mux.Vars(r)["project_id"]
	if projectID == "" {
		respond.Error(w, fmt.Errorf("Project ID is not valid"), http.StatusBadRequest)
		return
	}
	var req types.RemoveUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := h.srv.RemoveUserFromProject(r.Context(), req.UserID, projectID); err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}

	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: types.BaseResponse{
			Data: req,
		},
	})
}
