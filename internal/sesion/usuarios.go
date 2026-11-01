package sesion

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"mesa-ayuda/internal/respuesta"
)

// Usuario para login y tokens
type Usuario struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Nombre    string `json:"nombre" gorm:"uniqueIndex"`
	Rol       string `json:"rol"`
	ClaveHash string `json:"clave_hash"`
}

type Manejador struct {
	DB     *gorm.DB
	Tokens *Tokens
}

// CrearUsuario maneja POST /usuarios
func CrearUsuario(w http.ResponseWriter, r *http.Request) {
	var datos struct {
		Nombre string `json:"nombre"`
		Rol    string `json:"rol"`
		Clave  string `json:"clave"`
	}

	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		respuesta.Error(w, http.StatusBadRequest, "json_invalido", "el cuerpo no es JSON válido")
		return
	}

	// Generar hash de la clave
	hash, err := bcrypt.GenerateFromPassword([]byte(datos.Clave), bcrypt.DefaultCost)
	if err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "error_hash", "no se pudo generar el hash")
		return
	}

	// Crear usuario en memoria (sin DB real)
	u := Usuario{
		ID:        1, // 👈 aquí puedes luego usar GORM para autoincrementar
		Nombre:    datos.Nombre,
		Rol:       datos.Rol,
		ClaveHash: string(hash),
	}

	respuesta.Exito(w, http.StatusCreated, u)
}

// Login: POST /login
func (m *Manejador) Login(w http.ResponseWriter, r *http.Request) {
	var pedido struct {
		Usuario string `json:"usuario"`
		Clave   string `json:"clave"`
	}

	if err := json.NewDecoder(r.Body).Decode(&pedido); err != nil {
		respuesta.Error(w, http.StatusBadRequest, "json_invalido", "el cuerpo no es JSON válido")
		return
	}

	var u Usuario

	switch {
	case pedido.Usuario == "Ana Cliente" && pedido.Clave == "secreta":
		u = Usuario{
			ID:     1,
			Nombre: "Ana Cliente",
			Rol:    RolCliente,
		}

	case pedido.Usuario == "Luis Veterinario" && pedido.Clave == "veterinario123":
		u = Usuario{
			ID:     2,
			Nombre: "Luis Veterinario",
			Rol:    RolVeterinario,
		}

	default:
		respuesta.Error(w, http.StatusUnauthorized, "credenciales_invalidas", "usuario o contraseña incorrectos")
		return
	}

	token, err := m.Tokens.Emitir(u)
	if err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "error_token", "no se pudo emitir el token")
		return
	}

	respuesta.Exito(w, http.StatusOK, map[string]string{"token": token})
}
