package tickets

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"mesa-ayuda/internal/respuesta"

	"github.com/go-chi/chi/v5"
)

type Ticket struct {
	ID        uint
	Titulo    string
	Prioridad string
	Estado    string
}

// RepositorioMemoria implementa la interfaz Repositorio
type RepositorioMemoria struct {
	mu     sync.Mutex
	items  map[uint]Ticket
	ultimo uint
}

func NuevoRepositorioMemoria() *RepositorioMemoria {
	return &RepositorioMemoria{
		items: make(map[uint]Ticket),
	}
}

func (r *RepositorioMemoria) Listar(estado string) ([]Ticket, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var lista []Ticket
	for _, t := range r.items {
		if estado == "" || t.Estado == estado {
			lista = append(lista, t)
		}
	}
	return lista, nil
}

func (r *RepositorioMemoria) Buscar(id uint) (Ticket, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.items[id]
	if !ok {
		return Ticket{}, errors.New("no encontrado")
	}
	return t, nil
}

func (r *RepositorioMemoria) Crear(t *Ticket) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ultimo++
	t.ID = r.ultimo
	t.Estado = "abierto"
	r.items[t.ID] = *t
	return nil
}

func (r *RepositorioMemoria) Guardar(t *Ticket) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[t.ID]; !ok {
		return errors.New("no encontrado")
	}
	r.items[t.ID] = *t
	return nil
}

func (r *RepositorioMemoria) Borrar(id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[id]; !ok {
		return errors.New("no encontrado")
	}
	delete(r.items, id)
	return nil
}

// Servicio de tickets
type Servicio struct {
	repo Repositorio
}

func NuevoServicio(r Repositorio) *Servicio {
	return &Servicio{repo: r}
}

var prioridadesValidas = map[string]bool{"baja": true, "media": true, "alta": true}

var siguientes = map[string][]string{
	"abierto":    {"en_proceso"},
	"en_proceso": {"cerrado"},
	"cerrado":    {},
}

func puedePasar(de, a string) bool {
	for _, s := range siguientes[de] {
		if s == a {
			return true
		}
	}
	return false
}

func (s *Servicio) Listar(w http.ResponseWriter, r *http.Request) {
	estado := r.URL.Query().Get("estado")
	tickets, err := s.repo.Listar(estado)
	if err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "error_interno", err.Error())
		return
	}
	respuesta.Exito(w, http.StatusOK, tickets)
}

func (s *Servicio) Obtener(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil || id == 0 {
		respuesta.Error(w, http.StatusBadRequest, "id_invalido", "el id debe ser un entero positivo")
		return
	}
	t, err := s.repo.Buscar(uint(id))
	if err != nil {
		respuesta.Error(w, http.StatusNotFound, "no_encontrado", "no existe un ticket con ese id")
		return
	}
	respuesta.Exito(w, http.StatusOK, t)
}

func (s *Servicio) Crear(w http.ResponseWriter, r *http.Request) {
	var entrada struct {
		Titulo    string `json:"titulo"`
		Prioridad string `json:"prioridad"`
	}
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		respuesta.Error(w, http.StatusBadRequest, "cuerpo_malformado", "el cuerpo debe ser JSON con los campos titulo y prioridad")
		return
	}
	titulo := strings.TrimSpace(entrada.Titulo)
	if len(titulo) < 5 {
		respuesta.Error(w, http.StatusUnprocessableEntity, "datos_invalidos", "el título debe tener al menos 5 caracteres")
		return
	}
	if !prioridadesValidas[entrada.Prioridad] {
		respuesta.Error(w, http.StatusUnprocessableEntity, "datos_invalidos", "la prioridad debe ser baja, media o alta")
		return
	}

	nuevoTicket := &Ticket{
		Titulo:    titulo,
		Prioridad: entrada.Prioridad,
	}
	if err := s.repo.Crear(nuevoTicket); err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "error_interno", err.Error())
		return
	}
	respuesta.Exito(w, http.StatusCreated, nuevoTicket)
}

func (s *Servicio) ActualizarEstado(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil || id == 0 {
		respuesta.Error(w, http.StatusBadRequest, "id_invalido", "el id debe ser un entero positivo")
		return
	}

	var body struct {
		Estado string `json:"estado"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respuesta.Error(w, http.StatusBadRequest, "cuerpo_malformado", "el cuerpo debe ser JSON con campo estado")
		return
	}

	t, err := s.repo.Buscar(uint(id))
	if err != nil {
		respuesta.Error(w, http.StatusNotFound, "no_encontrado", "no existe un ticket con ese id")
		return
	}

	if _, existe := siguientes[body.Estado]; !existe {
		respuesta.Error(w, http.StatusUnprocessableEntity, "estado_invalido", "ese estado no existe en la máquina")
		return
	}

	if !puedePasar(t.Estado, body.Estado) {
		respuesta.Error(w, http.StatusUnprocessableEntity, "transicion_invalida", "no se puede pasar de "+t.Estado+" a "+body.Estado)
		return
	}

	t.Estado = body.Estado
	if err := s.repo.Guardar(&t); err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "error_interno", err.Error())
		return
	}

	respuesta.Exito(w, http.StatusOK, t)
}
