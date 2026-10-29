package nutricion

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSinTokenLasRutasProtegidasResponden401(t *testing.T) {
	protegidas := []struct{ metodo, ruta, cuerpo string }{
		{"GET", "/mascotas", ""},
		{"POST", "/mascotas", `{"nombre":"Firulais","especie":"perro"}`},
		{"GET", "/planes/1", ""},
		{"POST", "/planes", `{"descripcion":"Dieta baja en grasa","id_mascota":1}`},
		{"POST", "/detalles", `{"tipo_alimento":"pollo","cantidad":"200g","frecuencia":"2 veces al día"}`},
		{"PATCH", "/planes/1", `{"estado":"activo"}`},
		{"POST", "/citas", `{"id_usuario":1,"id_mascota":1,"fecha":"2026-10-20","motivo":"control veterinario"}`},
	}

	servidor := enrutadorCon("")

	for _, p := range protegidas {
		peticion := httptest.NewRequest(p.metodo, p.ruta, strings.NewReader(p.cuerpo))
		grabadora := httptest.NewRecorder()

		cayo := seCayo(t, func() { servidor.ServeHTTP(grabadora, peticion) })

		if cayo {
			t.Errorf("%s %s sin token: el guardia no la detuvo y la petición llegó hasta la base de datos", p.metodo, p.ruta)
			continue
		}
		if grabadora.Code != http.StatusUnauthorized {
			t.Errorf("%s %s sin token: se esperaba 401 y llegó %d", p.metodo, p.ruta, grabadora.Code)
		}
	}
}
