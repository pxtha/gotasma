package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/gotasma/internal/app/types"
	"github.com/gotasma/internal/pkg/http/respond"
)

type (
	service interface {
		Register(ctx context.Context, req *types.RegisterRequest) (*types.User, error)
		// Find users that role is Dev, also show PM info of those devs
		FindAllDev(ctx context.Context) ([]*types.UserInfo, error)
		CreateDev(ctx context.Context, req *types.RegisterRequest) (*types.User, error)
		Delete(ctx context.Context, id string) error
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

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req types.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Errorf("Fail to parse JSON to Register PM Request struct, %w", err)
		respond.Error(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	//Allways register as Project Manager
	req.Role = types.PM
	user, err := h.srv.Register(r.Context(), &req)
	if err != nil {
		logrus.Errorf("Fail to Register new PM account, %w", err)
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: user,
	})
}

func (h *Handler) CreateDev(w http.ResponseWriter, r *http.Request) {
	var req types.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Errorf("Fail to parse JSON to Register DEV Request struct, %w", err)
		respond.Error(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	//Allways create as DEV
	req.Role = types.DEV
	user, err := h.srv.CreateDev(r.Context(), &req)
	if err != nil {
		logrus.Errorf("Fail to Register new DEV account, %w", err)
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: user,
	})
}

func (h *Handler) FindAllDev(w http.ResponseWriter, r *http.Request) {
	users, err := h.srv.FindAllDev(r.Context())
	if err != nil {
		logrus.Errorf("Fail to find DEVs of project, %w", err)
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: users,
	})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["user_id"]
	if id == "" {
		logrus.Error("Fail to remove DEV from project due to empty ID")
		respond.Error(w, errors.New("invalid id"), http.StatusBadRequest)
		return
	}

	if err := h.srv.Delete(r.Context(), id); err != nil {
		logrus.Errorf("Fail to remove DEV from project, %v", err)
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}

	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: types.IDResponse{
			ID: id,
		},
	})
}
