package nutricion

import (
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"mesa-ayuda/internal/respuesta"
)

// Estados permitidos y sus transiciones
var siguientes = map[string][]string{
	"pendiente":  {"activo", "cancelado"},
	"activo":     {"finalizado", "cancelado"},
	"finalizado": {},
	"cancelado":  {},
}

// Función auxiliar para validar transiciones
func puedePasar(de, a string) bool {
	return slices.Contains(siguientes[de], a)
}

// Cambiar estado de un plan nutricional
func (m *Manejador) cambiarEstado(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		respuesta.Error(w, http.StatusBadRequest, "el id debe ser un número entero positivo", "")
		return
	}

	var pedido struct {
		Estado string `json:"estado"`
	}
	if err := json.NewDecoder(r.Body).Decode(&pedido); err != nil {
		respuesta.Error(w, http.StatusBadRequest, "el cuerpo no es JSON válido", err.Error())
		return
	}

	if _, existe := siguientes[pedido.Estado]; !existe {
		respuesta.Error(w, http.StatusUnprocessableEntity, "estado no permitido: "+pedido.Estado, "")
		return
	}

	var t PlanNutricional
	err = m.DB.First(&t, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		respuesta.Error(w, http.StatusNotFound, "no existe el plan nutricional", "")
		return
	}
	if err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "no se pudo leer", err.Error())
		return
	}

	if !puedePasar(t.Estado, pedido.Estado) {
		respuesta.Error(w, http.StatusUnprocessableEntity, "de "+t.Estado+" no se puede pasar a "+pedido.Estado, "")
		return
	}

	t.Estado = pedido.Estado
	if err := m.DB.Save(&t).Error; err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "no se pudo guardar", err.Error())
		return
	}

	respuesta.Exito(w, http.StatusOK, t)
}
