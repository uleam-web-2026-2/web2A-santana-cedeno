package main

import (
	"log"
	"net/http"
	"os"

	"mesa-ayuda/internal/middleware"
	"mesa-ayuda/internal/nutricion" // 👈 nuevo import para mascotas
	"mesa-ayuda/internal/respuesta"
	"mesa-ayuda/internal/sesion" // 👈 aquí está la corrección
	"mesa-ayuda/internal/tickets"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	// Middleware en orden
	r.Use(middleware.Registro)
	r.Use(middleware.Recuperacion)

	// Sobrescribir 404 y 405
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		respuesta.Error(w, http.StatusNotFound, "ruta_inexistente", "la ruta no existe")
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		respuesta.Error(w, http.StatusMethodNotAllowed, "metodo_no_permitido",
			"el método no está permitido en esta ruta")
	})

	// Rutas de tickets
	almacen := tickets.NuevoAlmacen()
	r.Get("/tickets", almacen.Listar)
	r.Post("/tickets", almacen.Crear)
	r.Get("/tickets/{id}", almacen.Obtener)
	r.Get("/explotar", tickets.Explotar)
	r.Put("/tickets/{id}/estado", almacen.ActualizarEstado)

	// Nueva ruta de usuarios (Petición 0)
	r.Post("/usuarios", sesion.CrearUsuario)

	// Nueva ruta de mascotas (Petición 1 y 3)
	r.Get("/mascotas", nutricion.ListarMascotas)

	// Nueva ruta de login (Petición 2)
	// 👉 Usa NuevosTokens para inicializar con un secreto
	secreto := os.Getenv("TOKEN_SECRETO")
	if secreto == "" {
		secreto = "mi-clave-super-secreta" // valor por defecto para pruebas
	}
	manejadorSesion := &sesion.Manejador{
		DB:     nil, // luego conectas GORM
		Tokens: sesion.NuevosTokens(secreto),
	}
	r.Post("/login", manejadorSesion.Login)

	log.Println("Servidor mesa de ayuda en :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
