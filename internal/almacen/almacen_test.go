package almacen

import (
	"testing"
)

func TestGuardarYLeer(t *testing.T) {
	ruta := t.TempDir() + "/datos.json"

	db, err := Nuevo(ruta)
	if err != nil {
		t.Fatal(err)
	}

	type Mascota struct {
		ID     int    `json:"id"`
		Nombre string `json:"nombre"`
	}

	original := []Mascota{
		{ID: 1, Nombre: "Firulais"},
	}

	if err := Guardar(db, "mascotas", original); err != nil {
		t.Fatal(err)
	}

	resultado, err := Leer[Mascota](db, "mascotas")
	if err != nil {
		t.Fatal(err)
	}

	if len(resultado) != 1 {
		t.Fatalf("se esperaba 1 mascota, se obtuvieron %d", len(resultado))
	}

	if resultado[0].Nombre != "Firulais" {
		t.Fatalf("nombre incorrecto: %s", resultado[0].Nombre)
	}
}

func TestModificar(t *testing.T) {
	db, err := Nuevo(t.TempDir() + "/datos.json")
	if err != nil {
		t.Fatal(err)
	}

	type Registro struct {
		ID int `json:"id"`
	}

	err = Modificar(db, "registros", func(registros []Registro) ([]Registro, error) {
		registros = append(registros, Registro{ID: 1})
		return registros, nil
	})

	if err != nil {
		t.Fatal(err)
	}

	resultado, err := Leer[Registro](db, "registros")
	if err != nil {
		t.Fatal(err)
	}

	if len(resultado) != 1 || resultado[0].ID != 1 {
		t.Fatal("el registro no se guardó correctamente")
	}
}
