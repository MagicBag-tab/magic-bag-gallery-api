package handlers

import (
	"database/sql"
	"encoding/json"
	"magic-bag-gallery-api/internal/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetArtistasHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT id_artista, nombre_completo, nacionalidad, id_reclutador
		FROM artista
	`)
	if err != nil {
		http.Error(w, "Error al obtener artistas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	artistas := []models.Artista{}
	for rows.Next() {
		var a models.Artista
		if err := rows.Scan(&a.ID, &a.NombreCompleto, &a.Nacionalidad, &a.IDReclutador); err != nil {
			http.Error(w, "Error al leer artista", http.StatusInternalServerError)
			return
		}
		artistas = append(artistas, a)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artistas)
}

func GetArtistaByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var a models.Artista
	err = db.QueryRow(`
		SELECT id_artista, nombre_completo, nacionalidad, id_reclutador
		FROM artista
		WHERE id_artista = $1
	`, id).Scan(&a.ID, &a.NombreCompleto, &a.Nacionalidad, &a.IDReclutador)
	if err == sql.ErrNoRows {
		http.Error(w, "Artista no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Error al obtener artista", http.StatusInternalServerError)
		return
	}

	pinturasRows, err := db.Query(`
		SELECT p.titulo
		FROM pintura p
		WHERE p.id_artista = $1
	`, id)
	if err != nil {
		http.Error(w, "Error al obtener pinturas del artista", http.StatusInternalServerError)
		return
	}
	defer pinturasRows.Close()

	for pinturasRows.Next() {
		var titulo string
		if err := pinturasRows.Scan(&titulo); err != nil {
			http.Error(w, "Error al leer pintura", http.StatusInternalServerError)
			return
		}
		a.Pinturas = append(a.Pinturas, titulo)
	}

	coleccionesRows, err := db.Query(`
		SELECT DISTINCT c.nombre
		FROM coleccion c
		JOIN pintura p ON p.id_coleccion = c.id_coleccion
		WHERE p.id_artista = $1
	`, id)
	if err != nil {
		http.Error(w, "Error al obtener colecciones del artista", http.StatusInternalServerError)
		return
	}
	defer coleccionesRows.Close()

	for coleccionesRows.Next() {
		var nombre string
		if err := coleccionesRows.Scan(&nombre); err != nil {
			http.Error(w, "Error al leer colección", http.StatusInternalServerError)
			return
		}
		a.Colecciones = append(a.Colecciones, nombre)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(a)
}

func CreateArtistaHandler(w http.ResponseWriter, r *http.Request) {
	var req models.ArtistaRequest
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
		INSERT INTO artista (nombre_completo, nacionalidad, id_reclutador)
		VALUES ($1, $2, $3)
	`, req.NombreCompleto, req.Nacionalidad, req.IDReclutador).Scan(&id)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al crear artista", http.StatusInternalServerError)
		return
	}

	for _, pinturaID := range req.IDPinturas {
		_, err := tx.Exec(`
			UPDATE pintura SET id_artista = $1 WHERE id_pintura = $2
		`, id, pinturaID)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error al asignar pintura al artista", http.StatusInternalServerError)
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
	json.NewEncoder(w).Encode(map[string]int{"id_artista": id})
}

func UpdateArtistaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req models.ArtistaRequest
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
		UPDATE artista SET nombre_completo = $1, nacionalidad = $2, id_reclutador = $3
		WHERE id_artista = $4
	`, req.NombreCompleto, req.Nacionalidad, req.IDReclutador, id)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al actualizar artista", http.StatusInternalServerError)
		return
	}

	for _, pinturaID := range req.IDPinturas {
		_, err := tx.Exec(`
			UPDATE pintura SET id_artista = $1 WHERE id_pintura = $2
		`, id, pinturaID)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error al asignar pintura al artista", http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		http.Error(w, "Error al confirmar transacción", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Artista actualizado"})
}

func DeleteArtistaHandler(w http.ResponseWriter, r *http.Request) {
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
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM artista WHERE id_artista = $1)`, id).Scan(&exists)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al verificar artista", http.StatusInternalServerError)
		return
	}
	if !exists {
		tx.Rollback()
		http.Error(w, "Artista no encontrado", http.StatusNotFound)
		return
	}

	_, err = tx.Exec(`
		DELETE FROM pintura_tecnica
		WHERE id_pintura IN (SELECT id_pintura FROM pintura WHERE id_artista = $1)
	`, id)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al eliminar técnicas relacionadas", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`DELETE FROM pintura WHERE id_artista = $1`, id)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al eliminar pinturas del artista", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`DELETE FROM artista WHERE id_artista = $1`, id)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al eliminar artista", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		http.Error(w, "Error al confirmar transacción", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Artista eliminado"})
}
