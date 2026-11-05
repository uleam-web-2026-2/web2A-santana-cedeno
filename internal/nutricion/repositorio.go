package nutricion

import (
	"errors"

	"gorm.io/gorm"
)

var ErrPlanNoEncontrado = errors.New("plan nutricional no encontrado")

// RepositorioPlanes define cómo se consultan y guardan los planes.
type RepositorioPlanes interface {
	Obtener(id uint) (PlanNutricional, error)
	Crear(plan *PlanNutricional) error
	Guardar(plan PlanNutricional) error
}

// RepositorioMemoria permite trabajar sin una base de datos.
type RepositorioMemoria struct {
	Planes map[uint]PlanNutricional
}

func NuevoRepositorioMemoria(planes map[uint]PlanNutricional) *RepositorioMemoria {
	if planes == nil {
		planes = make(map[uint]PlanNutricional)
	}

	return &RepositorioMemoria{Planes: planes}
}

func (r *RepositorioMemoria) Obtener(id uint) (PlanNutricional, error) {
	plan, existe := r.Planes[id]
	if !existe {
		return PlanNutricional{}, ErrPlanNoEncontrado
	}

	return plan, nil
}

func (r *RepositorioMemoria) Crear(plan *PlanNutricional) error {
	var nuevoID uint = 1

	for id := range r.Planes {
		if id >= nuevoID {
			nuevoID = id + 1
		}
	}

	plan.ID = nuevoID
	r.Planes[plan.ID] = *plan

	return nil
}

func (r *RepositorioMemoria) Guardar(plan PlanNutricional) error {
	r.Planes[plan.ID] = plan
	return nil
}

// RepositorioGORM permite trabajar con una base de datos mediante GORM.
type RepositorioGORM struct {
	DB *gorm.DB
}

func NuevoRepositorioGORM(db *gorm.DB) *RepositorioGORM {
	return &RepositorioGORM{DB: db}
}

func (r *RepositorioGORM) Obtener(id uint) (PlanNutricional, error) {
	var plan PlanNutricional

	err := r.DB.First(&plan, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return PlanNutricional{}, ErrPlanNoEncontrado
	}

	return plan, err
}

func (r *RepositorioGORM) Crear(plan *PlanNutricional) error {
	return r.DB.Create(plan).Error
}

func (r *RepositorioGORM) Guardar(plan PlanNutricional) error {
	return r.DB.Save(&plan).Error
}
