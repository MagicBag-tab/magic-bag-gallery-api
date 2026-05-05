package handlers

import (
	"database/sql"
	"encoding/json"
	"magic-bag-gallery-api/internal/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetEnviosHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT id_envio, id_venta, direccion_envio, fecha_envio, estado_envio
		FROM envio
	`)
	if err != nil {
		http.Error(w, "Error al obtener envíos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	envios := []models.Envio{}
	for rows.Next() {
		var e models.Envio
		if err := rows.Scan(&e.ID, &e.IDVenta, &e.DireccionEnvio, &e.FechaEnvio, &e.EstadoEnvio); err != nil {
			http.Error(w, "Error al leer envío", http.StatusInternalServerError)
			return
		}
		envios = append(envios, e)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(envios)
}

func GetEnvioByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var e models.Envio
	err = db.QueryRow(`
		SELECT id_envio, id_venta, direccion_envio, fecha_envio, estado_envio
		FROM envio
		WHERE id_envio = $1
	`, id).Scan(&e.ID, &e.IDVenta, &e.DireccionEnvio, &e.FechaEnvio, &e.EstadoEnvio)
	if err == sql.ErrNoRows {
		http.Error(w, "Envío no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Error al obtener envío", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e)
}

func CreateEnvioHandler(w http.ResponseWriter, r *http.Request) {
	var req models.EnvioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body inválido", http.StatusBadRequest)
		return
	}

	if req.IDVenta == 0 || req.DireccionEnvio == "" || req.FechaEnvio == "" || req.EstadoEnvio == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	var exists bool
	if err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM venta WHERE id_venta = $1)`, req.IDVenta).Scan(&exists); err != nil {
		http.Error(w, "Error al verificar venta", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "La venta especificada no existe", http.StatusBadRequest)
		return
	}

	var id int
	err := db.QueryRow(`
		INSERT INTO envio (id_venta, direccion_envio, fecha_envio, estado_envio)
		VALUES ($1, $2, $3, $4)
		RETURNING id_envio
	`, req.IDVenta, req.DireccionEnvio, req.FechaEnvio, req.EstadoEnvio).Scan(&id)
	if err != nil {
		http.Error(w, "Error al crear envío", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id_envio": id})
}

func UpdateEnvioHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req models.EnvioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body inválido", http.StatusBadRequest)
		return
	}

	if req.DireccionEnvio == "" || req.FechaEnvio == "" || req.EstadoEnvio == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`
		UPDATE envio
		SET id_venta = $1, direccion_envio = $2, fecha_envio = $3, estado_envio = $4
		WHERE id_envio = $5
	`, req.IDVenta, req.DireccionEnvio, req.FechaEnvio, req.EstadoEnvio, id)
	if err != nil {
		http.Error(w, "Error al actualizar envío", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Envío no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Envío actualizado"})
}

func DeleteEnvioHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`DELETE FROM envio WHERE id_envio = $1`, id)
	if err != nil {
		http.Error(w, "Error al eliminar envío", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Envío no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Envío eliminado"})
}
