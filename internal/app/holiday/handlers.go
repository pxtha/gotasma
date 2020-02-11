package holiday

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
		Create(ctx context.Context, req *types.HolidayRequest) (*types.Holiday, error)
		Delete(ctx context.Context, id string) error
		Update(ctx context.Context, id string, req *types.HolidayRequest) (*types.HolidayInfo, error)
		FindAll(ctx context.Context) ([]*types.Holiday, error)
	}
	Handler struct {
		srv service
	}
)

func NewHanlder(srv service) *Handler {
	return &Handler{
		srv: srv,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req types.HolidayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Errorf("Fail to parse JSON into struct Holiday Request, %v", err)
		respond.Error(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	holiday, err := h.srv.Create(r.Context(), &req)
	if err != nil {
		logrus.Errorf("Fail to create new Holiday, %v", err)
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: holiday,
	})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["holiday_id"]
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

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["holiday_id"]
	if id == "" {
		logrus.Error("Fail to Update holiday due to empty holiday ID ")
		respond.Error(w, errors.New("invalid id"), http.StatusBadRequest)
		return
	}

	var req types.HolidayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Errorf("Fail to parse JSON into struct Holiday Request, %v", err)
		respond.Error(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	holiday, err := h.srv.Update(r.Context(), id, &req)
	if err != nil {
		logrus.Errorf("Fail to update holiday due to %v", err)
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: holiday,
	})
}

func (h *Handler) FindAll(w http.ResponseWriter, r *http.Request) {
	holiday, err := h.srv.FindAll(r.Context())
	if err != nil {
		logrus.Error("Fail to get all holiday due to %v", err)
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: holiday,
	})
}
