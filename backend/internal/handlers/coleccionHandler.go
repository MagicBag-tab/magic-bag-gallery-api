package handlers

import (
	"database/sql"
	"encoding/json"
	"magic-bag-gallery-api/internal/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetColeccionesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT id_coleccion, nombre, descripcion, exclusiva, fecha_lanzamiento
		FROM coleccion
	`)
	if err != nil {
		http.Error(w, "Error al obtener colecciones", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	colecciones := []models.Coleccion{}
	for rows.Next() {
		var c models.Coleccion
		if err := rows.Scan(&c.ID, &c.Nombre, &c.Descripcion, &c.Exclusiva, &c.FechaLanzamiento); err != nil {
			http.Error(w, "Error al leer colección", http.StatusInternalServerError)
			return
		}
		colecciones = append(colecciones, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(colecciones)
}

func GetColeccionByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var c models.Coleccion
	err = db.QueryRow(`
		SELECT id_coleccion, nombre, descripcion, exclusiva, fecha_lanzamiento
		FROM coleccion
		WHERE id_coleccion = $1
	`, id).Scan(&c.ID, &c.Nombre, &c.Descripcion, &c.Exclusiva, &c.FechaLanzamiento)
	if err == sql.ErrNoRows {
		http.Error(w, "Colección no encontrada", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Error al obtener colección", http.StatusInternalServerError)
		return
	}

	pinturasRows, err := db.Query(`
		SELECT p.titulo
		FROM pintura p
		WHERE p.id_coleccion = $1
	`, id)
	if err != nil {
		http.Error(w, "Error al obtener pinturas", http.StatusInternalServerError)
		return
	}
	defer pinturasRows.Close()

	for pinturasRows.Next() {
		var titulo string
		if err := pinturasRows.Scan(&titulo); err != nil {
			http.Error(w, "Error al leer pintura", http.StatusInternalServerError)
			return
		}
		c.Pinturas = append(c.Pinturas, titulo)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func CreateColeccionHandler(w http.ResponseWriter, r *http.Request) {
	var req models.ColeccionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body inválido", http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Error al iniciar transacción", http.StatusInternalServerError)
		return
	}

	var id int
	err = tx.QueryRow(`
		INSERT INTO coleccion (nombre, descripcion, exclusiva, fecha_lanzamiento)
		VALUES ($1, $2, $3, $4)
		RETURNING id_coleccion
	`, req.Nombre, req.Descripcion, req.Exclusiva, req.FechaLanzamiento).Scan(&id)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al crear colección", http.StatusInternalServerError)
		return
	}

	for _, pinturaID := range req.IDPinturas {
		_, err := tx.Exec(`
			UPDATE pintura SET id_coleccion = $1 WHERE id_pintura = $2
		`, id, pinturaID)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error al asignar pintura a la colección", http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		http.Error(w, "Error al confirmar transacción", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id_coleccion": id})
}

func UpdateColeccionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req models.ColeccionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body inválido", http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Error al iniciar transacción", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`
		UPDATE coleccion SET nombre = $1, descripcion = $2, exclusiva = $3, fecha_lanzamiento = $4
		WHERE id_coleccion = $5
	`, req.Nombre, req.Descripcion, req.Exclusiva, req.FechaLanzamiento, id)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al actualizar colección", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`
		UPDATE pintura SET id_coleccion = NULL WHERE id_coleccion = $1
	`, id)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al desasignar pinturas", http.StatusInternalServerError)
		return
	}

	for _, pinturaID := range req.IDPinturas {
		_, err := tx.Exec(`
			UPDATE pintura SET id_coleccion = $1 WHERE id_pintura = $2
		`, id, pinturaID)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error al asignar pintura a la colección", http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		http.Error(w, "Error al confirmar transacción", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Colección actualizada"})
}

func DeleteColeccionHandler(w http.ResponseWriter, r *http.Request) {
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

	var exists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM coleccion WHERE id_coleccion = $1)`, id).Scan(&exists)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al verificar colección", http.StatusInternalServerError)
		return
	}
	if !exists {
		tx.Rollback()
		http.Error(w, "Colección no encontrada", http.StatusNotFound)
		return
	}

	_, err = tx.Exec(`
		UPDATE pintura SET id_coleccion = NULL WHERE id_coleccion = $1
	`, id)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al desasignar pinturas", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`DELETE FROM coleccion WHERE id_coleccion = $1`, id)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al eliminar colección", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		http.Error(w, "Error al confirmar transacción", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Colección eliminada"})
}
