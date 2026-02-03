package handler

import (
	"net/http"
	"strconv"

	"coviar_backend/internal/domain"
	"coviar_backend/internal/middleware"
	"coviar_backend/internal/service"
	"coviar_backend/pkg/httputil"
	"coviar_backend/pkg/router"
)

type ResponsableHandler struct {
	service *service.ResponsableService
}

func NewResponsableHandler(service *service.ResponsableService) *ResponsableHandler {
	return &ResponsableHandler{service: service}
}

func (h *ResponsableHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := router.GetParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	// Verificar permisos
	userID := r.Context().Value(middleware.UserIDKey).(int)
	userTipo := r.Context().Value(middleware.UserTipoKey).(string)

	canAccess, err := h.service.CanAccess(r.Context(), id, userID, userTipo)
	if err != nil {
		httputil.HandleServiceError(w, err)
		return
	}
	if !canAccess {
		httputil.RespondError(w, http.StatusForbidden, "no tiene permisos para acceder a este recurso")
		return
	}

	responsable, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		httputil.HandleServiceError(w, err)
		return
	}

	httputil.RespondJSON(w, http.StatusOK, responsable)
}

func (h *ResponsableHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := router.GetParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	// Verificar permisos
	userID := r.Context().Value(middleware.UserIDKey).(int)
	userTipo := r.Context().Value(middleware.UserTipoKey).(string)

	canAccess, err := h.service.CanAccess(r.Context(), id, userID, userTipo)
	if err != nil {
		httputil.HandleServiceError(w, err)
		return
	}
	if !canAccess {
		httputil.RespondError(w, http.StatusForbidden, "no tiene permisos para modificar este recurso")
		return
	}

	var req domain.ResponsableUpdateDTO
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := h.service.Update(r.Context(), id, &req); err != nil {
		httputil.HandleServiceError(w, err)
		return
	}

	httputil.RespondJSON(w, http.StatusOK, map[string]string{"mensaje": "Responsable actualizado"})
}

func (h *ResponsableHandler) DarDeBaja(w http.ResponseWriter, r *http.Request) {
	idStr := router.GetParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	// Verificar permisos
	userID := r.Context().Value(middleware.UserIDKey).(int)
	userTipo := r.Context().Value(middleware.UserTipoKey).(string)

	canAccess, err := h.service.CanAccess(r.Context(), id, userID, userTipo)
	if err != nil {
		httputil.HandleServiceError(w, err)
		return
	}
	if !canAccess {
		httputil.RespondError(w, http.StatusForbidden, "no tiene permisos para dar de baja este recurso")
		return
	}

	if err := h.service.DarDeBaja(r.Context(), id); err != nil {
		httputil.HandleServiceError(w, err)
		return
	}

	httputil.RespondJSON(w, http.StatusOK, map[string]string{"mensaje": "Responsable dado de baja exitosamente"})
}
