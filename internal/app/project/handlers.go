package project

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gotasma/internal/app/types"
	"github.com/gotasma/internal/pkg/http/respond"
	"github.com/sirupsen/logrus"
)

type (
	service interface {
		Create(ctx context.Context, req *types.CreateProjectRequest) (*types.Project, error)
		Delete(ctx context.Context, id string) error
		FindAllProjects(ctx context.Context) ([]*types.Project, error)
		FindAllDevs(context.Context, string) ([]*types.UserInfo, error)
		AddDevs(ctx context.Context, userID []string, projectID string) ([]string, error)
		RemoveDev(ctx context.Context, userID string, projectID string) error
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
		logrus.Errorf("Fail to parse JSON to Create Project Request struct, %v", err)
		respond.Error(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	project, err := h.srv.Create(r.Context(), &req)
	if err != nil {
		logrus.Errorf("Fail to parse Create Project due to, %v", err)
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: project,
	})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["project_id"]
	if id == "" {
		logrus.Error("Fail to delete holiday due to empty holiday ID ")
		respond.Error(w, errors.New("invalid id"), http.StatusBadRequest)
		return
	}
	if err := h.srv.Delete(r.Context(), id); err != nil {
		logrus.Errorf("Fail to delete holiday due to %v", err)
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: types.IDResponse{
			ID: id,
		},
	})
}

//Find all project
func (h *Handler) FindAllProjects(w http.ResponseWriter, r *http.Request) {

	projects, err := h.srv.FindAllProjects(r.Context())

	if err != nil {
		logrus.Errorf("Fail to get all project due to, %v", err)
		respond.Error(w, err, http.StatusBadRequest)
		return
	}

	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: projects,
	})
}

func (h *Handler) AddDevs(w http.ResponseWriter, r *http.Request) {
	projectID := mux.Vars(r)["project_id"]

	if projectID == "" {
		logrus.Errorf("Fail to assign project to user due to empty project ID")
		respond.Error(w, fmt.Errorf("Project ID is not valid"), http.StatusBadRequest)
		return
	}

	var req types.AddUsersRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Errorf("Fail to pasre JSON to Add Users Request struct, %v", err)
		respond.Error(w, err, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	assignedUser, err := h.srv.AddDevs(r.Context(), req.UserIDs, projectID)
	if err != nil {
		logrus.Errorf("Fail to assign project to user due to %v: user: %v", err, assignedUser)
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}

	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: assignedUser,
	})
}

func (h *Handler) RemoveDev(w http.ResponseWriter, r *http.Request) {

	projectID := mux.Vars(r)["project_id"]
	if projectID == "" {
		logrus.Errorf("Fail to remove project to user due to empty project ID")
		respond.Error(w, fmt.Errorf("Project ID is not valid"), http.StatusBadRequest)
		return
	}

	var req types.RemoveUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Errorf("Fail to pasre JSON to Remove Users Request struct, %v", err)
		respond.Error(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := h.srv.RemoveDev(r.Context(), req.UserID, projectID); err != nil {
		logrus.Errorf("Fail to remove user from project due to %v", err)
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}

	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: req,
	})
}

func (h *Handler) FindAllDevs(w http.ResponseWriter, r *http.Request) {
	projectID := mux.Vars(r)["project_id"]
	if projectID == "" {
		logrus.Errorf("Fail to get devs of project due to empty project ID")
		respond.Error(w, fmt.Errorf("Project ID is not valid"), http.StatusBadRequest)
		return
	}

	devs, err := h.srv.FindAllDevs(r.Context(), projectID)

	if err != nil {
		logrus.Errorf("Fail to get devs from project due to %v", err)
		respond.Error(w, err, http.StatusBadRequest)
		return
	}

	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: devs,
	})
}
