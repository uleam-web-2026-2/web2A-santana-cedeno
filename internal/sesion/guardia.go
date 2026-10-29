package sesion

import (
	"net/http"
)

// Guardia es un middleware que valida el token en las rutas protegidas.
func Guardia(tokens *Tokens) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Aquí deberías validar el token del encabezado Authorization.
			// Por ahora dejamos un "passthrough" para que compile.
			// Luego puedes implementar la lógica real.
			auth := r.Header.Get("Authorization")
			if auth == "" {
				http.Error(w, "token requerido", http.StatusUnauthorized)
				return
			}
			// Si hay token, continúa con la ruta
			next.ServeHTTP(w, r)
		})
	}
}
