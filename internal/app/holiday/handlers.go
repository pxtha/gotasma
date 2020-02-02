package holiday

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"praslar.com/gotasma/internal/app/types"
	"praslar.com/gotasma/internal/pkg/http/respond"
)

type (
	service interface {
		//TODO validate end > start
		Create(ctx context.Context, req *types.HolidayRequest) (*types.Holiday, error)
		//TODO check holiday creater id before delete
		Delete(ctx context.Context, id string) error
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
		respond.Error(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	holiday, err := h.srv.Create(r.Context(), &req)
	if err != nil {
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
		respond.Error(w, errors.New("invalid id"), http.StatusBadRequest)
		return
	}
	if err := h.srv.Delete(r.Context(), id); err != nil {
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
	holiday, err := h.srv.FindAll(r.Context())
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: holiday,
	})
}
