package holiday

import (
	"context"
	"encoding/json"
	"net/http"

	"praslar.com/gotasma/internal/app/types"
	"praslar.com/gotasma/internal/pkg/http/respond"
)

type (
	service interface {
		Create(ctx context.Context, req *types.HolidayRequest) (*types.Holiday, error)
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
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: holiday,
	})
}
