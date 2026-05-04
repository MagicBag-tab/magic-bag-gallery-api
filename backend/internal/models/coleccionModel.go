package models

type Coleccion struct {
	ID          int      `json:"id_coleccion"`
	Nombre      string   `json:"nombre"`
	Descripcion string   `json:"descripcion"`
	Artista     string   `json:"artista"`
	Pinturas    []string `json:"pinturas,omitempty"`
}

type ColeccionRequest struct {
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	IDArtista   int    `json:"id_artista"`
	IDPinturas  []int  `json:"id_pinturas,omitempty"`
}
