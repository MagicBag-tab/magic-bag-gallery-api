package handlers

import (
	"encoding/json"
	"magic-bag-gallery-api/internal/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetVentasHandler(w http.ResponseWriter, r *http.Request) {
	var ventas []models.Venta
	err := gormDB.Raw(`
		SELECT id_venta, id_cliente, id_empleado, TO_CHAR(fecha_venta, 'YYYY-MM-DD') AS fecha_venta, precio::text AS precio
		FROM venta
		ORDER BY fecha_venta DESC, id_venta DESC
	`).Scan(&ventas).Error
	if err != nil {
		http.Error(w, "Error al obtener ventas", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ventas)
}

func GetVentaByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	var venta models.Venta
	err = gormDB.Raw(`
		SELECT id_venta, id_cliente, id_empleado, TO_CHAR(fecha_venta, 'YYYY-MM-DD') AS fecha_venta, precio::text AS precio
		FROM venta
		WHERE id_venta = ?
	`, id).Scan(&venta).Error
	if err != nil || venta.ID == 0 {
		http.Error(w, "Venta no encontrada", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(venta)
}

func CreateVentaHandler(w http.ResponseWriter, r *http.Request) {
	var req models.VentaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body invalido", http.StatusBadRequest)
		return
	}

	if req.IDCliente == 0 || req.IDEmpleado == 0 || req.FechaVenta == "" || req.Precio == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	id, err := crearVentaConProcedure(req)
	if err != nil {
		http.Error(w, "Error al crear venta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id_venta": id})
}

func UpdateVentaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	var req models.VentaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body invalido", http.StatusBadRequest)
		return
	}

	result := gormDB.Model(&models.VentaEntity{}).Where("id_venta = ?", id).Updates(map[string]interface{}{
		"id_cliente":  req.IDCliente,
		"id_empleado": req.IDEmpleado,
		"fecha_venta": req.FechaVenta,
		"precio":      req.Precio,
	})
	if result.Error != nil {
		http.Error(w, "Error al actualizar venta", http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		http.Error(w, "Venta no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Venta actualizada"})
}

func DeleteVentaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	if err := gormDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id_venta = ?", id).Delete(&models.EnvioEntity{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id_venta = ?", id).Delete(&models.DetalleVentaEntity{}).Error; err != nil {
			return err
		}
		result := tx.Where("id_venta = ?", id).Delete(&models.VentaEntity{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	}); err != nil {
		if isNotFound(err) {
			http.Error(w, "Venta no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al eliminar venta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Venta eliminada"})
}

func GetDetallesVentaHandler(w http.ResponseWriter, r *http.Request) {
	detalles, err := loadDetallesVenta("")
	if err != nil {
		http.Error(w, "Error al obtener detalles de venta", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(detalles)
}

func GetDetalleVentaByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	detalles, err := loadDetallesVenta("WHERE dv.id_detalle_venta = ?", id)
	if err != nil {
		http.Error(w, "Error al obtener detalle de venta", http.StatusInternalServerError)
		return
	}
	if len(detalles) == 0 {
		http.Error(w, "Detalle de venta no encontrado", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(detalles[0])
}

func GetDetallesByVentaHandler(w http.ResponseWriter, r *http.Request) {
	idVenta, err := strconv.Atoi(mux.Vars(r)["id_venta"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	detalles, err := loadDetallesVenta("WHERE dv.id_venta = ?", idVenta)
	if err != nil {
		http.Error(w, "Error al obtener detalles de venta", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(detalles)
}

func CreateDetalleVentaHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DetalleVentaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body invalido", http.StatusBadRequest)
		return
	}

	record := models.DetalleVentaEntity{
		IDVenta:        req.IDVenta,
		IDPintura:      req.IDPintura,
		Cantidad:       req.Cantidad,
		PrecioUnitario: req.PrecioUnitario,
	}
	if err := gormDB.Create(&record).Error; err != nil {
		http.Error(w, "Error al crear detalle de venta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id_detalle_venta": record.IDDetalleVenta})
}

func UpdateDetalleVentaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	var req models.DetalleVentaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body invalido", http.StatusBadRequest)
		return
	}

	result := gormDB.Model(&models.DetalleVentaEntity{}).Where("id_detalle_venta = ?", id).Updates(map[string]interface{}{
		"id_venta":        req.IDVenta,
		"id_pintura":      req.IDPintura,
		"cantidad":        req.Cantidad,
		"precio_unitario": req.PrecioUnitario,
	})
	if result.Error != nil {
		http.Error(w, "Error al actualizar detalle de venta", http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		http.Error(w, "Detalle de venta no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Detalle de venta actualizado"})
}

func DeleteDetalleVentaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	result := gormDB.Where("id_detalle_venta = ?", id).Delete(&models.DetalleVentaEntity{})
	if result.Error != nil {
		http.Error(w, "Error al eliminar detalle de venta", http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		http.Error(w, "Detalle de venta no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Detalle de venta eliminado"})
}

func crearVentaConProcedure(req models.VentaRequest) (int, error) {
	type result struct {
		IDVenta int `gorm:"column:id_venta"`
	}
	var res result
	err := gormDB.Raw("SELECT * FROM sp_crear_venta(?, ?, ?, ?)", req.IDCliente, req.IDEmpleado, req.FechaVenta, req.Precio).Scan(&res).Error
	if err == nil && res.IDVenta != 0 {
		return res.IDVenta, nil
	}

	record := models.VentaEntity{
		IDCliente:  req.IDCliente,
		IDEmpleado: req.IDEmpleado,
		FechaVenta: req.FechaVenta,
		Precio:     req.Precio,
	}
	if err := gormDB.Create(&record).Error; err != nil {
		return 0, err
	}
	return record.IDVenta, nil
}

func loadDetallesVenta(where string, args ...interface{}) ([]models.DetalleVenta, error) {
	var detalles []models.DetalleVenta
	query := `
		SELECT
			dv.id_detalle_venta,
			dv.id_venta,
			dv.id_pintura,
			p.titulo AS titulo_pintura,
			dv.cantidad,
			dv.precio_unitario::text AS precio_unitario
		FROM detalle_venta dv
		JOIN pintura p ON dv.id_pintura = p.id_pintura
		` + where + `
		ORDER BY dv.id_detalle_venta
	`
	return detalles, gormDB.Raw(query, args...).Scan(&detalles).Error
}
