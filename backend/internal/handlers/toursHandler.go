package handlers

import (
	"database/sql"
	"encoding/json"
	"magic-bag-gallery-api/internal/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetToursHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT t.id_tour, t.id_guia,
		       u.nombre || ' ' || u.apellido AS nombre_guia,
		       t.nombre, t.descripcion,
		       t.fecha_inicio, t.fecha_fin, t.horario, t.precio
		FROM tour t
		JOIN empleado e ON t.id_guia = e.id_empleado
		JOIN usuario  u ON e.id_usuario = u.id_usuario
	`)
	if err != nil {
		http.Error(w, "Error al obtener tours", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tours := []models.Tour{}
	for rows.Next() {
		var t models.Tour
		if err := rows.Scan(
			&t.ID, &t.IDGuia, &t.NombreGuia,
			&t.Nombre, &t.Descripcion,
			&t.FechaInicio, &t.FechaFin, &t.Horario, &t.Precio,
		); err != nil {
			http.Error(w, "Error al leer tour", http.StatusInternalServerError)
			return
		}
		tours = append(tours, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tours)
}

func GetTourByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var t models.Tour
	err = db.QueryRow(`
		SELECT t.id_tour, t.id_guia,
		       u.nombre || ' ' || u.apellido AS nombre_guia,
		       t.nombre, t.descripcion,
		       t.fecha_inicio, t.fecha_fin, t.horario, t.precio
		FROM tour t
		JOIN empleado e ON t.id_guia = e.id_empleado
		JOIN usuario  u ON e.id_usuario = u.id_usuario
		WHERE t.id_tour = $1
	`, id).Scan(
		&t.ID, &t.IDGuia, &t.NombreGuia,
		&t.Nombre, &t.Descripcion,
		&t.FechaInicio, &t.FechaFin, &t.Horario, &t.Precio,
	)
	if err == sql.ErrNoRows {
		http.Error(w, "Tour no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Error al obtener tour", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func CreateTourHandler(w http.ResponseWriter, r *http.Request) {
	var req models.TourRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body inválido", http.StatusBadRequest)
		return
	}

	if req.IDGuia == 0 || req.Nombre == "" || req.FechaInicio == "" || req.FechaFin == "" || req.Horario == "" || req.Precio == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	var exists bool
	if err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM empleado WHERE id_empleado = $1)`, req.IDGuia).Scan(&exists); err != nil {
		http.Error(w, "Error al verificar guía", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "El guía especificado no existe", http.StatusBadRequest)
		return
	}

	var id int
	err := db.QueryRow(`
		INSERT INTO tour (id_guia, nombre, descripcion, fecha_inicio, fecha_fin, horario, precio)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id_tour
	`, req.IDGuia, req.Nombre, req.Descripcion, req.FechaInicio, req.FechaFin, req.Horario, req.Precio).Scan(&id)
	if err != nil {
		http.Error(w, "Error al crear tour", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id_tour": id})
}

func UpdateTourHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req models.TourRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body inválido", http.StatusBadRequest)
		return
	}

	if req.Nombre == "" || req.FechaInicio == "" || req.FechaFin == "" || req.Horario == "" || req.Precio == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`
		UPDATE tour
		SET id_guia = $1, nombre = $2, descripcion = $3,
		    fecha_inicio = $4, fecha_fin = $5, horario = $6, precio = $7
		WHERE id_tour = $8
	`, req.IDGuia, req.Nombre, req.Descripcion, req.FechaInicio, req.FechaFin, req.Horario, req.Precio, id)
	if err != nil {
		http.Error(w, "Error al actualizar tour", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Tour no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Tour actualizado"})
}

func DeleteTourHandler(w http.ResponseWriter, r *http.Request) {
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
	if err := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM tour WHERE id_tour = $1)`, id).Scan(&exists); err != nil {
		tx.Rollback()
		http.Error(w, "Error al verificar tour", http.StatusInternalServerError)
		return
	}
	if !exists {
		tx.Rollback()
		http.Error(w, "Tour no encontrado", http.StatusNotFound)
		return
	}

	if _, err := tx.Exec(`DELETE FROM cliente_tour WHERE id_tour = $1`, id); err != nil {
		tx.Rollback()
		http.Error(w, "Error al eliminar reservas del tour", http.StatusInternalServerError)
		return
	}

	if _, err := tx.Exec(`DELETE FROM tour WHERE id_tour = $1`, id); err != nil {
		tx.Rollback()
		http.Error(w, "Error al eliminar tour", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		http.Error(w, "Error al confirmar transacción", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Tour eliminado"})
}

func GetReservasHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT ct.id_cliente_tour, ct.id_cliente, ct.id_tour,
		       t.nombre AS nombre_tour, ct.fecha_reserva
		FROM cliente_tour ct
		JOIN tour t ON ct.id_tour = t.id_tour
	`)
	if err != nil {
		http.Error(w, "Error al obtener reservas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	reservas := []models.Reserva{}
	for rows.Next() {
		var res models.Reserva
		if err := rows.Scan(&res.ID, &res.IDCliente, &res.IDTour, &res.NombreTour, &res.FechaReserva); err != nil {
			http.Error(w, "Error al leer reserva", http.StatusInternalServerError)
			return
		}
		reservas = append(reservas, res)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservas)
}

func GetReservaByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var res models.Reserva
	err = db.QueryRow(`
		SELECT ct.id_cliente_tour, ct.id_cliente, ct.id_tour,
		       t.nombre AS nombre_tour, ct.fecha_reserva
		FROM cliente_tour ct
		JOIN tour t ON ct.id_tour = t.id_tour
		WHERE ct.id_cliente_tour = $1
	`, id).Scan(&res.ID, &res.IDCliente, &res.IDTour, &res.NombreTour, &res.FechaReserva)
	if err == sql.ErrNoRows {
		http.Error(w, "Reserva no encontrada", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Error al obtener reserva", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func CreateReservaHandler(w http.ResponseWriter, r *http.Request) {
	var req models.ReservaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body inválido", http.StatusBadRequest)
		return
	}

	if req.IDCliente == 0 || req.IDTour == 0 || req.FechaReserva == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	var clienteExists, tourExists bool
	if err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM cliente WHERE id_cliente = $1)`, req.IDCliente).Scan(&clienteExists); err != nil {
		http.Error(w, "Error al verificar cliente", http.StatusInternalServerError)
		return
	}
	if err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM tour WHERE id_tour = $1)`, req.IDTour).Scan(&tourExists); err != nil {
		http.Error(w, "Error al verificar tour", http.StatusInternalServerError)
		return
	}
	if !clienteExists {
		http.Error(w, "El cliente especificado no existe", http.StatusBadRequest)
		return
	}
	if !tourExists {
		http.Error(w, "El tour especificado no existe", http.StatusBadRequest)
		return
	}

	var id int
	err := db.QueryRow(`
		INSERT INTO cliente_tour (id_cliente, id_tour, fecha_reserva)
		VALUES ($1, $2, $3)
		RETURNING id_cliente_tour
	`, req.IDCliente, req.IDTour, req.FechaReserva).Scan(&id)
	if err != nil {
		http.Error(w, "Error al crear reserva", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id_cliente_tour": id})
}

func UpdateReservaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req models.ReservaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body inválido", http.StatusBadRequest)
		return
	}

	if req.IDCliente == 0 || req.IDTour == 0 || req.FechaReserva == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`
		UPDATE cliente_tour
		SET id_cliente = $1, id_tour = $2, fecha_reserva = $3
		WHERE id_cliente_tour = $4
	`, req.IDCliente, req.IDTour, req.FechaReserva, id)
	if err != nil {
		http.Error(w, "Error al actualizar reserva", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Reserva no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Reserva actualizada"})
}

func DeleteReservaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`DELETE FROM cliente_tour WHERE id_cliente_tour = $1`, id)
	if err != nil {
		http.Error(w, "Error al eliminar reserva", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Reserva no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Reserva eliminada"})
}
