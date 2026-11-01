package nutricion

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"mesa-ayuda/internal/respuesta"
)

func (m *Manejador) cambiarEstado(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id == 0 {
		respuesta.Error(w, http.StatusBadRequest, "id_invalido", "el id debe ser un número entero positivo")
		return
	}

	var pedido struct {
		Estado string `json:"estado"`
	}

	if err := json.NewDecoder(r.Body).Decode(&pedido); err != nil {
		respuesta.Error(w, http.StatusBadRequest, "json_invalido", "el cuerpo no es JSON válido")
		return
	}

	plan, err := m.Servicio.CambiarEstado(uint(id), pedido.Estado)

	switch {
	case errors.Is(err, ErrEstadoInvalido):
		respuesta.Error(w, http.StatusUnprocessableEntity, "estado_invalido", "el estado solicitado no está permitido")
		return

	case errors.Is(err, ErrPlanNoEncontrado):
		respuesta.Error(w, http.StatusNotFound, "plan_no_encontrado", "no existe el plan nutricional")
		return

	case errors.Is(err, ErrTransicionNoPermitida):
		respuesta.Error(w, http.StatusUnprocessableEntity, "transicion_no_permitida", err.Error())
		return

	case err != nil:
		respuesta.Error(w, http.StatusInternalServerError, "error_guardado", "no se pudo actualizar el plan")
		return
	}

	respuesta.Exito(w, http.StatusOK, plan)
}
