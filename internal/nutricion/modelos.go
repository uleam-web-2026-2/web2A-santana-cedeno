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
	ID          uint   `json:"id" gorm:"primaryKey"`
	IdMascota   uint   `json:"id_mascota"`
	Estado      string `json:"estado"`
	Descripcion string `json:"descripcion"`
	FechaInicio string `json:"fecha_inicio"`
	FechaFin    string `json:"fecha_fin"`
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
