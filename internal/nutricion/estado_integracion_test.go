package nutricion

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"mesa-ayuda/internal/sesion"

	"github.com/go-chi/chi/v5"
)

func TestCambiarEstadoTransicionNoPermitida(t *testing.T) {
	tokens := sesion.NuevosTokens("clave-de-prueba")

	token, err := tokens.Emitir(sesion.Usuario{
		ID:     1,
		Nombre: "Luis Veterinario",
		Rol:    sesion.RolVeterinario,
	})
	if err != nil {
		t.Fatalf("no se pudo generar el token: %v", err)
	}

	manejador := &Manejador{
		Tokens: tokens,
		Planes: map[uint]PlanNutricional{
			1: {
				ID:     1,
				Estado: "finalizado",
			},
		},
	}

	router := chi.NewRouter()
	manejador.Rutas(router, tokens)

	peticion := httptest.NewRequest(
		http.MethodPatch,
		"/planes/1",
		strings.NewReader(`{"estado":"activo"}`),
	)
	peticion.Header.Set("Authorization", "Bearer "+token)
	peticion.Header.Set("Content-Type", "application/json")

	respuesta := httptest.NewRecorder()
	router.ServeHTTP(respuesta, peticion)

	if respuesta.Code != http.StatusUnprocessableEntity {
		t.Fatalf(
			"se esperaba 422, se recibió %d. Respuesta: %s",
			respuesta.Code,
			respuesta.Body.String(),
		)
	}
}
