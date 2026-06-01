package handlers

import (
	"encoding/json"
	"magic-bag-gallery-api/internal/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetEnviosHandler(w http.ResponseWriter, r *http.Request) {
	var envios []models.Envio
	err := gormDB.Raw(`
		SELECT
			id_envio,
			id_venta,
			direccion_envio,
			TO_CHAR(fecha_envio, 'YYYY-MM-DD') AS fecha_envio,
			estado_envio
		FROM envio
		ORDER BY id_envio
	`).Scan(&envios).Error
	if err != nil {
		http.Error(w, "Error al obtener envios", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(envios)
}

func GetEnvioByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	var envio models.Envio
	err = gormDB.Raw(`
		SELECT
			id_envio,
			id_venta,
			direccion_envio,
			TO_CHAR(fecha_envio, 'YYYY-MM-DD') AS fecha_envio,
			estado_envio
		FROM envio
		WHERE id_envio = ?
	`, id).Scan(&envio).Error
	if err != nil || envio.ID == 0 {
		http.Error(w, "Envio no encontrado", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(envio)
}

func CreateEnvioHandler(w http.ResponseWriter, r *http.Request) {
	var req models.EnvioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body invalido", http.StatusBadRequest)
		return
	}

	if req.IDVenta == 0 || req.DireccionEnvio == "" || req.FechaEnvio == "" || req.EstadoEnvio == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	record := models.EnvioEntity{
		IDVenta:        req.IDVenta,
		DireccionEnvio: req.DireccionEnvio,
		FechaEnvio:     req.FechaEnvio,
		EstadoEnvio:    req.EstadoEnvio,
	}
	if err := gormDB.Create(&record).Error; err != nil {
		http.Error(w, "Error al crear envio", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id_envio": record.IDEnvio})
}

func UpdateEnvioHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	var req models.EnvioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body invalido", http.StatusBadRequest)
		return
	}

	result := gormDB.Model(&models.EnvioEntity{}).Where("id_envio = ?", id).Updates(map[string]interface{}{
		"id_venta":        req.IDVenta,
		"direccion_envio": req.DireccionEnvio,
		"fecha_envio":     req.FechaEnvio,
		"estado_envio":    req.EstadoEnvio,
	})
	if result.Error != nil {
		http.Error(w, "Error al actualizar envio", http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		http.Error(w, "Envio no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Envio actualizado"})
}

func DeleteEnvioHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	result := gormDB.Where("id_envio = ?", id).Delete(&models.EnvioEntity{})
	if result.Error != nil {
		http.Error(w, "Error al eliminar envio", http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		http.Error(w, "Envio no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Envio eliminado"})
}
