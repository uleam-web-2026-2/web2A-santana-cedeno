package nutricion

import (
	"errors"
	"fmt"
	"slices"
)

var ErrEstadoInvalido = errors.New("estado solicitado no permitido")
var ErrTransicionNoPermitida = errors.New("transición de estado no permitida")

var siguientes = map[string][]string{
	"pendiente":  {"activo", "cancelado"},
	"activo":     {"finalizado", "cancelado"},
	"finalizado": {},
	"cancelado":  {},
}

type ServicioPlanes struct {
	Repositorio RepositorioPlanes
}

func NuevoServicioPlanes(repo RepositorioPlanes) *ServicioPlanes {
	return &ServicioPlanes{Repositorio: repo}
}

func (s *ServicioPlanes) CambiarEstado(id uint, nuevoEstado string) (PlanNutricional, error) {
	if _, existe := siguientes[nuevoEstado]; !existe {
		return PlanNutricional{}, ErrEstadoInvalido
	}

	plan, err := s.Repositorio.Obtener(id)
	if err != nil {
		return PlanNutricional{}, err
	}

	if !slices.Contains(siguientes[plan.Estado], nuevoEstado) {
		return PlanNutricional{}, fmt.Errorf(
			"%w: de %s no se puede pasar a %s",
			ErrTransicionNoPermitida,
			plan.Estado,
			nuevoEstado,
		)
	}

	plan.Estado = nuevoEstado

	if err := s.Repositorio.Guardar(plan); err != nil {
		return PlanNutricional{}, err
	}

	return plan, nil
}
