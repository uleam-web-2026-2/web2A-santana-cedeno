package tickets

type Repositorio interface {
	Listar(estado string) ([]Ticket, error)
	Buscar(id uint) (Ticket, error)
	Crear(t *Ticket) error
	Guardar(t *Ticket) error
	Borrar(id uint) error
}
