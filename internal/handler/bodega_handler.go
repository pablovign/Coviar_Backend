package handler

import (
	"net/http"
	"strconv"

	"coviar_backend/internal/domain"
	"coviar_backend/internal/service"
	"coviar_backend/pkg/httputil"
	"coviar_backend/pkg/router"
)

type BodegaHandler struct {
	service *service.BodegaService
}

func NewBodegaHandler(service *service.BodegaService) *BodegaHandler {
	return &BodegaHandler{service: service}
}

func (h *BodegaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := router.GetParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	bodega, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		httputil.HandleServiceError(w, err)
		return
	}

	httputil.RespondJSON(w, http.StatusOK, bodega)
}

func (h *BodegaHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := router.GetParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var req domain.BodegaUpdateDTO
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := h.service.Update(r.Context(), id, &req); err != nil {
		httputil.HandleServiceError(w, err)
		return
	}

	httputil.RespondJSON(w, http.StatusOK, map[string]string{"mensaje": "Bodega actualizada"})
}
