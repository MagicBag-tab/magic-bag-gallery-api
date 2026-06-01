package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"magic-bag-gallery-api/internal/middleware"
	"magic-bag-gallery-api/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Correo     string `json:"correo_electronico"`
	Contrasena string `json:"contrasena"`
}

type RegisterClienteRequest struct {
	Nombre            string `json:"nombre"`
	Apellido          string `json:"apellido"`
	CorreoElectronico string `json:"correo_electronico"`
	Telefono          string `json:"telefono"`
	Contrasena        string `json:"contrasena"`
	TipoCliente       string `json:"tipo_cliente"`
}

type RegisterEmpleadoRequest struct {
	Nombre            string `json:"nombre"`
	Apellido          string `json:"apellido"`
	CorreoElectronico string `json:"correo_electronico"`
	Telefono          string `json:"telefono"`
	Contrasena        string `json:"contrasena"`
	TipoEmpleado      string `json:"tipo_empleado"`
}

type authResponse struct {
	Role         string `json:"role"`
	IDUsuario    int    `json:"id_usuario"`
	IDCliente    int    `json:"id_cliente,omitempty"`
	IDEmpleado   int    `json:"id_empleado,omitempty"`
	Nombre       string `json:"nombre"`
	TipoEmpleado string `json:"tipo_empleado,omitempty"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body invalido", http.StatusBadRequest)
		return
	}

	var usuario models.UsuarioEntity
	if err := gormDB.Where("correo_electronico = ?", req.Correo).First(&usuario).Error; err != nil {
		http.Error(w, "Credenciales invalidas", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usuario.Contrasena), []byte(req.Contrasena)); err != nil {
		http.Error(w, "Credenciales invalidas", http.StatusUnauthorized)
		return
	}

	role := "cliente"
	tipoEmpleado := ""
	idEmpleado := 0
	var empleado models.EmpleadoEntity
	if err := gormDB.Where("id_usuario = ?", usuario.IDUsuario).First(&empleado).Error; err == nil {
		role = "empleado"
		tipoEmpleado = empleado.TipoEmpleado
		idEmpleado = empleado.IDEmpleado
	}

	if err := middleware.SaveSessionUser(w, r, usuario.IDUsuario, role, usuario.Nombre, tipoEmpleado); err != nil {
		http.Error(w, "Error al guardar sesion", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authResponse{
		Role:         role,
		IDUsuario:    usuario.IDUsuario,
		IDEmpleado:   idEmpleado,
		Nombre:       usuario.Nombre,
		TipoEmpleado: tipoEmpleado,
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if err := middleware.ClearSession(w, r); err != nil {
		http.Error(w, "Error al cerrar sesion", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Sesion cerrada"})
}

func SessionHandler(w http.ResponseWriter, r *http.Request) {
	idUsuario, _ := r.Context().Value(middleware.UserIDKey).(int)
	role, _ := r.Context().Value(middleware.RoleKey).(string)
	nombre, _ := r.Context().Value(middleware.NombreKey).(string)
	tipoEmpleado, _ := r.Context().Value(middleware.TipoEmpleadoKey).(string)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authResponse{
		Role:         role,
		IDUsuario:    idUsuario,
		Nombre:       nombre,
		TipoEmpleado: tipoEmpleado,
	})
}

func RegisterClienteHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterClienteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body invalido", http.StatusBadRequest)
		return
	}

	if req.Nombre == "" || req.Apellido == "" || req.CorreoElectronico == "" || req.Contrasena == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}
	if req.TipoCliente == "" {
		req.TipoCliente = "regular"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Contrasena), 12)
	if err != nil {
		http.Error(w, "Error al procesar contrasena", http.StatusInternalServerError)
		return
	}

	idUsuario, idCliente, err := registrarClienteConProcedure(req, string(hash))
	if err != nil {
		http.Error(w, "El correo ya esta registrado", http.StatusConflict)
		return
	}

	if err := middleware.SaveSessionUser(w, r, idUsuario, "cliente", req.Nombre, ""); err != nil {
		http.Error(w, "Error al guardar sesion", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(authResponse{
		Role:      "cliente",
		IDUsuario: idUsuario,
		IDCliente: idCliente,
		Nombre:    req.Nombre,
	})
}

func RegisterEmpleadoHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterEmpleadoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Body invalido", http.StatusBadRequest)
		return
	}

	if req.Nombre == "" || req.Apellido == "" || req.CorreoElectronico == "" || req.Contrasena == "" || req.TipoEmpleado == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Contrasena), 12)
	if err != nil {
		http.Error(w, "Error al procesar contrasena", http.StatusInternalServerError)
		return
	}

	idUsuario, idEmpleado, err := registrarEmpleadoConProcedure(req, string(hash))
	if err != nil {
		http.Error(w, "El correo ya esta registrado", http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(authResponse{
		Role:         "empleado",
		IDUsuario:    idUsuario,
		IDEmpleado:   idEmpleado,
		Nombre:       req.Nombre,
		TipoEmpleado: req.TipoEmpleado,
	})
}

func registrarClienteConProcedure(req RegisterClienteRequest, hash string) (int, int, error) {
	type result struct {
		IDUsuario int `gorm:"column:id_usuario"`
		IDCliente int `gorm:"column:id_cliente"`
	}
	var res result
	err := gormDB.Raw(
		"SELECT * FROM sp_registrar_cliente(?, ?, ?, ?, ?, ?)",
		req.Nombre, req.Apellido, req.CorreoElectronico, req.Telefono, hash, req.TipoCliente,
	).Scan(&res).Error
	if err == nil && res.IDUsuario != 0 && res.IDCliente != 0 {
		return res.IDUsuario, res.IDCliente, nil
	}

	return registrarClienteGORM(req, hash)
}

func registrarClienteGORM(req RegisterClienteRequest, hash string) (int, int, error) {
	var usuario models.UsuarioEntity
	var cliente models.ClienteEntity
	err := gormDB.Transaction(func(tx *gorm.DB) error {
		usuario = models.UsuarioEntity{
			Nombre:            req.Nombre,
			Apellido:          req.Apellido,
			CorreoElectronico: req.CorreoElectronico,
			Telefono:          req.Telefono,
			Contrasena:        hash,
		}
		if err := tx.Create(&usuario).Error; err != nil {
			return err
		}
		cliente = models.ClienteEntity{IDUsuario: usuario.IDUsuario, TipoCliente: req.TipoCliente}
		return tx.Create(&cliente).Error
	})
	return usuario.IDUsuario, cliente.IDCliente, err
}

func registrarEmpleadoConProcedure(req RegisterEmpleadoRequest, hash string) (int, int, error) {
	type result struct {
		IDUsuario  int `gorm:"column:id_usuario"`
		IDEmpleado int `gorm:"column:id_empleado"`
	}
	var res result
	err := gormDB.Raw(
		"SELECT * FROM sp_registrar_empleado(?, ?, ?, ?, ?, ?)",
		req.Nombre, req.Apellido, req.CorreoElectronico, req.Telefono, hash, req.TipoEmpleado,
	).Scan(&res).Error
	if err == nil && res.IDUsuario != 0 && res.IDEmpleado != 0 {
		return res.IDUsuario, res.IDEmpleado, nil
	}

	return registrarEmpleadoGORM(req, hash)
}

func registrarEmpleadoGORM(req RegisterEmpleadoRequest, hash string) (int, int, error) {
	var usuario models.UsuarioEntity
	var empleado models.EmpleadoEntity
	err := gormDB.Transaction(func(tx *gorm.DB) error {
		usuario = models.UsuarioEntity{
			Nombre:            req.Nombre,
			Apellido:          req.Apellido,
			CorreoElectronico: req.CorreoElectronico,
			Telefono:          req.Telefono,
			Contrasena:        hash,
		}
		if err := tx.Create(&usuario).Error; err != nil {
			return err
		}
		empleado = models.EmpleadoEntity{IDUsuario: usuario.IDUsuario, TipoEmpleado: req.TipoEmpleado}
		return tx.Create(&empleado).Error
	})
	return usuario.IDUsuario, empleado.IDEmpleado, err
}

func isNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
