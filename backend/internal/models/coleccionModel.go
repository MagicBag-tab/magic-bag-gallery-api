package models

type Coleccion struct {
	ID               int      `json:"id_coleccion"`
	Nombre           string   `json:"nombre"`
	Descripcion      string   `json:"descripcion"`
	Exclusiva        bool     `json:"exclusiva"`
	FechaLanzamiento string   `json:"fecha_lanzamiento"`
	Pinturas         []string `json:"pinturas,omitempty"`
}

type ColeccionRequest struct {
	Nombre           string `json:"nombre"`
	Descripcion      string `json:"descripcion"`
	Exclusiva        bool   `json:"exclusiva"`
	FechaLanzamiento string `json:"fecha_lanzamiento"`
	IDPinturas       []int  `json:"id_pinturas,omitempty"`
}
