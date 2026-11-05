package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
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
	// Autenticación
	// ---------------------------------------------------------
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

	// Semilla inicial para conservar los planes usados en las pruebas.
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

	// Ruta exclusiva para veterinarios.
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
