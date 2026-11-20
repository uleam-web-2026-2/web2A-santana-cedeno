package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"mesa-ayuda/internal/middleware"
	"mesa-ayuda/internal/nutricion"
	"mesa-ayuda/internal/respuesta"
	"mesa-ayuda/internal/sesion"
	"mesa-ayuda/internal/tickets"
)

func main() {
	r := chi.NewRouter()

	// CORS arriba de todo.
	// Para el taller se permite únicamente el origen de desarrollo indicado.
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.Registro)
	r.Use(middleware.Recuperacion)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		respuesta.Error(
			w,
			http.StatusNotFound,
			"ruta_inexistente",
			"la ruta no existe",
		)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		respuesta.Error(
			w,
			http.StatusMethodNotAllowed,
			"metodo_no_permitido",
			"el método no está permitido en esta ruta",
		)
	})

	// ---------------------------------------------------------
	// Tickets: repositorio en memoria
	// ---------------------------------------------------------
	repoTickets := tickets.NuevoRepositorioMemoria()
	servicioTickets := tickets.NuevoServicio(repoTickets)

	r.Get("/tickets", servicioTickets.Listar)
	r.Post("/tickets", servicioTickets.Crear)
	r.Get("/tickets/{id}", servicioTickets.Obtener)
	r.Put("/tickets/{id}/estado", servicioTickets.ActualizarEstado)

	// ---------------------------------------------------------
	// Variables de entorno y autenticación
	// ---------------------------------------------------------
	secreto := os.Getenv("TOKEN_SECRETO")
	if secreto == "" {
		log.Fatal("falta la variable de entorno TOKEN_SECRETO")
	}

	tokens := sesion.NuevosTokens(secreto)

	manejadorSesion := &sesion.Manejador{
		DB:     nil,
		Tokens: tokens,
	}

	r.With(httprate.LimitByIP(5, time.Minute)).
		Post("/login", manejadorSesion.Login)

	// ---------------------------------------------------------
	// Conexión a PostgreSQL
	// ---------------------------------------------------------
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("falta la variable de entorno DATABASE_URL")
	}

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("no se pudo conectar con PostgreSQL: ", err)
	}

	// GORM continúa creando y actualizando la tabla.
	if err := db.AutoMigrate(&nutricion.PlanNutricional{}); err != nil {
		log.Fatal("no se pudo migrar PlanNutricional: ", err)
	}

	// ---------------------------------------------------------
	// Datos iniciales de prueba
	// ---------------------------------------------------------
	var cantidadPlanes int64

	if err := db.Model(&nutricion.PlanNutricional{}).
		Count(&cantidadPlanes).Error; err != nil {
		log.Fatal("no se pudieron contar los planes: ", err)
	}

	if cantidadPlanes == 0 {
		planesIniciales := []nutricion.PlanNutricional{
			{
				IdMascota:   1,
				Estado:      "finalizado",
				Descripcion: "Plan de prueba finalizado",
			},
			{
				IdMascota:   1,
				Estado:      "pendiente",
				Descripcion: "Plan de prueba pendiente",
			},
		}

		if err := db.Create(&planesIniciales).Error; err != nil {
			log.Fatal("no se pudieron crear los planes iniciales: ", err)
		}
	}

	// ---------------------------------------------------------
	// Repositorio híbrido: GORM + sqlc
	// ---------------------------------------------------------
	conexion, err := db.DB()
	if err != nil {
		log.Fatal("no se pudo obtener la conexión SQL: ", err)
	}

	repoGORM := nutricion.NuevoRepositorioGORM(db)

	repoSqlc := nutricion.NuevoRepositorioPlanesSqlc(
		repoGORM,
		conexion,
	)

	servicioPlanes := nutricion.NuevoServicioPlanes(repoSqlc)

	manejadorNutricion := &nutricion.Manejador{
		DB:       db,
		Tokens:   tokens,
		Servicio: servicioPlanes,
	}

	manejadorNutricion.Rutas(r, tokens)

	// ---------------------------------------------------------
	// Ruta exclusiva para veterinarios
	// ---------------------------------------------------------
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
