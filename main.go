package main

import (
	"log"
	"net/http"

	"mesa-ayuda/internal/middleware"
	"mesa-ayuda/internal/respuesta"
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

	almacen := tickets.NuevoAlmacen()
	r.Get("/tickets", almacen.Listar)
	r.Post("/tickets", almacen.Crear)
	r.Get("/tickets/{id}", almacen.Obtener)
	r.Get("/explotar", tickets.Explotar)

	log.Println("Servidor mesa de ayuda en :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
