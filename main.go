package main

import (
	"log"
	"net/http"
	"os"

	"mesa-ayuda/internal/middleware"
	"mesa-ayuda/internal/nutricion"
	"mesa-ayuda/internal/respuesta"
	"mesa-ayuda/internal/sesion"
	"mesa-ayuda/internal/tickets"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Registro)
	r.Use(middleware.Recuperacion)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		respuesta.Error(w, http.StatusNotFound, "ruta_inexistente", "la ruta no existe")
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		respuesta.Error(w, http.StatusMethodNotAllowed, "metodo_no_permitido", "el método no está permitido en esta ruta")
	})

	// Tickets
	almacen := tickets.NuevoAlmacen()
	r.Get("/tickets", almacen.Listar)
	r.Post("/tickets", almacen.Crear)
	r.Get("/tickets/{id}", almacen.Obtener)
	r.Put("/tickets/{id}/estado", almacen.ActualizarEstado)

	// Autenticación
	secreto := os.Getenv("TOKEN_SECRETO")
	if secreto == "" {
		secreto = "mi-clave-super-secreta"
	}

	tokens := sesion.NuevosTokens(secreto)

	manejadorSesion := &sesion.Manejador{
		DB:     nil,
		Tokens: tokens,
	}

	r.Post("/login", manejadorSesion.Login)

	// Rutas protegidas de nutrición
	manejadorNutricion := &nutricion.Manejador{
		DB:     nil,
		Tokens: tokens,
		Planes: map[uint]nutricion.PlanNutricional{
			1: {
				ID:          1,
				IdMascota:   1,
				Estado:      "finalizado",
				Descripcion: "Plan de prueba finalizado",
			},
			2: {
				ID:          2,
				IdMascota:   1,
				Estado:      "pendiente",
				Descripcion: "Plan de prueba pendiente",
			},
		},
	}

	manejadorNutricion.Rutas(r, tokens)

	// Ruta exclusiva para veterinarios
	r.With(
		sesion.Guardia(tokens),
		sesion.RequiereRol(sesion.RolVeterinario),
	).Get("/veterinario", func(w http.ResponseWriter, r *http.Request) {
		respuesta.Exito(w, http.StatusOK, map[string]string{
			"mensaje": "Acceso autorizado para veterinarios",
		})
	})

	log.Println("Servidor mesa de ayuda en :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
