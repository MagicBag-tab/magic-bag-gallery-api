package models

type Artista struct {
	ID             int      `json:"id_artista"`
	NombreCompleto string   `json:"nombre_completo"`
	Nacionalidad   string   `json:"nacionalidad"`
	IDReclutador   int      `json:"id_reclutador"`
	Pinturas       []string `json:"pinturas,omitempty"`
	Colecciones    []string `json:"colecciones,omitempty"`
}

type ArtistaRequest struct {
	NombreCompleto string `json:"nombre_completo"`
	Nacionalidad   string `json:"nacionalidad"`
	IDReclutador   int    `json:"id_reclutador"`
	IDPinturas     []int  `json:"id_pinturas,omitempty"`
	IDColecciones  []int  `json:"id_colecciones,omitempty"`
}
