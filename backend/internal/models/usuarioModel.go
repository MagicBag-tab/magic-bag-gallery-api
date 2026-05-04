package models

type Usuario struct {
	ID                int    `json:"id_usuario"`
	Nombre            string `json:"nombre"`
	Apellido          string `json:"apellido"`
	CorreoElectronico string `json:"correo_electronico"`
	Telefono          string `json:"telefono"`
}

type UsuarioRequest struct {
	Nombre            string `json:"nombre"`
	Apellido          string `json:"apellido"`
	CorreoElectronico string `json:"correo_electronico"`
	Telefono          string `json:"telefono"`
	Contrasena        string `json:"contrasena,omitempty"`
}
