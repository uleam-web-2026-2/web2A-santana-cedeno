package nutricion

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"

	"mesa-ayuda/internal/sesion"
)

const rolDePrueba = sesion.RolVeterinario

var tokensDePrueba = sesion.NuevosTokens("secreto-solo-de-las-pruebas")

func enrutadorCon(rol string) http.Handler {
	r := chi.NewRouter()
	(&Manejador{DB: nil}).Rutas(r, tokensDePrueba)
	if rol == "" {
		return r
	}
	token, err := tokensDePrueba.Emitir(
		sesion.Usuario{ID: 1, Nombre: "prueba", Rol: rol})
	if err != nil {
		panic(err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, p *http.Request) {
		p.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(w, p)
	})
}

func seCayo(t *testing.T, f func()) (cayo bool) {
	t.Helper()
	defer func() {
		if recover() != nil {
			cayo = true
		}
	}()
	f()
	return false
}
