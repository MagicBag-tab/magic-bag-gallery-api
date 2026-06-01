package handlers

import (
	"encoding/json"
	"magic-bag-gallery-api/internal/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetTecnicasHandler(w http.ResponseWriter, r *http.Request) {
	var records []models.TecnicaEntity
	if err := gormDB.Order("id_tecnica").Find(&records).Error; err != nil {
		http.Error(w, "Error al obtener tecnicas", http.StatusInternalServerError)
		return
	}

	tecnicas := make([]models.Tecnica, 0, len(records))
	for _, record := range records {
		tecnicas = append(tecnicas, models.Tecnica{
			ID:          record.IDTecnica,
			Nombre:      record.Nombre,
			Descripcion: record.Descripcion,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tecnicas)
}

func GetTecnicaByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	var record models.TecnicaEntity
	if err := gormDB.First(&record, "id_tecnica = ?", id).Error; err != nil {
		if isNotFound(err) {
			http.Error(w, "Tecnica no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al obtener tecnica", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Tecnica{
		ID:          record.IDTecnica,
		Nombre:      record.Nombre,
		Descripcion: record.Descripcion,
	})
}

func CreateTecnicaHandler(w http.ResponseWriter, r *http.Request) {
	var req models.TecnicaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body invalido", http.StatusBadRequest)
		return
	}

	if req.Nombre == "" || req.Descripcion == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	record := models.TecnicaEntity{Nombre: req.Nombre, Descripcion: req.Descripcion}
	if err := gormDB.Create(&record).Error; err != nil {
		http.Error(w, "Error al crear tecnica", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id_tecnica": record.IDTecnica})
}

func UpdateTecnicaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	var req models.TecnicaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body invalido", http.StatusBadRequest)
		return
	}

	result := gormDB.Model(&models.TecnicaEntity{}).
		Where("id_tecnica = ?", id).
		Updates(map[string]interface{}{"nombre": req.Nombre, "descripcion": req.Descripcion})
	if result.Error != nil {
		http.Error(w, "Error al actualizar tecnica", http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		http.Error(w, "Tecnica no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Tecnica actualizada"})
}

func DeleteTecnicaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	if err := gormDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id_tecnica = ?", id).Delete(&models.PinturaTecnicaEntity{}).Error; err != nil {
			return err
		}
		return tx.Where("id_tecnica = ?", id).Delete(&models.TecnicaEntity{}).Error
	}); err != nil {
		http.Error(w, "Error al eliminar tecnica", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Tecnica eliminada"})
}
