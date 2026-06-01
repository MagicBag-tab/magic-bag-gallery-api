package handlers

import (
	"encoding/json"
	"magic-bag-gallery-api/internal/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetUsuariosHandler(w http.ResponseWriter, r *http.Request) {
	var usuarios []models.Usuario
	err := gormDB.Raw(`
		SELECT
			u.id_usuario,
			u.nombre,
			u.apellido,
			u.correo_electronico,
			u.telefono,
			e.id_empleado,
			e.tipo_empleado
		FROM usuario u
		LEFT JOIN empleado e ON e.id_usuario = u.id_usuario
		ORDER BY u.nombre, u.apellido
	`).Scan(&usuarios).Error
	if err != nil {
		http.Error(w, "Error al obtener usuarios", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuarios)
}

func GetUsuarioByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	var usuario models.Usuario
	err = gormDB.Raw(`
		SELECT
			u.id_usuario,
			u.nombre,
			u.apellido,
			u.correo_electronico,
			u.telefono,
			e.id_empleado,
			e.tipo_empleado
		FROM usuario u
		LEFT JOIN empleado e ON e.id_usuario = u.id_usuario
		WHERE u.id_usuario = ?
	`, id).Scan(&usuario).Error
	if err != nil || usuario.ID == 0 {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

func CreateUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UsuarioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body invalido", http.StatusBadRequest)
		return
	}

	if req.Nombre == "" || req.Apellido == "" || req.CorreoElectronico == "" || req.Telefono == "" || req.Contrasena == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Contrasena), 12)
	if err != nil {
		http.Error(w, "Error al procesar contrasena", http.StatusInternalServerError)
		return
	}

	record := models.UsuarioEntity{
		Nombre:            req.Nombre,
		Apellido:          req.Apellido,
		CorreoElectronico: req.CorreoElectronico,
		Telefono:          req.Telefono,
		Contrasena:        string(hash),
	}
	if err := gormDB.Create(&record).Error; err != nil {
		http.Error(w, "El correo ya esta registrado", http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id_usuario": record.IDUsuario})
}

func UpdateUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	var req models.UsuarioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body invalido", http.StatusBadRequest)
		return
	}

	if req.Nombre == "" || req.Apellido == "" || req.CorreoElectronico == "" || req.Telefono == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	updates := map[string]interface{}{
		"nombre":             req.Nombre,
		"apellido":           req.Apellido,
		"correo_electronico": req.CorreoElectronico,
		"telefono":           req.Telefono,
	}
	if req.Contrasena != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Contrasena), 12)
		if err != nil {
			http.Error(w, "Error al procesar contrasena", http.StatusInternalServerError)
			return
		}
		updates["contrasena"] = string(hash)
	}

	result := gormDB.Model(&models.UsuarioEntity{}).Where("id_usuario = ?", id).Updates(updates)
	if result.Error != nil {
		http.Error(w, "Error al actualizar usuario", http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Usuario actualizado"})
}

func DeleteUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	if err := gormDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id_usuario = ?", id).Delete(&models.ClienteEntity{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id_usuario = ?", id).Delete(&models.EmpleadoEntity{}).Error; err != nil {
			return err
		}
		result := tx.Where("id_usuario = ?", id).Delete(&models.UsuarioEntity{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	}); err != nil {
		if isNotFound(err) {
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al eliminar usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Usuario eliminado"})
}
