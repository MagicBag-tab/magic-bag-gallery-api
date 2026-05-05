package handlers

import (
	"database/sql"
	"encoding/json"
	"magic-bag-gallery-api/internal/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetTecnicasHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT id_tecnica, nombre, descripcion
		FROM tecnica
	`)
	if err != nil {
		http.Error(w, "Error al obtener técnicas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tecnicas := []models.Tecnica{}
	for rows.Next() {
		var t models.Tecnica
		if err := rows.Scan(&t.ID, &t.Nombre, &t.Descripcion); err != nil {
			http.Error(w, "Error al leer técnica", http.StatusInternalServerError)
			return
		}
		tecnicas = append(tecnicas, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tecnicas)
}

func GetTecnicaByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var t models.Tecnica
	err = db.QueryRow(`
		SELECT id_tecnica, nombre, descripcion
		FROM tecnica
		WHERE id_tecnica = $1
	`, id).Scan(&t.ID, &t.Nombre, &t.Descripcion)
	if err == sql.ErrNoRows {
		http.Error(w, "Técnica no encontrada", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Error al obtener técnica", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func CreateTecnicaHandler(w http.ResponseWriter, r *http.Request) {
	var req models.TecnicaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body inválido", http.StatusBadRequest)
		return
	}

	if req.Nombre == "" || req.Descripcion == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	var id int
	err := db.QueryRow(`
		INSERT INTO tecnica (nombre, descripcion)
		VALUES ($1, $2)
		RETURNING id_tecnica
	`, req.Nombre, req.Descripcion).Scan(&id)
	if err != nil {
		http.Error(w, "Error al crear técnica", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id_tecnica": id})
}

func UpdateTecnicaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req models.TecnicaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body inválido", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`
		UPDATE tecnica SET nombre = $1, descripcion = $2
		WHERE id_tecnica = $3
	`, req.Nombre, req.Descripcion, id)
	if err != nil {
		http.Error(w, "Error al actualizar técnica", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Técnica actualizada"})
}

func DeleteTecnicaHandler(w http.ResponseWriter, r *http.Request) {
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

	_, err = tx.Exec(`
		DELETE FROM pintura_tecnica WHERE id_tecnica = $1
	`, id)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al eliminar relaciones de técnica", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`
		DELETE FROM tecnica WHERE id_tecnica = $1
	`, id)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al eliminar técnica", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		http.Error(w, "Error al confirmar transacción", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Técnica eliminada"})
}
