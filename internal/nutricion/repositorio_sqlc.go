package nutricion

import (
	"context"
	"database/sql"
	"errors"

	"mesa-ayuda/internal/almacen/filas"
)

type RepositorioPlanesSqlc struct {
	*RepositorioGORM
	Conexion *sql.DB
	q        *filas.Queries
}

// Comprobación en tiempo de compilación:
// RepositorioPlanesSqlc debe cumplir RepositorioPlanes.
var _ RepositorioPlanes = (*RepositorioPlanesSqlc)(nil)

func NuevoRepositorioPlanesSqlc(
	g *RepositorioGORM,
	conexion *sql.DB,
) *RepositorioPlanesSqlc {
	return &RepositorioPlanesSqlc{
		RepositorioGORM: g,
		Conexion:        conexion,
		q:               filas.New(conexion),
	}
}

// Obtener usa sqlc para buscar un plan por su ID.
func (r *RepositorioPlanesSqlc) Obtener(id uint) (PlanNutricional, error) {
	fila, err := r.q.BuscarPlan(
		context.Background(),
		int64(id),
	)

	if errors.Is(err, sql.ErrNoRows) {
		return PlanNutricional{}, ErrPlanNoEncontrado
	}

	if err != nil {
		return PlanNutricional{}, err
	}

	return PlanNutricional{
		ID:          uint(fila.ID),
		IdMascota:   uint(fila.IDMascota),
		Estado:      fila.Estado,
		Descripcion: fila.Descripcion,
		FechaInicio: fila.FechaInicio,
		FechaFin:    fila.FechaFin,
	}, nil
}

// Guardar usa sqlc para actualizar el plan.
func (r *RepositorioPlanesSqlc) Guardar(plan PlanNutricional) error {
	filasAfectadas, err := r.q.GuardarPlan(
		context.Background(),
		filas.GuardarPlanParams{
			ID:          int64(plan.ID),
			IDMascota:   int64(plan.IdMascota),
			Estado:      plan.Estado,
			Descripcion: plan.Descripcion,
			FechaInicio: plan.FechaInicio,
			FechaFin:    plan.FechaFin,
		},
	)

	if err == nil && filasAfectadas == 0 {
		return ErrPlanNoEncontrado
	}

	return err
}
