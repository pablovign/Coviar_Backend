package handler

import (
	"net/http"
	"strconv"

	"coviar_backend/internal/service"
	"coviar_backend/pkg/httputil"
)

type AdminHandler struct {
	service *service.AdminService
}

func NewAdminHandler(service *service.AdminService) *AdminHandler {
	return &AdminHandler{service: service}
}

// GetStats GET /api/admin/stats
func (h *AdminHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetStats(r.Context())
	if err != nil {
		httputil.HandleServiceError(w, err)
		return
	}

	httputil.RespondJSON(w, http.StatusOK, stats)
}

// GetAllEvaluaciones GET /api/admin/evaluaciones
func (h *AdminHandler) GetAllEvaluaciones(w http.ResponseWriter, r *http.Request) {
	// Parámetros de query opcionales
	estado := r.URL.Query().Get("estado")
	idBodegaStr := r.URL.Query().Get("id_bodega")

	var idBodega int
	if idBodegaStr != "" {
		var err error
		idBodega, err = strconv.Atoi(idBodegaStr)
		if err != nil {
			httputil.RespondError(w, http.StatusBadRequest, "id_bodega inválido")
			return
		}
	}

	evaluaciones, err := h.service.GetAllEvaluaciones(r.Context(), estado, idBodega)
	if err != nil {
		httputil.HandleServiceError(w, err)
		return
	}

	httputil.RespondJSON(w, http.StatusOK, evaluaciones)
}
