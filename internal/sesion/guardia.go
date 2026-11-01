package sesion

import (
	"context"
	"net/http"
	"strings"

	"mesa-ayuda/internal/respuesta"
)

type claveSesion struct{}

func SesionDesdeContexto(ctx context.Context) (Sesion, bool) {
	s, ok := ctx.Value(claveSesion{}).(Sesion)
	return s, ok
}

func Guardia(tokens *Tokens) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			partes := strings.SplitN(auth, " ", 2)

			if len(partes) != 2 || !strings.EqualFold(partes[0], "Bearer") {
				respuesta.Error(w, http.StatusUnauthorized, "no_autorizado", "token requerido o inválido")
				return
			}

			sesion, err := tokens.Leer(partes[1])
			if err != nil {
				respuesta.Error(w, http.StatusUnauthorized, "token_invalido", "el token no es válido o ha vencido")
				return
			}

			ctx := context.WithValue(r.Context(), claveSesion{}, sesion)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequiereRol(rol string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s, ok := SesionDesdeContexto(r.Context())
			if !ok {
				respuesta.Error(w, http.StatusUnauthorized, "no_autorizado", "sesión requerida")
				return
			}

			if s.Rol != rol {
				respuesta.Error(w, http.StatusForbidden, "rol_no_autorizado", "no tienes permisos para esta operación")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
