package nutricion

type PlanNutricionalPeticion struct {
	IdMascota   uint   `json:"id_mascota"`
	Descripcion string `json:"descripcion"`
	FechaInicio string `json:"fecha_inicio"`
	FechaFin    string `json:"fecha_fin"`
}

func (p PlanNutricionalPeticion) aEntidad() PlanNutricional {
	return PlanNutricional{
		IdMascota:   p.IdMascota,
		Descripcion: p.Descripcion,
		FechaInicio: p.FechaInicio,
		FechaFin:    p.FechaFin,
	}
}

type PlanNutricionalRespuesta struct {
	ID          uint   `json:"id"`
	IdMascota   uint   `json:"id_mascota"`
	Estado      string `json:"estado"`
	Descripcion string `json:"descripcion"`
	FechaInicio string `json:"fecha_inicio"`
	FechaFin    string `json:"fecha_fin"`
}

func nuevaRespuesta(p PlanNutricional) PlanNutricionalRespuesta {
	return PlanNutricionalRespuesta{
		ID:          p.ID,
		IdMascota:   p.IdMascota,
		Estado:      p.Estado,
		Descripcion: p.Descripcion,
		FechaInicio: p.FechaInicio,
		FechaFin:    p.FechaFin,
	}
}
