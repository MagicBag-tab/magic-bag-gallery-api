package handlers

import (
	"database/sql"
	"encoding/json"
	"magic-bag-gallery-api/internal/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetPinturasHandler(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT p.id_pintura, p.titulo, p.descripcion, p.fecha_creacion,
		       p.precio, p.exclusiva, p.imagen_path, p.imagen_tipo, p.imagen_nombre,
		       a.nombre_completo as artista,
		       COALESCE(c.nombre, '') as coleccion
		FROM pintura p
		JOIN artista a ON p.id_artista = a.id_artista
		LEFT JOIN coleccion c ON p.id_coleccion = c.id_coleccion
	`
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Error al obtener pinturas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	pinturas := []models.Pintura{}
	for rows.Next() {
		var p models.Pintura
		if err := rows.Scan(
			&p.ID, &p.Titulo, &p.Descripcion, &p.FechaCreacion,
			&p.Precio, &p.Exclusiva, &p.ImagenPath, &p.ImagenTipo, &p.ImagenNombre,
			&p.Artista, &p.Coleccion,
		); err != nil {
			http.Error(w, "Error al leer pintura", http.StatusInternalServerError)
			return
		}

		tecRows, err := db.Query(`
			SELECT t.nombre FROM tecnica t
			JOIN pintura_tecnica pt ON t.id_tecnica = pt.id_tecnica
			WHERE pt.id_pintura = $1
		`, p.ID)
		if err != nil {
			http.Error(w, "Error al obtener técnicas", http.StatusInternalServerError)
			return
		}
		for tecRows.Next() {
			var tec string
			if err := tecRows.Scan(&tec); err != nil {
				tecRows.Close()
				http.Error(w, "Error al leer técnica", http.StatusInternalServerError)
				return
			}
			p.Tecnicas = append(p.Tecnicas, tec)
		}
		tecRows.Close()

		pinturas = append(pinturas, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pinturas)
}

func GetPinturaByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var p models.Pintura
	err = db.QueryRow(`
		SELECT p.id_pintura, p.titulo, p.descripcion, p.fecha_creacion,
		       p.precio, p.exclusiva, p.imagen_path, p.imagen_tipo, p.imagen_nombre,
		       a.nombre as artista,
		       COALESCE(c.nombre, '') as coleccion
		FROM pintura p
		JOIN artista a ON p.id_artista = a.id_artista
		LEFT JOIN coleccion c ON p.id_coleccion = c.id_coleccion
		WHERE p.id_pintura = $1
	`, id).Scan(
		&p.ID, &p.Titulo, &p.Descripcion, &p.FechaCreacion,
		&p.Precio, &p.Exclusiva, &p.ImagenPath, &p.ImagenTipo, &p.ImagenNombre,
		&p.Artista, &p.Coleccion,
	)
	if err == sql.ErrNoRows {
		http.Error(w, "Pintura no encontrada", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Error al obtener pintura", http.StatusInternalServerError)
		return
	}

	tecRows, err := db.Query(`
		SELECT t.nombre FROM tecnica t
		JOIN pintura_tecnica pt ON t.id_tecnica = pt.id_tecnica
		WHERE pt.id_pintura = $1
	`, p.ID)
	if err != nil {
		http.Error(w, "Error al obtener técnicas", http.StatusInternalServerError)
		return
	}
	defer tecRows.Close()
	for tecRows.Next() {
		var tec string
		tecRows.Scan(&tec)
		p.Tecnicas = append(p.Tecnicas, tec)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func CreatePinturaHandler(w http.ResponseWriter, r *http.Request) {
	var req models.PinturaRequest
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
		INSERT INTO pintura (titulo, descripcion, fecha_creacion, precio, exclusiva,
		                     imagen_path, imagen_tipo, imagen_nombre, id_artista, id_coleccion)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, req.Titulo, req.Descripcion, req.FechaCreacion, req.Precio, req.Exclusiva,
		req.ImagenPath, req.ImagenTipo, req.ImagenNombre, req.IDArtista, req.IDColeccion,
	).Scan(&id)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al crear pintura", http.StatusInternalServerError)
		return
	}

	for _, idTec := range req.Tecnicas {
		_, err := tx.Exec(`
			INSERT INTO pintura_tecnica (id_pintura, id_tecnica) VALUES ($1, $2)
		`, id, idTec)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error al asignar técnica", http.StatusInternalServerError)
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
	json.NewEncoder(w).Encode(map[string]int{"id_pintura": id})
}

func UpdatePinturaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req models.PinturaRequest
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
		UPDATE pintura SET titulo=$1, descripcion=$2, fecha_creacion=$3, precio=$4,
		                   exclusiva=$5, imagen_path=$6, imagen_tipo=$7, imagen_nombre=$8,
		                   id_artista=$9, id_coleccion=$10
		WHERE id_pintura=$11
	`, req.Titulo, req.Descripcion, req.FechaCreacion, req.Precio, req.Exclusiva,
		req.ImagenPath, req.ImagenTipo, req.ImagenNombre, req.IDArtista, req.IDColeccion, id,
	)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al actualizar pintura", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`DELETE FROM pintura_tecnica WHERE id_pintura = $1`, id)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al actualizar técnicas", http.StatusInternalServerError)
		return
	}
	for _, idTec := range req.Tecnicas {
		_, err := tx.Exec(`
			INSERT INTO pintura_tecnica (id_pintura, id_tecnica) VALUES ($1, $2)
		`, id, idTec)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error al asignar técnica", http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		http.Error(w, "Error al confirmar transacción", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Pintura actualizada"})
}

