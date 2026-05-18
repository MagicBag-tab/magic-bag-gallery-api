package handlers

import (
	"encoding/json"
	"net/http"

	"magic-bag-gallery-api/internal/middleware"
)

func getUserIDFromContext(r *http.Request) (int, bool) {
	raw := r.Context().Value(middleware.UserIDKey)
	if raw == nil {
		return 0, false
	}
	switch v := raw.(type) {
	case float64:
		return int(v), true
	case int:
		return v, true
	}
	return 0, false
}

func GetMiTipoEmpleadoHandler(w http.ResponseWriter, r *http.Request) {
	idUsuario, ok := getUserIDFromContext(r)
	if !ok {
		http.Error(w, "Token inválido", http.StatusUnauthorized)
		return
	}

	var tipoEmpleado string
	err := db.QueryRow(`
		SELECT tipo_empleado FROM empleado WHERE id_usuario = $1
	`, idUsuario).Scan(&tipoEmpleado)
	if err != nil {
		http.Error(w, "Empleado no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"tipo_empleado": tipoEmpleado})
}

func GetMisReservasHandler(w http.ResponseWriter, r *http.Request) {
	idUsuario, ok := getUserIDFromContext(r)
	if !ok {
		http.Error(w, "Token inválido", http.StatusUnauthorized)
		return
	}

	var idCliente int
	err := db.QueryRow(`
		SELECT id_cliente FROM cliente WHERE id_usuario = $1 LIMIT 1
	`, idUsuario).Scan(&idCliente)
	if err != nil {
		http.Error(w, "Cliente no encontrado", http.StatusNotFound)
		return
	}

	rows, err := db.Query(`
		SELECT
			ct.id_cliente_tour,
			ct.id_tour,
			t.nombre     AS nombre_tour,
			ct.fecha_reserva,
			t.horario,
			t.precio
		FROM cliente_tour ct
		JOIN tour t ON ct.id_tour = t.id_tour
		WHERE ct.id_cliente = $1
		ORDER BY ct.fecha_reserva DESC
	`, idCliente)
	if err != nil {
		http.Error(w, "Error al obtener reservas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Reserva struct {
		ID           int     `json:"id_cliente_tour"`
		IDTour       int     `json:"id_tour"`
		NombreTour   string  `json:"nombre_tour"`
		FechaReserva string  `json:"fecha_reserva"`
		Horario      string  `json:"horario"`
		Precio       float64 `json:"precio"`
	}

	reservas := []Reserva{}
	for rows.Next() {
		var res Reserva
		if err := rows.Scan(
			&res.ID, &res.IDTour, &res.NombreTour,
			&res.FechaReserva, &res.Horario, &res.Precio,
		); err != nil {
			http.Error(w, "Error al leer reserva", http.StatusInternalServerError)
			return
		}
		reservas = append(reservas, res)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservas)
}

func GetMisVentasHandler(w http.ResponseWriter, r *http.Request) {
	idUsuario, ok := getUserIDFromContext(r)
	if !ok {
		http.Error(w, "Token inválido", http.StatusUnauthorized)
		return
	}

	var idCliente int
	err := db.QueryRow(`
		SELECT id_cliente FROM cliente WHERE id_usuario = $1 LIMIT 1
	`, idUsuario).Scan(&idCliente)
	if err != nil {
		http.Error(w, "Cliente no encontrado", http.StatusNotFound)
		return
	}

	rows, err := db.Query(`
		SELECT
			v.id_venta,
			v.fecha_venta,
			v.precio                                            AS total,
			STRING_AGG(p.titulo, ', ' ORDER BY p.titulo)       AS pinturas,
			COALESCE(e.estado_envio, 'Sin envío')              AS estado_envio
		FROM venta v
		LEFT JOIN detalle_venta dv ON dv.id_venta = v.id_venta
		LEFT JOIN pintura p        ON p.id_pintura = dv.id_pintura
		LEFT JOIN envio e          ON e.id_venta   = v.id_venta
		WHERE v.id_cliente = $1
		GROUP BY v.id_venta, v.fecha_venta, v.precio, e.estado_envio
		ORDER BY v.fecha_venta DESC
	`, idCliente)
	if err != nil {
		http.Error(w, "Error al obtener ventas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Venta struct {
		IDVenta     int     `json:"id_venta"`
		FechaVenta  string  `json:"fecha_venta"`
		Total       float64 `json:"total"`
		Pinturas    string  `json:"pinturas"`
		EstadoEnvio string  `json:"estado_envio"`
	}

	ventas := []Venta{}
	for rows.Next() {
		var v Venta
		if err := rows.Scan(&v.IDVenta, &v.FechaVenta, &v.Total, &v.Pinturas, &v.EstadoEnvio); err != nil {
			http.Error(w, "Error al leer venta", http.StatusInternalServerError)
			return
		}
		ventas = append(ventas, v)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ventas)
}
