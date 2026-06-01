package handlers

import (
	"encoding/json"
	"magic-bag-gallery-api/internal/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetToursHandler(w http.ResponseWriter, r *http.Request) {
	tours, err := loadTours("")
	if err != nil {
		http.Error(w, "Error al obtener tours", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tours)
}

func GetTourByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	tours, err := loadTours("WHERE t.id_tour = ?", id)
	if err != nil {
		http.Error(w, "Error al obtener tour", http.StatusInternalServerError)
		return
	}
	if len(tours) == 0 {
		http.Error(w, "Tour no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tours[0])
}

func CreateTourHandler(w http.ResponseWriter, r *http.Request) {
	var req models.TourRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body invalido", http.StatusBadRequest)
		return
	}

	if req.IDGuia == 0 || req.Nombre == "" || req.FechaInicio == "" || req.FechaFin == "" || req.Horario == "" || req.Precio == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	if !empleadoExiste(req.IDGuia) {
		http.Error(w, "El guia especificado no existe", http.StatusBadRequest)
		return
	}

	record := models.TourEntity{
		IDGuia:      req.IDGuia,
		Nombre:      req.Nombre,
		Descripcion: req.Descripcion,
		FechaInicio: req.FechaInicio,
		FechaFin:    req.FechaFin,
		Horario:     req.Horario,
		Precio:      req.Precio,
	}
	if err := gormDB.Create(&record).Error; err != nil {
		http.Error(w, "Error al crear tour", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id_tour": record.IDTour})
}

func UpdateTourHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	var req models.TourRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body invalido", http.StatusBadRequest)
		return
	}

	if req.Nombre == "" || req.FechaInicio == "" || req.FechaFin == "" || req.Horario == "" || req.Precio == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	result := gormDB.Model(&models.TourEntity{}).Where("id_tour = ?", id).Updates(map[string]interface{}{
		"id_guia":      req.IDGuia,
		"nombre":       req.Nombre,
		"descripcion":  req.Descripcion,
		"fecha_inicio": req.FechaInicio,
		"fecha_fin":    req.FechaFin,
		"horario":      req.Horario,
		"precio":       req.Precio,
	})
	if result.Error != nil {
		http.Error(w, "Error al actualizar tour", http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		http.Error(w, "Tour no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Tour actualizado"})
}

func DeleteTourHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	if err := gormDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id_tour = ?", id).Delete(&models.ClienteTourEntity{}).Error; err != nil {
			return err
		}
		result := tx.Where("id_tour = ?", id).Delete(&models.TourEntity{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	}); err != nil {
		if isNotFound(err) {
			http.Error(w, "Tour no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al eliminar tour", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Tour eliminado"})
}

func GetReservasHandler(w http.ResponseWriter, r *http.Request) {
	var reservas []models.Reserva
	err := gormDB.Raw(`
		SELECT
			ct.id_cliente_tour,
			ct.id_cliente,
			ct.id_tour,
			uc.nombre || ' ' || uc.apellido AS nombre_cliente,
			t.nombre AS nombre_tour,
			TO_CHAR(ct.fecha_reserva, 'YYYY-MM-DD') AS fecha_reserva
		FROM cliente_tour ct
		JOIN cliente c ON ct.id_cliente = c.id_cliente
		JOIN usuario uc ON c.id_usuario = uc.id_usuario
		JOIN tour t ON ct.id_tour = t.id_tour
		ORDER BY ct.fecha_reserva DESC
	`).Scan(&reservas).Error
	if err != nil {
		http.Error(w, "Error al obtener reservas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservas)
}

func GetReservaByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	var reserva models.Reserva
	err = gormDB.Raw(`
		SELECT
			ct.id_cliente_tour,
			ct.id_cliente,
			ct.id_tour,
			uc.nombre || ' ' || uc.apellido AS nombre_cliente,
			t.nombre AS nombre_tour,
			TO_CHAR(ct.fecha_reserva, 'YYYY-MM-DD') AS fecha_reserva
		FROM cliente_tour ct
		JOIN cliente c ON ct.id_cliente = c.id_cliente
		JOIN usuario uc ON c.id_usuario = uc.id_usuario
		JOIN tour t ON ct.id_tour = t.id_tour
		WHERE ct.id_cliente_tour = ?
	`, id).Scan(&reserva).Error
	if err != nil || reserva.ID == 0 {
		http.Error(w, "Reserva no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reserva)
}

func CreateReservaHandler(w http.ResponseWriter, r *http.Request) {
	var req models.ReservaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body invalido", http.StatusBadRequest)
		return
	}

	if req.IDCliente == 0 {
		idUsuario, ok := getUserIDFromContext(r)
		if !ok {
			http.Error(w, "Sesion invalida", http.StatusUnauthorized)
			return
		}
		var cliente models.ClienteEntity
		if err := gormDB.First(&cliente, "id_usuario = ?", idUsuario).Error; err != nil {
			http.Error(w, "Cliente no encontrado", http.StatusNotFound)
			return
		}
		req.IDCliente = cliente.IDCliente
	}

	if req.IDCliente == 0 || req.IDTour == 0 || req.FechaReserva == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	id, err := crearReservaConProcedure(req)
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
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	var req models.ReservaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body invalido", http.StatusBadRequest)
		return
	}

	result := gormDB.Model(&models.ClienteTourEntity{}).Where("id_cliente_tour = ?", id).Updates(map[string]interface{}{
		"id_cliente":    req.IDCliente,
		"id_tour":       req.IDTour,
		"fecha_reserva": req.FechaReserva,
	})
	if result.Error != nil {
		http.Error(w, "Error al actualizar reserva", http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		http.Error(w, "Reserva no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Reserva actualizada"})
}

func DeleteReservaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	result := gormDB.Where("id_cliente_tour = ?", id).Delete(&models.ClienteTourEntity{})
	if result.Error != nil {
		http.Error(w, "Error al eliminar reserva", http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		http.Error(w, "Reserva no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Reserva eliminada"})
}

func loadTours(where string, args ...interface{}) ([]models.Tour, error) {
	query := `
		SELECT
			t.id_tour,
			t.id_guia,
			u.nombre || ' ' || u.apellido AS nombre_guia,
			t.nombre,
			t.descripcion,
			TO_CHAR(t.fecha_inicio, 'YYYY-MM-DD') AS fecha_inicio,
			TO_CHAR(t.fecha_fin, 'YYYY-MM-DD') AS fecha_fin,
			t.horario,
			t.precio::text AS precio
		FROM tour t
		JOIN empleado e ON t.id_guia = e.id_empleado
		JOIN usuario u ON e.id_usuario = u.id_usuario
		` + where + `
		ORDER BY t.fecha_inicio, t.id_tour
	`
	var tours []models.Tour
	return tours, gormDB.Raw(query, args...).Scan(&tours).Error
}

func empleadoExiste(idEmpleado int) bool {
	var total int64
	err := gormDB.Model(&models.EmpleadoEntity{}).Where("id_empleado = ?", idEmpleado).Count(&total).Error
	return err == nil && total > 0
}

func crearReservaConProcedure(req models.ReservaRequest) (int, error) {
	type result struct {
		IDClienteTour int `gorm:"column:id_cliente_tour"`
	}
	var res result
	err := gormDB.Raw("SELECT * FROM sp_crear_reserva(?, ?, ?)", req.IDCliente, req.IDTour, req.FechaReserva).Scan(&res).Error
	if err == nil && res.IDClienteTour != 0 {
		return res.IDClienteTour, nil
	}

	record := models.ClienteTourEntity{
		IDCliente:    req.IDCliente,
		IDTour:       req.IDTour,
		FechaReserva: req.FechaReserva,
	}
	if err := gormDB.Create(&record).Error; err != nil {
		return 0, err
	}
	return record.IDClienteTour, nil
}
