package project

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/gotasma/internal/app/types"
	"github.com/gotasma/internal/pkg/http/respond"
)

type (
	service interface {
		Create(ctx context.Context, req *types.CreateProjectRequest) (*types.Project, error)

		//Client use SAVE api, sent all tasks of project, only update tasks has new update_at time ==> Client has to take care of task info strictly
		Save(ctx context.Context, id string, req *types.SaveProject) (*types.ProjectHistory, error)
		Update(ctx context.Context, id string, req *types.UpdateProject) (*types.ProjectHistory, error)
		Delete(ctx context.Context, id string) error

		FindAllProjects(ctx context.Context) ([]*types.ProjectInfo, error)
		FindByID(context.Context, string) (*types.Project, error)

		//TODO: check userID exist
		AddDevs(ctx context.Context, userID []string, projectID string) ([]string, error)
		RemoveDev(ctx context.Context, userID string, projectID string) error

		//Holiday services
		AddHoliday(ctx context.Context, holidayID string, projectID string) (string, error)
		RemoveHoliday(ctx context.Context, holidayID string, projectID string) (string, error)
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

func (h *Handler) Save(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["project_id"]

	if id == "" {
		logrus.Error("Fail to Save Project due to empty project ID ")
		respond.Error(w, errors.New("invalid id"), http.StatusBadRequest)
		return
	}

	var req types.SaveProject

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Errorf("Fail to parse JSON to Save Project Request struct, %v", err)
		respond.Error(w, err, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	//Update project tasks
	info, err := h.srv.Save(r.Context(), id, &req)

	if err != nil {
		logrus.Errorf("Fail to Save Project due to, %v", err)
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}

	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: info,
	})
}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["project_id"]

	if id == "" {
		logrus.Error("Fail to Update Project due to empty project ID ")
		respond.Error(w, errors.New("invalid id"), http.StatusBadRequest)
		return
	}

	var req types.UpdateProject

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Errorf("Fail to parse JSON to Update Project Request struct, %v", err)
		respond.Error(w, err, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	//Update project tasks
	project, err := h.srv.Update(r.Context(), id, &req)

	if err != nil {
		logrus.Errorf("Fail to Save Project due to, %v", err)
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}

	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: project,
	})
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
		logrus.Errorf("Fail to Create Project due to, %v", err)
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
		logrus.Error("Fail to delete project due to empty project ID ")
		respond.Error(w, errors.New("invalid id"), http.StatusBadRequest)
		return
	}
	if err := h.srv.Delete(r.Context(), id); err != nil {
		logrus.Errorf("Fail to delete project due to %v", err)
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: types.IDResponse{
			ID: id,
		},
	})
}

func (h *Handler) FindAll(w http.ResponseWriter, r *http.Request) {

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

func (h *Handler) FindByID(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["project_id"]
	if id == "" {
		logrus.Error("Fail to get project due to empty project ID ")
		respond.Error(w, errors.New("invalid id"), http.StatusBadRequest)
		return
	}

	project, err := h.srv.FindByID(r.Context(), id)

	if err != nil {
		logrus.Errorf("Fail to get project due to, %v", err)
		respond.Error(w, err, http.StatusBadRequest)
		return
	}

	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: project,
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
func (h *Handler) AddHoliday(w http.ResponseWriter, r *http.Request) {
	projectID := mux.Vars(r)["project_id"]

	if projectID == "" {
		logrus.Errorf("Fail to assign project to user due to empty project ID")
		respond.Error(w, fmt.Errorf("Project ID is not valid"), http.StatusBadRequest)
		return
	}

	var req types.AddHolidayRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Errorf("Fail to pasre JSON to Add Users Request struct, %v", err)
		respond.Error(w, err, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	holiday, err := h.srv.AddHoliday(r.Context(), req.UserID, projectID)
	if err != nil {
		logrus.Errorf("Fail to assign holiday to project due to %v", err)
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}

	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: holiday,
	})
}

func (h *Handler) RemoveHoliday(w http.ResponseWriter, r *http.Request) {

	projectID := mux.Vars(r)["project_id"]
	if projectID == "" {
		logrus.Errorf("Fail to remove project to user due to empty project ID")
		respond.Error(w, fmt.Errorf("Project ID is not valid"), http.StatusBadRequest)
		return
	}

	var req types.RemoveHolidayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Errorf("Fail to pasre JSON to Remove Users Request struct, %v", err)
		respond.Error(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	holiday, err := h.srv.RemoveHoliday(r.Context(), req.UserID, projectID)

	if err != nil {
		logrus.Errorf("Fail to remove holiday from project due to %v", err)
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}

	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: holiday,
	})
}
