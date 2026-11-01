package almacen

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

// Almacen administra la persistencia de los datos en un archivo JSON.
type Almacen struct {
	ruta string
	mu   sync.Mutex
}

// Nuevo crea o abre el archivo de almacenamiento.
func Nuevo(ruta string) (*Almacen, error) {
	directorio := filepath.Dir(ruta)

	if err := os.MkdirAll(directorio, 0755); err != nil {
		return nil, err
	}

	a := &Almacen{ruta: ruta}

	if _, err := os.Stat(ruta); errors.Is(err, os.ErrNotExist) {
		if err := a.escribir(map[string]json.RawMessage{}); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	// Comprobar que el archivo contiene JSON válido.
	if _, err := a.leer(); err != nil {
		return nil, err
	}

	return a, nil
}

// leer recupera todo el contenido del archivo.
func (a *Almacen) leer() (map[string]json.RawMessage, error) {
	contenido, err := os.ReadFile(a.ruta)
	if err != nil {
		return nil, err
	}

	datos := make(map[string]json.RawMessage)

	if len(contenido) == 0 {
		return datos, nil
	}

	if err := json.Unmarshal(contenido, &datos); err != nil {
		return nil, err
	}

	if datos == nil {
		datos = make(map[string]json.RawMessage)
	}

	return datos, nil
}

// escribir guarda los datos utilizando un archivo temporal.
func (a *Almacen) escribir(datos map[string]json.RawMessage) error {
	contenido, err := json.MarshalIndent(datos, "", "  ")
	if err != nil {
		return err
	}

	temporal := a.ruta + ".tmp"

	if err := os.WriteFile(temporal, contenido, 0644); err != nil {
		return err
	}

	if err := os.Rename(temporal, a.ruta); err != nil {
		os.Remove(temporal)
		return err
	}

	return nil
}

// Leer recupera una colección del archivo JSON.
func Leer[T any](a *Almacen, clave string) ([]T, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	datos, err := a.leer()
	if err != nil {
		return nil, err
	}

	contenido, existe := datos[clave]
	if !existe {
		return []T{}, nil
	}

	var registros []T

	if err := json.Unmarshal(contenido, &registros); err != nil {
		return nil, err
	}

	if registros == nil {
		registros = []T{}
	}

	return registros, nil
}

// Guardar almacena una colección completa.
func Guardar[T any](a *Almacen, clave string, registros []T) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	datos, err := a.leer()
	if err != nil {
		return err
	}

	contenido, err := json.Marshal(registros)
	if err != nil {
		return err
	}

	datos[clave] = contenido

	return a.escribir(datos)
}

// Modificar permite leer y actualizar una colección bajo el mismo bloqueo.
func Modificar[T any](
	a *Almacen,
	clave string,
	actualizar func([]T) ([]T, error),
) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	datos, err := a.leer()
	if err != nil {
		return err
	}

	var registros []T

	if contenido, existe := datos[clave]; existe {
		if err := json.Unmarshal(contenido, &registros); err != nil {
			return err
		}
	}

	if registros == nil {
		registros = []T{}
	}

	nuevos, err := actualizar(registros)
	if err != nil {
		return err
	}

	contenido, err := json.Marshal(nuevos)
	if err != nil {
		return err
	}

	datos[clave] = contenido

	return a.escribir(datos)
}
