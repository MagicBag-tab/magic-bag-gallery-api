package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Correo     string `json:"correo_electronico"`
	Contrasena string `json:"contrasena"`
}

type RegisterClienteRequest struct {
	Nombre            string `json:"nombre"`
	Apellido          string `json:"apellido"`
	CorreoElectronico string `json:"correo_electronico"`
	Telefono          string `json:"telefono"`
	Contrasena        string `json:"contrasena"`
	TipoCliente       string `json:"tipo_cliente"`
}

type RegisterEmpleadoRequest struct {
	Nombre            string `json:"nombre"`
	Apellido          string `json:"apellido"`
	CorreoElectronico string `json:"correo_electronico"`
	Telefono          string `json:"telefono"`
	Contrasena        string `json:"contrasena"`
	TipoEmpleado      string `json:"tipo_empleado"`
}

func generarToken(idUsuario int, role string) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	claims := jwt.MapClaims{
		"sub":  idUsuario,
		"role": role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body inválido", http.StatusBadRequest)
		return
	}

	var idUsuario int
	var hashContrasena string
	err := db.QueryRow(`
		SELECT id_usuario, contrasena FROM usuario WHERE correo_electronico = $1
	`, req.Correo).Scan(&idUsuario, &hashContrasena)
	if err != nil {
		http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashContrasena), []byte(req.Contrasena)); err != nil {
		http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
		return
	}

	var role string
	var dummy int
	if err := db.QueryRow(`SELECT id_empleado FROM empleado WHERE id_usuario = $1`, idUsuario).Scan(&dummy); err == nil {
		role = "empleado"
	} else {
		role = "cliente"
	}

	tokenStr, err := generarToken(idUsuario, role)
	if err != nil {
		http.Error(w, "Error al generar token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenStr, "role": role})
}

func RegisterClienteHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterClienteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body inválido", http.StatusBadRequest)
		return
	}

	if req.Nombre == "" || req.Apellido == "" || req.CorreoElectronico == "" || req.Contrasena == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}
	if req.TipoCliente == "" {
		req.TipoCliente = "regular"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Contrasena), 12)
	if err != nil {
		http.Error(w, "Error al procesar contraseña", http.StatusInternalServerError)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Error al iniciar transacción", http.StatusInternalServerError)
		return
	}

	var idUsuario int
	err = tx.QueryRow(`
		INSERT INTO usuario (nombre, apellido, correo_electronico, telefono, contrasena)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id_usuario
	`, req.Nombre, req.Apellido, req.CorreoElectronico, req.Telefono, string(hash)).Scan(&idUsuario)
	if err != nil {
		tx.Rollback()
		http.Error(w, "El correo ya está registrado", http.StatusConflict)
		return
	}

	var idCliente int
	err = tx.QueryRow(`
		INSERT INTO cliente (id_usuario, tipo_cliente) VALUES ($1, $2) RETURNING id_cliente
	`, idUsuario, req.TipoCliente).Scan(&idCliente)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al crear cliente", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		http.Error(w, "Error al confirmar transacción", http.StatusInternalServerError)
		return
	}

	tokenStr, err := generarToken(idUsuario, "cliente")
	if err != nil {
		http.Error(w, "Error al generar token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":      tokenStr,
		"role":       "cliente",
		"id_usuario": idUsuario,
		"id_cliente": idCliente,
	})
}

func RegisterEmpleadoHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterEmpleadoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body inválido", http.StatusBadRequest)
		return
	}

	if req.Nombre == "" || req.Apellido == "" || req.CorreoElectronico == "" || req.Contrasena == "" || req.TipoEmpleado == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Contrasena), 12)
	if err != nil {
		http.Error(w, "Error al procesar contraseña", http.StatusInternalServerError)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Error al iniciar transacción", http.StatusInternalServerError)
		return
	}

	var idUsuario int
	err = tx.QueryRow(`
		INSERT INTO usuario (nombre, apellido, correo_electronico, telefono, contrasena)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id_usuario
	`, req.Nombre, req.Apellido, req.CorreoElectronico, req.Telefono, string(hash)).Scan(&idUsuario)
	if err != nil {
		tx.Rollback()
		http.Error(w, "El correo ya está registrado", http.StatusConflict)
		return
	}

	var idEmpleado int
	err = tx.QueryRow(`
		INSERT INTO empleado (id_usuario, tipo_empleado) VALUES ($1, $2) RETURNING id_empleado
	`, idUsuario, req.TipoEmpleado).Scan(&idEmpleado)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al crear empleado", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		http.Error(w, "Error al confirmar transacción", http.StatusInternalServerError)
		return
	}

	tokenStr, err := generarToken(idUsuario, "empleado")
	if err != nil {
		http.Error(w, "Error al generar token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":       tokenStr,
		"role":        "empleado",
		"id_usuario":  idUsuario,
		"id_empleado": idEmpleado,
	})
}
