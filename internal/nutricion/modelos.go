package nutricion

type Usuario struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
	Correo string `json:"correo"`
}

type Mascota struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	IdUsuario   uint   `json:"id_usuario"`
	Nombre      string `json:"nombre"`
	Especie     string `json:"especie"`
	Raza        string `json:"raza"`
	Peso        string `json:"peso"`
	EstadoSalud string `json:"estado_salud"`
}

type PlanNutricional struct {
	ID          uint `gorm:"primaryKey"`
	IdMascota   uint
	Estado      string
	Descripcion string
	FechaInicio string
	FechaFin    string
}

type DetalleNutricional struct {
	ID                uint   `json:"id" gorm:"primaryKey"`
	IdPlanNutricional uint   `json:"id_plan_nutricional"`
	TipoAlimento      string `json:"tipo_alimento"`
	Cantidad          string `json:"cantidad"`
	Frecuencia        string `json:"frecuencia"`
	Suplementos       string `json:"suplementos"`
}

type Cita struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	IdUsuario uint   `json:"id_usuario"`
	IdMascota uint   `json:"id_mascota"`
	Fecha     string `json:"fecha"`
	Motivo    string `json:"motivo"`
}
