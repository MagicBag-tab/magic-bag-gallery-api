package handlers

import (
	"encoding/json"
	"magic-bag-gallery-api/internal/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetColeccionesHandler(w http.ResponseWriter, r *http.Request) {
	var records []models.ColeccionEntity
	if err := gormDB.Order("id_coleccion").Find(&records).Error; err != nil {
		http.Error(w, "Error al obtener colecciones", http.StatusInternalServerError)
		return
	}

	colecciones := make([]models.Coleccion, 0, len(records))
	for _, record := range records {
		colecciones = append(colecciones, coleccionFromEntity(record))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(colecciones)
}

func GetColeccionByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	var record models.ColeccionEntity
	if err := gormDB.First(&record, "id_coleccion = ?", id).Error; err != nil {
		if isNotFound(err) {
			http.Error(w, "Coleccion no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al obtener coleccion", http.StatusInternalServerError)
		return
	}

	coleccion := coleccionFromEntity(record)
	if err := gormDB.Model(&models.PinturaEntity{}).
		Where("id_coleccion = ?", id).
		Pluck("titulo", &coleccion.Pinturas).Error; err != nil {
		http.Error(w, "Error al obtener pinturas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(coleccion)
}

func CreateColeccionHandler(w http.ResponseWriter, r *http.Request) {
	var req models.ColeccionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body invalido", http.StatusBadRequest)
		return
	}

	record := models.ColeccionEntity{
		Nombre:           req.Nombre,
		Descripcion:      req.Descripcion,
		Exclusiva:        req.Exclusiva,
		FechaLanzamiento: req.FechaLanzamiento,
	}

	if err := gormDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		return asignarPinturasAColeccion(tx, record.IDColeccion, req.IDPinturas)
	}); err != nil {
		http.Error(w, "Error al crear coleccion", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id_coleccion": record.IDColeccion})
}

func UpdateColeccionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	var req models.ColeccionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body invalido", http.StatusBadRequest)
		return
	}

	if err := gormDB.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.ColeccionEntity{}).Where("id_coleccion = ?", id).Updates(map[string]interface{}{
			"nombre":            req.Nombre,
			"descripcion":       req.Descripcion,
			"exclusiva":         req.Exclusiva,
			"fecha_lanzamiento": req.FechaLanzamiento,
		})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		if err := tx.Model(&models.PinturaEntity{}).
			Where("id_coleccion = ?", id).
			Update("id_coleccion", nil).Error; err != nil {
			return err
		}
		return asignarPinturasAColeccion(tx, id, req.IDPinturas)
	}); err != nil {
		if isNotFound(err) {
			http.Error(w, "Coleccion no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al actualizar coleccion", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Coleccion actualizada"})
}

func DeleteColeccionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	if err := gormDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.PinturaEntity{}).
			Where("id_coleccion = ?", id).
			Update("id_coleccion", nil).Error; err != nil {
			return err
		}
		result := tx.Where("id_coleccion = ?", id).Delete(&models.ColeccionEntity{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	}); err != nil {
		if isNotFound(err) {
			http.Error(w, "Coleccion no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al eliminar coleccion", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Coleccion eliminada"})
}

func coleccionFromEntity(record models.ColeccionEntity) models.Coleccion {
	return models.Coleccion{
		ID:               record.IDColeccion,
		Nombre:           record.Nombre,
		Descripcion:      record.Descripcion,
		Exclusiva:        record.Exclusiva,
		FechaLanzamiento: record.FechaLanzamiento,
	}
}

func asignarPinturasAColeccion(tx *gorm.DB, idColeccion int, pinturas []int) error {
	if len(pinturas) == 0 {
		return nil
	}
	return tx.Model(&models.PinturaEntity{}).
		Where("id_pintura IN ?", pinturas).
		Update("id_coleccion", idColeccion).Error
}
