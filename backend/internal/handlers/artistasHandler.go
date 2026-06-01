package handlers

import (
	"encoding/json"
	"magic-bag-gallery-api/internal/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetArtistasHandler(w http.ResponseWriter, r *http.Request) {
	var rows []struct {
		ID               int    `gorm:"column:id_artista"`
		NombreCompleto   string `gorm:"column:nombre_completo"`
		Nacionalidad     string `gorm:"column:nacionalidad"`
		IDReclutador     int    `gorm:"column:id_reclutador"`
		NombreReclutador string `gorm:"column:nombre_reclutador"`
	}

	err := gormDB.Table("artista a").
		Select("a.id_artista, a.nombre_completo, a.nacionalidad, a.id_reclutador, u.nombre || ' ' || u.apellido AS nombre_reclutador").
		Joins("JOIN empleado e ON a.id_reclutador = e.id_empleado").
		Joins("JOIN usuario u ON e.id_usuario = u.id_usuario").
		Order("a.nombre_completo").
		Scan(&rows).Error
	if err != nil {
		http.Error(w, "Error al obtener artistas", http.StatusInternalServerError)
		return
	}

	artistas := make([]models.Artista, 0, len(rows))
	for _, row := range rows {
		artistas = append(artistas, models.Artista{
			ID:               row.ID,
			NombreCompleto:   row.NombreCompleto,
			Nacionalidad:     row.Nacionalidad,
			IDReclutador:     row.IDReclutador,
			NombreReclutador: row.NombreReclutador,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artistas)
}

func GetArtistaByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	var row struct {
		ID               int    `gorm:"column:id_artista"`
		NombreCompleto   string `gorm:"column:nombre_completo"`
		Nacionalidad     string `gorm:"column:nacionalidad"`
		IDReclutador     int    `gorm:"column:id_reclutador"`
		NombreReclutador string `gorm:"column:nombre_reclutador"`
	}

	err = gormDB.Table("artista a").
		Select("a.id_artista, a.nombre_completo, a.nacionalidad, a.id_reclutador, u.nombre || ' ' || u.apellido AS nombre_reclutador").
		Joins("JOIN empleado e ON a.id_reclutador = e.id_empleado").
		Joins("JOIN usuario u ON e.id_usuario = u.id_usuario").
		Where("a.id_artista = ?", id).
		Scan(&row).Error
	if err != nil || row.ID == 0 {
		http.Error(w, "Artista no encontrado", http.StatusNotFound)
		return
	}

	artista := models.Artista{
		ID:               row.ID,
		NombreCompleto:   row.NombreCompleto,
		Nacionalidad:     row.Nacionalidad,
		IDReclutador:     row.IDReclutador,
		NombreReclutador: row.NombreReclutador,
	}

	if err := gormDB.Model(&models.PinturaEntity{}).Where("id_artista = ?", id).Pluck("titulo", &artista.Pinturas).Error; err != nil {
		http.Error(w, "Error al obtener pinturas del artista", http.StatusInternalServerError)
		return
	}

	if err := gormDB.Table("coleccion c").
		Joins("JOIN pintura p ON p.id_coleccion = c.id_coleccion").
		Where("p.id_artista = ?", id).
		Distinct("c.nombre").
		Pluck("c.nombre", &artista.Colecciones).Error; err != nil {
		http.Error(w, "Error al obtener colecciones del artista", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artista)
}

func CreateArtistaHandler(w http.ResponseWriter, r *http.Request) {
	var req models.ArtistaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body invalido", http.StatusBadRequest)
		return
	}

	if !reclutadorExiste(req.IDReclutador) {
		http.Error(w, "El reclutador especificado no existe", http.StatusBadRequest)
		return
	}

	record := models.ArtistaEntity{
		NombreCompleto: req.NombreCompleto,
		Nacionalidad:   req.Nacionalidad,
		IDReclutador:   req.IDReclutador,
	}

	if err := gormDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		return asignarPinturasAArtista(tx, record.IDArtista, req.IDPinturas)
	}); err != nil {
		http.Error(w, "Error al crear artista", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id_artista": record.IDArtista})
}

func UpdateArtistaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	var req models.ArtistaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body invalido", http.StatusBadRequest)
		return
	}

	if !reclutadorExiste(req.IDReclutador) {
		http.Error(w, "El reclutador especificado no existe", http.StatusBadRequest)
		return
	}

	if err := gormDB.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.ArtistaEntity{}).Where("id_artista = ?", id).Updates(map[string]interface{}{
			"nombre_completo": req.NombreCompleto,
			"nacionalidad":    req.Nacionalidad,
			"id_reclutador":   req.IDReclutador,
		})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return asignarPinturasAArtista(tx, id, req.IDPinturas)
	}); err != nil {
		if isNotFound(err) {
			http.Error(w, "Artista no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al actualizar artista", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Artista actualizado"})
}

func DeleteArtistaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	if err := gormDB.Transaction(func(tx *gorm.DB) error {
		pinturas := tx.Model(&models.PinturaEntity{}).Select("id_pintura").Where("id_artista = ?", id)
		if err := tx.Where("id_pintura IN (?)", pinturas).Delete(&models.PinturaTecnicaEntity{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id_artista = ?", id).Delete(&models.PinturaEntity{}).Error; err != nil {
			return err
		}
		result := tx.Where("id_artista = ?", id).Delete(&models.ArtistaEntity{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	}); err != nil {
		if isNotFound(err) {
			http.Error(w, "Artista no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al eliminar artista", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Artista eliminado"})
}

func reclutadorExiste(idReclutador int) bool {
	var total int64
	err := gormDB.Model(&models.EmpleadoEntity{}).
		Where("id_empleado = ? AND tipo_empleado = ?", idReclutador, "reclutador").
		Count(&total).Error
	return err == nil && total > 0
}

func asignarPinturasAArtista(tx *gorm.DB, idArtista int, pinturas []int) error {
	if len(pinturas) == 0 {
		return nil
	}
	return tx.Model(&models.PinturaEntity{}).
		Where("id_pintura IN ?", pinturas).
		Update("id_artista", idArtista).Error
}
