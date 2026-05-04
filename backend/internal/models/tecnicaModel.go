package models

type Tecnica struct {
	ID          int    `json:"id_tecnica"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
}

type TecnicaRequest struct {
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
}
