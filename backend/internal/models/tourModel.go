package models

type Tour struct {
	ID          int    `json:"id_tour"`
	IDGuia      int    `json:"id_guia"`
	NombreGuia  string `json:"nombre_guia,omitempty"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	FechaInicio string `json:"fecha_inicio"`
	FechaFin    string `json:"fecha_fin"`
	Horario     string `json:"horario"`
	Precio      string `json:"precio"`
}

type TourRequest struct {
	IDGuia      int    `json:"id_guia"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	FechaInicio string `json:"fecha_inicio"`
	FechaFin    string `json:"fecha_fin"`
	Horario     string `json:"horario"`
	Precio      string `json:"precio"`
}

type Reserva struct {
	ID           int    `json:"id_cliente_tour"`
	IDCliente    int    `json:"id_cliente"`
	IDTour       int    `json:"id_tour"`
	NombreTour   string `json:"nombre_tour,omitempty"`
	FechaReserva string `json:"fecha_reserva"`
}

type ReservaRequest struct {
	IDCliente    int    `json:"id_cliente"`
	IDTour       int    `json:"id_tour"`
	FechaReserva string `json:"fecha_reserva"`
}
