package models

type Pintura struct {
	ID            int      `json:"id_pintura"`
	Titulo        string   `json:"titulo"`
	Descripcion   string   `json:"descripcion"`
	FechaCreacion string   `json:"fecha_creacion"`
	Precio        float64  `json:"precio"`
	Exclusiva     bool     `json:"exclusiva"`
	ImagenPath    string   `json:"imagen_path"`
	ImagenTipo    string   `json:"imagen_tipo"`
	ImagenNombre  string   `json:"imagen_nombre"`
	Artista       string   `json:"artista"`
	Coleccion     string   `json:"coleccion,omitempty"`
	Tecnicas      []string `json:"tecnicas,omitempty"`
}

type PinturaRequest struct {
	Titulo        string  `json:"titulo"`
	Descripcion   string  `json:"descripcion"`
	FechaCreacion string  `json:"fecha_creacion"`
	Precio        float64 `json:"precio"`
	Exclusiva     bool    `json:"exclusiva"`
	ImagenPath    string  `json:"imagen_path"`
	ImagenTipo    string  `json:"imagen_tipo"`
	ImagenNombre  string  `json:"imagen_nombre"`
	IDArtista     int     `json:"id_artista"`
	IDColeccion   *int    `json:"id_coleccion"` 
	Tecnicas      []int   `json:"tecnicas"`     
}