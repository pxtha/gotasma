package project

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gotasma/internal/app/types"
	"github.com/gotasma/internal/pkg/http/respond"
	
	"github.com/sirupsen/logrus"
	"github.com/gorilla/mux"

)

type (
	service interface {
		Create(ctx context.Context, req *types.CreateProjectRequest) (*types.Project, error)
		Save(ctx context.Context, id string, req *types.SaveProject) ([]*types.Task, error)
		Update(ctx context.Context, id string, req *types.UpdateProject) (*types.UpdateProject, error)
		Delete(ctx context.Context, id string) error

		FindByID(context.Context, string) (*types.Project, error)

		FindAllProjects(ctx context.Context) ([]*types.ProjectInfo, error)

		//User services
		AddDev(ctx context.Context, req *types.AddUsersRequest, projectID string) (string, error)
		RemoveDev(ctx context.Context, req *types.RemoveUserRequest, projectID string) (string, error)
		FindAllDevs(context.Context, string) ([]*types.UserInfo, error)

		//Holiday services
		AddHoliday(ctx context.Context, req *types.AddHolidayRequest, projectID string) (string, error)
		RemoveHoliday(ctx context.Context, req *types.RemoveHolidayRequest, projectID string) (string, error)
		FindAllHolidays(context.Context, string) ([]*types.HolidayInfo, error)

		//Tasks service
		FindAllTasks(context.Context, string) ([]*types.Task, error)
		AssignDev(ctx context.Context, projectID string, req *types.AssignDev) (*types.WorkLoadInfo, error)
		UnAssignDev(ctx context.Context, projectID string, req *types.UnAssignDev) (*types.WorkLoadInfo, error)
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

//Manage devs of project
func (h *Handler) AddDev(w http.ResponseWriter, r *http.Request) {
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

	assignedUser, err := h.srv.AddDev(r.Context(), &req, projectID)
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

	user, err := h.srv.RemoveDev(r.Context(), &req, projectID)
	if err != nil {
		logrus.Errorf("Fail to remove user from project due to %v", err)
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}

	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: user,
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

//Manage holidays of project
func (h *Handler) AddHoliday(w http.ResponseWriter, r *http.Request) {
	projectID := mux.Vars(r)["project_id"]

	if projectID == "" {
		logrus.Errorf("Fail to assign project to user due to empty project ID")
		respond.Error(w, fmt.Errorf("Project ID is not valid"), http.StatusBadRequest)
		return
	}

	var req types.AddHolidayRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Errorf("Fail to pasre JSON to Add holiday Request struct, %v", err)
		respond.Error(w, err, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	holiday, err := h.srv.AddHoliday(r.Context(), &req, projectID)
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
		logrus.Errorf("Fail to remove holiday from user due to empty project ID")
		respond.Error(w, fmt.Errorf("Project ID is not valid"), http.StatusBadRequest)
		return
	}

	var req types.RemoveHolidayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Errorf("Fail to pasre JSON to Remove holiday Request struct, %v", err)
		respond.Error(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	holiday, err := h.srv.RemoveHoliday(r.Context(), &req, projectID)
	if err != nil {
		logrus.Errorf("Fail to remove holiday from project due to %v", err)
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: holiday,
	})
}

func (h *Handler) FindAllHolidays(w http.ResponseWriter, r *http.Request) {
	projectID := mux.Vars(r)["project_id"]
	if projectID == "" {
		logrus.Errorf("Fail to get holidays of project due to empty project ID")
		respond.Error(w, fmt.Errorf("Project ID is not valid"), http.StatusBadRequest)
		return
	}

	holidays, err := h.srv.FindAllHolidays(r.Context(), projectID)

	if err != nil {
		logrus.Errorf("Fail to holidays devs from project due to %v", err)
		respond.Error(w, err, http.StatusBadRequest)
		return
	}

	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: holidays,
	})
}

//Manage tasks of project
func (h *Handler) FindAllTasks(w http.ResponseWriter, r *http.Request) {
	projectID := mux.Vars(r)["project_id"]
	if projectID == "" {
		logrus.Errorf("Fail to get tasks of project due to empty project ID")
		respond.Error(w, fmt.Errorf("Project ID is not valid"), http.StatusBadRequest)
		return
	}

	tasks, err := h.srv.FindAllTasks(r.Context(), projectID)

	if err != nil {
		logrus.Errorf("Fail to get tasks from project due to %v", err)
		respond.Error(w, err, http.StatusBadRequest)
		return
	}

	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: tasks,
	})
}

func (h *Handler) AssignDev(w http.ResponseWriter, r *http.Request) {
	projectID := mux.Vars(r)["project_id"]

	if projectID == "" {
		logrus.Errorf("Fail to assign task to user due to empty project ID")
		respond.Error(w, fmt.Errorf("Project ID is not valid"), http.StatusBadRequest)
		return
	}

	var req types.AssignDev

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Errorf("Fail to pasre JSON to assign task Request struct, %v", err)
		respond.Error(w, err, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	workload, err := h.srv.AssignDev(r.Context(), projectID, &req)
	if err != nil {
		logrus.Errorf("Fail to assign task to user due to %v", err)
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}

	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: workload,
	})
}

func (h *Handler) UnAssignDev(w http.ResponseWriter, r *http.Request) {
	projectID := mux.Vars(r)["project_id"]

	if projectID == "" {
		logrus.Errorf("Fail to assign task to user due to empty project ID")
		respond.Error(w, fmt.Errorf("Project ID is not valid"), http.StatusBadRequest)
		return
	}

	var req types.UnAssignDev

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Errorf("Fail to pasre JSON to assign task Request struct, %v", err)
		respond.Error(w, err, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	task, err := h.srv.UnAssignDev(r.Context(), projectID, &req)
	if err != nil {
		logrus.Errorf("Fail to unassign task to user due to %v", err)
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}

	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: task,
	})
}
