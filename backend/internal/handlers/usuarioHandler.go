package handlers

import (
	"database/sql"
	"encoding/json"
	"magic-bag-gallery-api/internal/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func GetUsuariosHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT id_usuario, nombre, apellido, correo_electronico, telefono
		FROM usuario
	`)
	if err != nil {
		http.Error(w, "Error al obtener usuarios", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	usuarios := []models.Usuario{}
	for rows.Next() {
		var u models.Usuario
		if err := rows.Scan(&u.ID, &u.Nombre, &u.Apellido, &u.CorreoElectronico, &u.Telefono); err != nil {
			http.Error(w, "Error al leer usuario", http.StatusInternalServerError)
			return
		}
		usuarios = append(usuarios, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuarios)
}

func GetUsuarioByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var u models.Usuario
	err = db.QueryRow(`
		SELECT id_usuario, nombre, apellido, correo_electronico, telefono
		FROM usuario
		WHERE id_usuario = $1
	`, id).Scan(&u.ID, &u.Nombre, &u.Apellido, &u.CorreoElectronico, &u.Telefono)
	if err == sql.ErrNoRows {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Error al obtener usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func CreateUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UsuarioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body inválido", http.StatusBadRequest)
		return
	}

	if req.Nombre == "" || req.Apellido == "" || req.CorreoElectronico == "" || req.Telefono == "" || req.Contrasena == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Contrasena), 12)
	if err != nil {
		http.Error(w, "Error al procesar contraseña", http.StatusInternalServerError)
		return
	}

	var id int
	err = db.QueryRow(`
		INSERT INTO usuario (nombre, apellido, correo_electronico, telefono, contrasena)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id_usuario
	`, req.Nombre, req.Apellido, req.CorreoElectronico, req.Telefono, string(hash)).Scan(&id)
	if err != nil {
		http.Error(w, "El correo ya está registrado", http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id_usuario": id})
}

func UpdateUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req models.UsuarioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body inválido", http.StatusBadRequest)
		return
	}

	if req.Nombre == "" || req.Apellido == "" || req.CorreoElectronico == "" || req.Telefono == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	if req.Contrasena != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Contrasena), 12)
		if err != nil {
			http.Error(w, "Error al procesar contraseña", http.StatusInternalServerError)
			return
		}
		_, err = db.Exec(`
			UPDATE usuario SET nombre = $1, apellido = $2, correo_electronico = $3, telefono = $4, contrasena = $5
			WHERE id_usuario = $6
		`, req.Nombre, req.Apellido, req.CorreoElectronico, req.Telefono, string(hash), id)
	} else {
		_, err = db.Exec(`
			UPDATE usuario SET nombre = $1, apellido = $2, correo_electronico = $3, telefono = $4
			WHERE id_usuario = $5
		`, req.Nombre, req.Apellido, req.CorreoElectronico, req.Telefono, id)
	}
	if err != nil {
		http.Error(w, "Error al actualizar usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Usuario actualizado"})
}

func DeleteUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Error al iniciar transacción", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`DELETE FROM cliente WHERE id_usuario = $1`, id)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al eliminar cliente asociado", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`DELETE FROM empleado WHERE id_usuario = $1`, id)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al eliminar empleado asociado", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`DELETE FROM usuario WHERE id_usuario = $1`, id)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al eliminar usuario", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		http.Error(w, "Error al confirmar transacción", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Usuario eliminado"})
}
