package handlers

import (
	"encoding/json"
	"fmt"
	"magic-bag-gallery-api/internal/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type pinturaRow struct {
	ID            int     `gorm:"column:id_pintura"`
	Titulo        string  `gorm:"column:titulo"`
	Descripcion   string  `gorm:"column:descripcion"`
	FechaCreacion string  `gorm:"column:fecha_creacion"`
	Precio        float64 `gorm:"column:precio"`
	Exclusiva     bool    `gorm:"column:exclusiva"`
	ImagenPath    string  `gorm:"column:imagen_path"`
	ImagenTipo    string  `gorm:"column:imagen_tipo"`
	ImagenNombre  string  `gorm:"column:imagen_nombre"`
	Artista       string  `gorm:"column:artista"`
	Coleccion     string  `gorm:"column:coleccion"`
	Tecnicas      string  `gorm:"column:tecnicas"`
}

func GetPinturasHandler(w http.ResponseWriter, r *http.Request) {
	pinturas, err := loadPinturas("")
	if err != nil {
		http.Error(w, "Error al obtener pinturas", http.StatusInternalServerError)
		return
	}
	respondPinturas(w, pinturas)
}

func GetPinturaByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	pinturas, err := loadPinturas("WHERE p.id_pintura = ?", id)
	if err != nil {
		http.Error(w, "Error al obtener pintura", http.StatusInternalServerError)
		return
	}
	if len(pinturas) == 0 {
		http.Error(w, "Pintura no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pinturas[0])
}

func CreatePinturaHandler(w http.ResponseWriter, r *http.Request) {
	var req models.PinturaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body invalido", http.StatusBadRequest)
		return
	}

	record := models.PinturaEntity{
		IDArtista:     req.IDArtista,
		Titulo:        req.Titulo,
		Descripcion:   req.Descripcion,
		Precio:        req.Precio,
		FechaCreacion: req.FechaCreacion,
		ImagenPath:    req.ImagenPath,
		ImagenTipo:    req.ImagenTipo,
		ImagenNombre:  req.ImagenNombre,
		Exclusiva:     req.Exclusiva,
		IDColeccion:   req.IDColeccion,
	}

	if err := gormDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		return reemplazarTecnicasPintura(tx, record.IDPintura, req.Tecnicas)
	}); err != nil {
		http.Error(w, "Error al crear pintura", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id_pintura": record.IDPintura})
}

func UpdatePinturaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	var req models.PinturaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body invalido", http.StatusBadRequest)
		return
	}

	if err := gormDB.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.PinturaEntity{}).Where("id_pintura = ?", id).Updates(map[string]interface{}{
			"titulo":         req.Titulo,
			"descripcion":    req.Descripcion,
			"fecha_creacion": req.FechaCreacion,
			"precio":         req.Precio,
			"exclusiva":      req.Exclusiva,
			"imagen_path":    req.ImagenPath,
			"imagen_tipo":    req.ImagenTipo,
			"imagen_nombre":  req.ImagenNombre,
			"id_artista":     req.IDArtista,
			"id_coleccion":   req.IDColeccion,
		})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return reemplazarTecnicasPintura(tx, id, req.Tecnicas)
	}); err != nil {
		if isNotFound(err) {
			http.Error(w, "Pintura no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al actualizar pintura", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Pintura actualizada"})
}

func DeletePinturaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	if err := eliminarPinturaConProcedure(id); err != nil {
		if isNotFound(err) {
			http.Error(w, "Pintura no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al eliminar pintura", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Pintura eliminada"})
}

func GetPinturasByArtistaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id_artista"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}
	pinturas, err := loadPinturas("WHERE p.id_artista = ?", id)
	if err != nil {
		http.Error(w, "Error al obtener pinturas", http.StatusInternalServerError)
		return
	}
	respondPinturas(w, pinturas)
}

func GetPinturasByColeccionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id_coleccion"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}
	pinturas, err := loadPinturas("WHERE p.id_coleccion = ?", id)
	if err != nil {
		http.Error(w, "Error al obtener pinturas", http.StatusInternalServerError)
		return
	}
	respondPinturas(w, pinturas)
}

func GetPinturasByTecnicaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id_tecnica"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}
	pinturas, err := loadPinturas("WHERE EXISTS (SELECT 1 FROM pintura_tecnica ptx WHERE ptx.id_pintura = p.id_pintura AND ptx.id_tecnica = ?)", id)
	if err != nil {
		http.Error(w, "Error al obtener pinturas", http.StatusInternalServerError)
		return
	}
	respondPinturas(w, pinturas)
}

func loadPinturas(where string, args ...interface{}) ([]models.Pintura, error) {
	query := fmt.Sprintf(`
		SELECT
			p.id_pintura,
			p.titulo,
			p.descripcion,
			TO_CHAR(p.fecha_creacion, 'YYYY-MM-DD') AS fecha_creacion,
			p.precio::float8 AS precio,
			p.exclusiva,
			p.imagen_path,
			p.imagen_tipo,
			p.imagen_nombre,
			a.nombre_completo AS artista,
			COALESCE(c.nombre, '') AS coleccion,
			COALESCE(STRING_AGG(t.nombre, ', ' ORDER BY t.nombre), '') AS tecnicas
		FROM pintura p
		JOIN artista a ON p.id_artista = a.id_artista
		LEFT JOIN coleccion c ON p.id_coleccion = c.id_coleccion
		LEFT JOIN pintura_tecnica pt ON p.id_pintura = pt.id_pintura
		LEFT JOIN tecnica t ON pt.id_tecnica = t.id_tecnica
		%s
		GROUP BY p.id_pintura, p.titulo, p.descripcion, p.fecha_creacion,
			p.precio, p.exclusiva, p.imagen_path, p.imagen_tipo, p.imagen_nombre,
			a.nombre_completo, c.nombre
		ORDER BY p.id_pintura
	`, where)

	var rows []pinturaRow
	if err := gormDB.Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}

	pinturas := make([]models.Pintura, 0, len(rows))
	for _, row := range rows {
		p := models.Pintura{
			ID:            row.ID,
			Titulo:        row.Titulo,
			Descripcion:   row.Descripcion,
			FechaCreacion: row.FechaCreacion,
			Precio:        row.Precio,
			Exclusiva:     row.Exclusiva,
			ImagenPath:    row.ImagenPath,
			ImagenTipo:    row.ImagenTipo,
			ImagenNombre:  row.ImagenNombre,
			Artista:       row.Artista,
			Coleccion:     row.Coleccion,
		}
		if row.Tecnicas != "" {
			p.Tecnicas = strings.Split(row.Tecnicas, ", ")
		}
		pinturas = append(pinturas, p)
	}
	return pinturas, nil
}

func respondPinturas(w http.ResponseWriter, pinturas []models.Pintura) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pinturas)
}

func reemplazarTecnicasPintura(tx *gorm.DB, idPintura int, tecnicas []int) error {
	if err := tx.Where("id_pintura = ?", idPintura).Delete(&models.PinturaTecnicaEntity{}).Error; err != nil {
		return err
	}
	for _, idTecnica := range tecnicas {
		if err := tx.Create(&models.PinturaTecnicaEntity{IDPintura: idPintura, IDTecnica: idTecnica}).Error; err != nil {
			return err
		}
	}
	return nil
}

func eliminarPinturaConProcedure(id int) error {
	if err := gormDB.Exec("SELECT sp_eliminar_pintura(?)", id).Error; err == nil {
		return nil
	}
	return gormDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id_pintura = ?", id).Delete(&models.PinturaTecnicaEntity{}).Error; err != nil {
			return err
		}
		result := tx.Where("id_pintura = ?", id).Delete(&models.PinturaEntity{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