func DeletePinturaHandler(w http.ResponseWriter, r *http.Request) {
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

	_, err = tx.Exec(`DELETE FROM pintura_tecnica WHERE id_pintura = $1`, id)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al eliminar técnicas", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`DELETE FROM pintura WHERE id_pintura = $1`, id)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al eliminar pintura", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		http.Error(w, "Error al confirmar transacción", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Pintura eliminada"})
}

func GetPinturasByArtistaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id_artista"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	rows, err := db.Query(`
		SELECT p.id_pintura, p.titulo, p.descripcion, p.fecha_creacion,
		       p.precio, p.exclusiva, p.imagen_path, p.imagen_tipo, p.imagen_nombre,
		       a.nombre_completo as artista,
		       COALESCE(c.nombre, '') as coleccion
		FROM pintura p
		JOIN artista a ON p.id_artista = a.id_artista
		LEFT JOIN coleccion c ON p.id_coleccion = c.id_coleccion
		WHERE p.id_artista = $1
	`, id)
	if err != nil {
		http.Error(w, "Error al obtener pinturas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	pinturas := []models.Pintura{}
	for rows.Next() {
		var p models.Pintura
		rows.Scan(
			&p.ID, &p.Titulo, &p.Descripcion, &p.FechaCreacion,
			&p.Precio, &p.Exclusiva, &p.ImagenPath, &p.ImagenTipo, &p.ImagenNombre,
			&p.Artista, &p.Coleccion,
		)
		pinturas = append(pinturas, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pinturas)
}

func GetPinturasByColeccionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id_coleccion"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	rows, err := db.Query(`
		SELECT p.id_pintura, p.titulo, p.descripcion, p.fecha_creacion,
		       p.precio, p.exclusiva, p.imagen_path, p.imagen_tipo, p.imagen_nombre,
		       a.nombre_completo as artista,
		       COALESCE(c.nombre, '') as coleccion
		FROM pintura p
		JOIN artista a ON p.id_artista = a.id_artista
		LEFT JOIN coleccion c ON p.id_coleccion = c.id_coleccion
		WHERE p.id_coleccion = $1
	`, id)
	if err != nil {
		http.Error(w, "Error al obtener pinturas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	pinturas := []models.Pintura{}
	for rows.Next() {
		var p models.Pintura
		rows.Scan(
			&p.ID, &p.Titulo, &p.Descripcion, &p.FechaCreacion,
			&p.Precio, &p.Exclusiva, &p.ImagenPath, &p.ImagenTipo, &p.ImagenNombre,
			&p.Artista, &p.Coleccion,
		)
		pinturas = append(pinturas, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pinturas)
}

func GetPinturasByTecnicaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id_tecnica"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	rows, err := db.Query(`
		SELECT p.id_pintura, p.titulo, p.descripcion, p.fecha_creacion,
		       p.precio, p.exclusiva, p.imagen_path, p.imagen_tipo, p.imagen_nombre,
		       a.nombre_completo as artista,
		       COALESCE(c.nombre, '') as coleccion
		FROM pintura p
		JOIN artista a ON p.id_artista = a.id_artista
		LEFT JOIN coleccion c ON p.id_coleccion = c.id_coleccion
		JOIN pintura_tecnica pt ON p.id_pintura = pt.id_pintura
		WHERE pt.id_tecnica = $1
	`, id)
	if err != nil {
		http.Error(w, "Error al obtener pinturas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	pinturas := []models.Pintura{}
	for rows.Next() {
		var p models.Pintura
		rows.Scan(
			&p.ID, &p.Titulo, &p.Descripcion, &p.FechaCreacion,
			&p.Precio, &p.Exclusiva, &p.ImagenPath, &p.ImagenTipo, &p.ImagenNombre,
			&p.Artista, &p.Coleccion,
		)
		pinturas = append(pinturas, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pinturas)
}
