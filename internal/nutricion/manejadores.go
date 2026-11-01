package nutricion

import (
	"net/http"

	"mesa-ayuda/internal/sesion"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Manejador struct {
	DB       *gorm.DB
	Tokens   *sesion.Tokens
	Planes   map[uint]PlanNutricional
	Servicio *ServicioPlanes
}

func (m *Manejador) Rutas(r chi.Router, tokens *sesion.Tokens) {
	m.Tokens = tokens

	// Seleccionar el repositorio según exista o no una conexión a la base de datos.
	if m.Servicio == nil {
		var repositorio RepositorioPlanes

		if m.DB != nil {
			repositorio = NuevoRepositorioGORM(m.DB)
		} else {
			repositorio = NuevoRepositorioMemoria(m.Planes)
		}

		m.Servicio = NuevoServicioPlanes(repositorio)
	}

	r.Group(func(r chi.Router) {
		r.Use(sesion.Guardia(tokens))

		r.Get("/mascotas", m.listarMascotas)
		r.Post("/mascotas", m.crearMascota)

		r.Get("/planes/{id}", m.verPlan)
		r.Post("/planes", m.crearPlan)
		r.With(sesion.RequiereRol(sesion.RolVeterinario)).
  		  Patch("/planes/{id}", m.cambiarEstado)

		r.Post("/detalles", m.crearDetalle)
		r.Post("/citas", m.crearCita)
	})
}

func (m *Manejador) listarMascotas(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("listarMascotas OK"))
}

func (m *Manejador) crearMascota(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("crearMascota OK"))
}

func (m *Manejador) verPlan(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("verPlan OK"))
}

func (m *Manejador) crearPlan(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("crearPlan OK"))
}

func (m *Manejador) crearDetalle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("crearDetalle OK"))
}

func (m *Manejador) crearCita(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("crearCita OK"))
}
