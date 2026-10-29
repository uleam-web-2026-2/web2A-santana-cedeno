package nutricion

import (
	"mesa-ayuda/internal/respuesta"
	"net/http"
)

// ListarMascotas maneja GET /mascotas
func ListarMascotas(w http.ResponseWriter, r *http.Request) {
	// Verificar si viene el header Authorization
	token := r.Header.Get("Authorization")
	if token == "" {
		respuesta.Error(w, http.StatusUnauthorized, "no_autorizado", "token requerido")
		return
	}

	// Si hay token, aquí luego devolverás la lista real (Petición 3)
	respuesta.Exito(w, http.StatusOK, []map[string]string{
		{"id": "1", "nombre": "Firulais", "tipo": "perro"},
		{"id": "2", "nombre": "Mishi", "tipo": "gato"},
	})
}
