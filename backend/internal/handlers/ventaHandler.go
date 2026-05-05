package handlers

import (
	"database/sql"
	"encoding/json"
	"magic-bag-gallery-api/internal/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetVentasHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT id_venta, id_cliente, id_empleado, fecha_venta, precio
		FROM venta
	`)
	if err != nil {
		http.Error(w, "Error al obtener ventas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	ventas := []models.Venta{}
	for rows.Next() {
		var v models.Venta
		if err := rows.Scan(&v.ID, &v.IDCliente, &v.IDEmpleado, &v.FechaVenta, &v.Precio); err != nil {
			http.Error(w, "Error al leer venta", http.StatusInternalServerError)
			return
		}
		ventas = append(ventas, v)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ventas)
}

func GetVentaByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var v models.Venta
	err = db.QueryRow(`
		SELECT id_venta, id_cliente, id_empleado, fecha_venta, precio
		FROM venta
		WHERE id_venta = $1
	`, id).Scan(&v.ID, &v.IDCliente, &v.IDEmpleado, &v.FechaVenta, &v.Precio)
	if err == sql.ErrNoRows {
		http.Error(w, "Venta no encontrada", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Error al obtener venta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func CreateVentaHandler(w http.ResponseWriter, r *http.Request) {
	var req models.VentaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body inválido", http.StatusBadRequest)
		return
	}

	if req.IDCliente == 0 || req.IDEmpleado == 0 || req.FechaVenta == "" || req.Precio == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	var clienteExists, empleadoExists bool
	if err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM cliente WHERE id_cliente = $1)`, req.IDCliente).Scan(&clienteExists); err != nil {
		http.Error(w, "Error al verificar cliente", http.StatusInternalServerError)
		return
	}
	if err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM empleado WHERE id_empleado = $1)`, req.IDEmpleado).Scan(&empleadoExists); err != nil {
		http.Error(w, "Error al verificar empleado", http.StatusInternalServerError)
		return
	}
	if !clienteExists {
		http.Error(w, "El cliente especificado no existe", http.StatusBadRequest)
		return
	}
	if !empleadoExists {
		http.Error(w, "El empleado especificado no existe", http.StatusBadRequest)
		return
	}

	var id int
	err := db.QueryRow(`
		INSERT INTO venta (id_cliente, id_empleado, fecha_venta, precio)
		VALUES ($1, $2, $3, $4)
		RETURNING id_venta
	`, req.IDCliente, req.IDEmpleado, req.FechaVenta, req.Precio).Scan(&id)
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
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req models.VentaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body inválido", http.StatusBadRequest)
		return
	}

	if req.IDCliente == 0 || req.IDEmpleado == 0 || req.FechaVenta == "" || req.Precio == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`
		UPDATE venta
		SET id_cliente = $1, id_empleado = $2, fecha_venta = $3, precio = $4
		WHERE id_venta = $5
	`, req.IDCliente, req.IDEmpleado, req.FechaVenta, req.Precio, id)
	if err != nil {
		http.Error(w, "Error al actualizar venta", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Venta no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Venta actualizada"})
}

func DeleteVentaHandler(w http.ResponseWriter, r *http.Request) {
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
	if err := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM venta WHERE id_venta = $1)`, id).Scan(&exists); err != nil {
		tx.Rollback()
		http.Error(w, "Error al verificar venta", http.StatusInternalServerError)
		return
	}
	if !exists {
		tx.Rollback()
		http.Error(w, "Venta no encontrada", http.StatusNotFound)
		return
	}

	if _, err := tx.Exec(`DELETE FROM envio WHERE id_venta = $1`, id); err != nil {
		tx.Rollback()
		http.Error(w, "Error al eliminar envíos asociados", http.StatusInternalServerError)
		return
	}

	if _, err := tx.Exec(`DELETE FROM detalle_venta WHERE id_venta = $1`, id); err != nil {
		tx.Rollback()
		http.Error(w, "Error al eliminar detalles de venta", http.StatusInternalServerError)
		return
	}

	if _, err := tx.Exec(`DELETE FROM venta WHERE id_venta = $1`, id); err != nil {
		tx.Rollback()
		http.Error(w, "Error al eliminar venta", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		http.Error(w, "Error al confirmar transacción", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Venta eliminada"})
}

func GetDetallesVentaHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT dv.id_detalle_venta, dv.id_venta, dv.id_pintura,
		       p.titulo AS titulo_pintura, dv.cantidad, dv.precio_unitario
		FROM detalle_venta dv
		JOIN pintura p ON dv.id_pintura = p.id_pintura
	`)
	if err != nil {
		http.Error(w, "Error al obtener detalles de venta", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	detalles := []models.DetalleVenta{}
	for rows.Next() {
		var d models.DetalleVenta
		if err := rows.Scan(&d.ID, &d.IDVenta, &d.IDPintura, &d.TituloPintura, &d.Cantidad, &d.PrecioUnitario); err != nil {
			http.Error(w, "Error al leer detalle de venta", http.StatusInternalServerError)
			return
		}
		detalles = append(detalles, d)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(detalles)
}

func GetDetalleVentaByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var d models.DetalleVenta
	err = db.QueryRow(`
		SELECT dv.id_detalle_venta, dv.id_venta, dv.id_pintura,
		       p.titulo AS titulo_pintura, dv.cantidad, dv.precio_unitario
		FROM detalle_venta dv
		JOIN pintura p ON dv.id_pintura = p.id_pintura
		WHERE dv.id_detalle_venta = $1
	`, id).Scan(&d.ID, &d.IDVenta, &d.IDPintura, &d.TituloPintura, &d.Cantidad, &d.PrecioUnitario)
	if err == sql.ErrNoRows {
		http.Error(w, "Detalle de venta no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Error al obtener detalle de venta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func GetDetallesByVentaHandler(w http.ResponseWriter, r *http.Request) {
	idVenta, err := strconv.Atoi(mux.Vars(r)["id_venta"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	rows, err := db.Query(`
		SELECT dv.id_detalle_venta, dv.id_venta, dv.id_pintura,
		       p.titulo AS titulo_pintura, dv.cantidad, dv.precio_unitario
		FROM detalle_venta dv
		JOIN pintura p ON dv.id_pintura = p.id_pintura
		WHERE dv.id_venta = $1
	`, idVenta)
	if err != nil {
		http.Error(w, "Error al obtener detalles de venta", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	detalles := []models.DetalleVenta{}
	for rows.Next() {
		var d models.DetalleVenta
		if err := rows.Scan(&d.ID, &d.IDVenta, &d.IDPintura, &d.TituloPintura, &d.Cantidad, &d.PrecioUnitario); err != nil {
			http.Error(w, "Error al leer detalle de venta", http.StatusInternalServerError)
			return
		}
		detalles = append(detalles, d)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(detalles)
}

func CreateDetalleVentaHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DetalleVentaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body inválido", http.StatusBadRequest)
		return
	}

	if req.IDVenta == 0 || req.IDPintura == 0 || req.Cantidad == 0 || req.PrecioUnitario == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	var ventaExists, pinturaExists bool
	if err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM venta WHERE id_venta = $1)`, req.IDVenta).Scan(&ventaExists); err != nil {
		http.Error(w, "Error al verificar venta", http.StatusInternalServerError)
		return
	}
	if err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM pintura WHERE id_pintura = $1)`, req.IDPintura).Scan(&pinturaExists); err != nil {
		http.Error(w, "Error al verificar pintura", http.StatusInternalServerError)
		return
	}
	if !ventaExists {
		http.Error(w, "La venta especificada no existe", http.StatusBadRequest)
		return
	}
	if !pinturaExists {
		http.Error(w, "La pintura especificada no existe", http.StatusBadRequest)
		return
	}

	var id int
	err := db.QueryRow(`
		INSERT INTO detalle_venta (id_venta, id_pintura, cantidad, precio_unitario)
		VALUES ($1, $2, $3, $4)
		RETURNING id_detalle_venta
	`, req.IDVenta, req.IDPintura, req.Cantidad, req.PrecioUnitario).Scan(&id)
	if err != nil {
		http.Error(w, "Error al crear detalle de venta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id_detalle_venta": id})
}

func UpdateDetalleVentaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req models.DetalleVentaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body inválido", http.StatusBadRequest)
		return
	}

	if req.IDVenta == 0 || req.IDPintura == 0 || req.Cantidad == 0 || req.PrecioUnitario == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`
		UPDATE detalle_venta
		SET id_venta = $1, id_pintura = $2, cantidad = $3, precio_unitario = $4
		WHERE id_detalle_venta = $5
	`, req.IDVenta, req.IDPintura, req.Cantidad, req.PrecioUnitario, id)
	if err != nil {
		http.Error(w, "Error al actualizar detalle de venta", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Detalle de venta no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Detalle de venta actualizado"})
}

func DeleteDetalleVentaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`DELETE FROM detalle_venta WHERE id_detalle_venta = $1`, id)
	if err != nil {
		http.Error(w, "Error al eliminar detalle de venta", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Detalle de venta no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Detalle de venta eliminado"})
}
